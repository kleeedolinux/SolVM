package vm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const (
	maxCacheSize  = 100
	httpTimeout   = 30 * time.Second
	maxModuleSize = 10 << 20
	luaExtension  = ".lua"
	modulesDir    = "modules"
)

type ModuleCache struct {
	code      string
	timestamp time.Time
	size      int
}

type ImportModule struct {
	vm         *SolVM
	loaded     map[string]bool
	cache      map[string]*ModuleCache
	mu         sync.RWMutex
	httpClient *http.Client
	cacheSize  int
}

func NewImportModule(vm *SolVM) *ImportModule {
	return &ImportModule{
		vm:     vm,
		loaded: make(map[string]bool),
		cache:  make(map[string]*ModuleCache),
		httpClient: &http.Client{
			Timeout: httpTimeout,
			Transport: &http.Transport{
				MaxIdleConns:        10,
				IdleConnTimeout:     30 * time.Second,
				DisableCompression:  false,
				MaxIdleConnsPerHost: 2,
			},
		},
	}
}

func (im *ImportModule) Register() {
	im.vm.RegisterFunction("import", im.importModule)
}

func (im *ImportModule) importModule(L *lua.LState) int {
	modulePath := L.CheckString(1)

	if modulePath == "" {
		L.RaiseError("Module path cannot be empty")
		return 0
	}

	im.mu.RLock()
	if im.loaded[modulePath] {
		im.mu.RUnlock()
		return 0
	}
	im.mu.RUnlock()

	if strings.HasSuffix(modulePath, "/") {
		return im.importFolder(L, modulePath)
	}

	code, err := im.loadModuleWithCache(modulePath)
	if err != nil {
		L.RaiseError("Failed to import module '%s': %v", modulePath, err)
		return 0
	}

	moduleState := L.NewTable()
	L.SetGlobal(modulePath, moduleState)

	if err := L.DoString(code); err != nil {
		L.RaiseError("Failed to execute module '%s': %v", modulePath, err)
		return 0
	}

	ret := L.Get(-1)
	L.Pop(1)

	if ret.Type() == lua.LTTable {
		im.copyTableContents(ret.(*lua.LTable), moduleState)
	}

	im.mu.Lock()
	im.loaded[modulePath] = true
	im.mu.Unlock()

	return 0
}

func (im *ImportModule) importFolder(L *lua.LState, folderPath string) int {
	if !strings.HasSuffix(folderPath, "/") {
		folderPath += "/"
	}

	moduleDir := filepath.Join(modulesDir, folderPath)
	entries, err := os.ReadDir(moduleDir)
	if err != nil {
		L.RaiseError("Failed to read module folder '%s': %v", folderPath, err)
		return 0
	}

	moduleState := L.NewTable()
	L.SetGlobal(folderPath, moduleState)

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), luaExtension) {
			continue
		}

		filePath := filepath.Join(folderPath, entry.Name())
		code, err := im.loadModuleWithCache(filePath)
		if err != nil {
			L.RaiseError("Failed to import module '%s': %v", filePath, err)
			continue
		}

		if err := L.DoString(code); err != nil {
			L.RaiseError("Failed to execute module '%s': %v", filePath, err)
			continue
		}

		ret := L.Get(-1)
		L.Pop(1)

		if ret.Type() == lua.LTTable {
			moduleName := strings.TrimSuffix(entry.Name(), luaExtension)
			subTable := L.NewTable()
			moduleState.RawSetString(moduleName, subTable)
			im.copyTableContents(ret.(*lua.LTable), subTable)
		}

		im.mu.Lock()
		im.loaded[filePath] = true
		im.mu.Unlock()
	}

	return 0
}

func (im *ImportModule) loadModuleWithCache(modulePath string) (string, error) {
	im.mu.RLock()
	if cached, exists := im.cache[modulePath]; exists {
		im.mu.RUnlock()
		return cached.code, nil
	}
	im.mu.RUnlock()

	code, err := im.loadModule(modulePath)
	if err != nil {
		return "", err
	}

	im.mu.Lock()
	defer im.mu.Unlock()

	if im.cacheSize >= maxCacheSize {
		im.evictOldestCache()
	}

	im.cache[modulePath] = &ModuleCache{
		code:      code,
		timestamp: time.Now(),
		size:      len(code),
	}
	im.cacheSize++

	return code, nil
}

func (im *ImportModule) loadModule(modulePath string) (string, error) {

	if im.isURL(modulePath) {
		return im.loadFromURL(modulePath)
	}

	if !strings.HasSuffix(modulePath, luaExtension) {
		modulePath += luaExtension
	}

	if code, err := im.readFileWithLimit(modulePath); err == nil {
		return code, nil
	}

	moduleDir := filepath.Join(modulesDir, modulePath)
	if code, err := im.readFileWithLimit(moduleDir); err == nil {
		return code, nil
	}

	return "", fmt.Errorf("module not found: %s", modulePath)
}

func (im *ImportModule) loadFromURL(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := im.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch module from URL: %w", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {

		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch module: HTTP %d", resp.StatusCode)
	}

	limitedReader := io.LimitReader(resp.Body, maxModuleSize)
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", fmt.Errorf("failed to read module content: %w", err)
	}

	return string(body), nil
}

func (im *ImportModule) readFileWithLimit(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {

		}
	}()

	info, err := file.Stat()
	if err != nil {
		return "", err
	}

	if info.Size() > maxModuleSize {
		return "", fmt.Errorf("module file too large: %d bytes", info.Size())
	}

	limitedReader := io.LimitReader(file, maxModuleSize)
	content, err := io.ReadAll(limitedReader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (im *ImportModule) copyTableContents(src, dst *lua.LTable) {
	src.ForEach(func(key, value lua.LValue) {
		dst.RawSet(key, value)
	})
}

func (im *ImportModule) isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

func (im *ImportModule) evictOldestCache() {
	var oldestKey string
	var oldestTime time.Time

	for key, entry := range im.cache {
		if oldestKey == "" || entry.timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.timestamp
		}
	}

	if oldestKey != "" {
		delete(im.cache, oldestKey)
		im.cacheSize--
	}
}

func (im *ImportModule) ClearCache() {
	im.mu.Lock()
	defer im.mu.Unlock()

	im.cache = make(map[string]*ModuleCache)
	im.cacheSize = 0
}

func (im *ImportModule) GetCacheStats() (size int, entries int) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	totalSize := 0
	for _, entry := range im.cache {
		totalSize += entry.size
	}

	return totalSize, len(im.cache)
}
