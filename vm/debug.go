package vm

import (
	"fmt"
	"os"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type DebugModule struct {
	vm       *SolVM
	watchers map[string]*time.Ticker
	mu       sync.RWMutex
	lastMod  map[string]time.Time
}

func NewDebugModule(vm *SolVM) *DebugModule {
	return &DebugModule{
		vm:       vm,
		watchers: make(map[string]*time.Ticker),
		lastMod:  make(map[string]time.Time),
	}
}

func (dm *DebugModule) Register() {
	dm.vm.RegisterFunction("watch_file", dm.watchFile)
	dm.vm.RegisterFunction("reload_script", dm.reloadScript)
	dm.vm.RegisterFunction("trace", dm.trace)
}

func (dm *DebugModule) watchFile(L *lua.LState) int {
	filePath := L.CheckString(1)
	callback := L.CheckFunction(2)

	
	info, err := os.Stat(filePath)
	if err != nil {
		L.RaiseError("Failed to watch file: %v", err)
		return 0
	}

	dm.mu.Lock()
	dm.lastMod[filePath] = info.ModTime()
	dm.mu.Unlock()

	
	ticker := time.NewTicker(1 * time.Second)
	dm.watchers[filePath] = ticker

	go func() {
		for range ticker.C {
			info, err := os.Stat(filePath)
			if err != nil {
				continue
			}

			dm.mu.RLock()
			lastMod := dm.lastMod[filePath]
			dm.mu.RUnlock()

			if info.ModTime().After(lastMod) {
				dm.mu.Lock()
				dm.lastMod[filePath] = info.ModTime()
				dm.mu.Unlock()

				
				L2 := lua.NewState()
				defer L2.Close()

				fn := L2.NewFunctionFromProto(callback.Proto)
				L2.Push(fn)
				if err := L2.PCall(0, 0, nil); err != nil {
					dm.vm.monitor.handleError(fmt.Errorf("Watch callback error: %v", err))
				}
			}
		}
	}()

	return 0
}

func (dm *DebugModule) reloadScript(L *lua.LState) int {
	
	scriptPath := L.GetGlobal("_SCRIPT_PATH").String()
	if scriptPath == "" {
		L.RaiseError("No script path available for reloading")
		return 0
	}

	
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		L.RaiseError("Failed to read script: %v", err)
		return 0
	}

	
	L2 := lua.NewState()
	defer L2.Close()

	
	dm.vm.RegisterCustomFunctions()

	
	if err := L2.DoString(string(content)); err != nil {
		L.RaiseError("Failed to reload script: %v", err)
		return 0
	}

	return 0
}

func (dm *DebugModule) trace(L *lua.LState) int {
	
	debug := L.GetGlobal("debug")
	if debug.Type() != lua.LTTable {
		L.RaiseError("debug library not available")
		return 0
	}

	traceback := debug.(*lua.LTable).RawGetString("traceback")
	if traceback.Type() != lua.LTFunction {
		L.RaiseError("debug.traceback not available")
		return 0
	}

	L.Push(traceback)
	if err := L.PCall(0, 1, nil); err != nil {
		L.RaiseError("Failed to get stack trace: %v", err)
		return 0
	}

	
	fmt.Println(L.Get(-1).String())
	return 0
}

func (dm *DebugModule) StopWatching(filePath string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	if ticker, exists := dm.watchers[filePath]; exists {
		ticker.Stop()
		delete(dm.watchers, filePath)
		delete(dm.lastMod, filePath)
	}
}
