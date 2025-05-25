package vm

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type ImportModule struct {
	vm     *SolVM
	loaded map[string]bool
}

func NewImportModule(vm *SolVM) *ImportModule {
	return &ImportModule{
		vm:     vm,
		loaded: make(map[string]bool),
	}
}

func (im *ImportModule) Register() {
	im.vm.RegisterFunction("import", im.importModule)
}

func (im *ImportModule) importModule(L *lua.LState) int {
	modulePath := L.CheckString(1)

	if im.loaded[modulePath] {
		return 0
	}

	code, err := im.loadModule(modulePath)
	if err != nil {
		L.RaiseError("Failed to import module: %v", err)
		return 0
	}

	
	moduleState := L.NewTable()
	L.SetGlobal(modulePath, moduleState)

	
	if err := L.DoString(code); err != nil {
		L.RaiseError("Failed to execute module: %v", err)
		return 0
	}

	
	ret := L.Get(-1)
	L.Pop(1)

	
	if ret.Type() == lua.LTTable {
		retTable := ret.(*lua.LTable)
		retTable.ForEach(func(key, value lua.LValue) {
			moduleState.RawSet(key, value)
		})
	}

	im.loaded[modulePath] = true
	return 0
}

func (im *ImportModule) loadModule(modulePath string) (string, error) {
	
	if strings.HasPrefix(modulePath, "http://") || strings.HasPrefix(modulePath, "https://") {
		return im.loadFromURL(modulePath)
	}

	if !strings.HasSuffix(modulePath, ".lua") {
		modulePath += ".lua"
	}

	
	code, err := os.ReadFile(modulePath)
	if err == nil {
		return string(code), nil
	}

	
	moduleDir := filepath.Join("modules", modulePath)
	code, err = os.ReadFile(moduleDir)
	if err == nil {
		return string(code), nil
	}

	return "", fmt.Errorf("module not found: %s", modulePath)
}

func (im *ImportModule) loadFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch module from URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch module: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read module content: %v", err)
	}

	return string(body), nil
}
