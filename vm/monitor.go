package vm

import (
	"fmt"
	"runtime"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type MonitorModule struct {
	vm            *SolVM
	startMem      runtime.MemStats
	lastMem       runtime.MemStats
	goroutineMap  map[int]string
	goroutineMu   sync.RWMutex
	errorHandlers []func(error)
}

func NewMonitorModule(vm *SolVM) *MonitorModule {
	monitor := &MonitorModule{
		vm:           vm,
		goroutineMap: make(map[int]string),
	}
	runtime.ReadMemStats(&monitor.startMem)
	monitor.lastMem = monitor.startMem
	return monitor
}

func (mm *MonitorModule) Register() {
	mm.vm.RegisterFunction("on_error", mm.registerErrorHandler)
	mm.vm.RegisterFunction("check_memory", mm.checkMemory)
	mm.vm.RegisterFunction("get_goroutines", mm.getGoroutines)
}

func (mm *MonitorModule) registerErrorHandler(L *lua.LState) int {
	fn := L.CheckFunction(1)

	handler := func(err error) {
		L2 := lua.NewState()
		defer L2.Close()

		fn2 := L2.NewFunctionFromProto(fn.Proto)
		L2.Push(fn2)
		L2.Push(lua.LString(mm.formatError(err)))
		L2.PCall(1, 0, nil)
	}

	mm.errorHandlers = append(mm.errorHandlers, handler)
	return 0
}

func (mm *MonitorModule) checkMemory(L *lua.LState) int {
	var currentMem runtime.MemStats
	runtime.ReadMemStats(&currentMem)

	allocDiff := currentMem.Alloc - mm.lastMem.Alloc
	totalAllocDiff := currentMem.TotalAlloc - mm.lastMem.TotalAlloc
	sysDiff := currentMem.Sys - mm.lastMem.Sys

	mm.lastMem = currentMem

	
	stats := L.NewTable()
	stats.RawSetString("alloc_diff", lua.LNumber(allocDiff))
	stats.RawSetString("total_alloc_diff", lua.LNumber(totalAllocDiff))
	stats.RawSetString("sys_diff", lua.LNumber(sysDiff))
	stats.RawSetString("num_gc", lua.LNumber(currentMem.NumGC))
	stats.RawSetString("goroutines", lua.LNumber(runtime.NumGoroutine()))

	L.Push(stats)
	return 1
}

func (mm *MonitorModule) getGoroutines(L *lua.LState) int {
	mm.goroutineMu.RLock()
	defer mm.goroutineMu.RUnlock()

	goroutines := L.NewTable()
	for id, name := range mm.goroutineMap {
		goroutines.RawSetInt(id, lua.LString(name))
	}

	L.Push(goroutines)
	return 1
}

func (mm *MonitorModule) trackGoroutine(id int, name string) {
	mm.goroutineMu.Lock()
	defer mm.goroutineMu.Unlock()
	mm.goroutineMap[id] = name
}

func (mm *MonitorModule) untrackGoroutine(id int) {
	mm.goroutineMu.Lock()
	defer mm.goroutineMu.Unlock()
	delete(mm.goroutineMap, id)
}

func (mm *MonitorModule) handleError(err error) {
	if err == nil {
		return
	}

	
	for _, handler := range mm.errorHandlers {
		handler(err)
	}
}

func (mm *MonitorModule) formatError(err error) string {
	if err == nil {
		return ""
	}

	switch e := err.(type) {
	case *lua.ApiError:
		return fmt.Sprintf("Lua Error: %s", e.Error())
	default:
		return err.Error()
	}
}
