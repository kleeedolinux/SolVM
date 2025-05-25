package vm

import (
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type ConcurrencyModule struct {
	vm       *SolVM
	channels map[string]chan lua.LValue
	mu       sync.RWMutex
	wg       sync.WaitGroup
}

func NewConcurrencyModule(vm *SolVM) *ConcurrencyModule {
	return &ConcurrencyModule{
		vm:       vm,
		channels: make(map[string]chan lua.LValue),
	}
}

func (cm *ConcurrencyModule) Register() {
	cm.vm.RegisterFunction("go", cm.goFunc)
	cm.vm.RegisterFunction("chan", cm.createChannel)
	cm.vm.RegisterFunction("send", cm.sendToChannel)
	cm.vm.RegisterFunction("receive", cm.receiveFromChannel)
	cm.vm.RegisterFunction("select", cm.selectChannel)
	cm.vm.RegisterFunction("wait", cm.waitForGoroutines)
}

func (cm *ConcurrencyModule) goFunc(L *lua.LState) int {
	fn := L.CheckFunction(1)

	cm.wg.Add(1)
	go func() {
		defer cm.wg.Done()

		
		L2 := lua.NewState()
		defer L2.Close()

		
		fn2 := L2.NewFunctionFromProto(fn.Proto)
		L2.Push(fn2)

		
		L2.SetGlobal("print", L.GetGlobal("print"))
		L2.SetGlobal("sleep", L.GetGlobal("sleep"))
		L2.SetGlobal("send", L.GetGlobal("send"))
		L2.SetGlobal("receive", L.GetGlobal("receive"))
		L2.SetGlobal("select", L.GetGlobal("select"))

		
		if err := L2.PCall(0, 0, nil); err != nil {
			L.RaiseError("Goroutine error: %v", err)
		}
	}()

	return 0
}

func (cm *ConcurrencyModule) createChannel(L *lua.LState) int {
	name := L.CheckString(1)
	bufferSize := L.OptInt(2, 0)

	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.channels[name]; exists {
		L.RaiseError("Channel %s already exists", name)
		return 0
	}

	cm.channels[name] = make(chan lua.LValue, bufferSize)
	return 0
}

func (cm *ConcurrencyModule) sendToChannel(L *lua.LState) int {
	name := L.CheckString(1)
	value := L.CheckAny(2)

	cm.mu.RLock()
	ch, exists := cm.channels[name]
	cm.mu.RUnlock()

	if !exists {
		L.RaiseError("Channel %s does not exist", name)
		return 0
	}

	select {
	case ch <- value:
		L.Push(lua.LTrue)
	case <-time.After(time.Second):
		L.Push(lua.LFalse)
	}

	return 1
}

func (cm *ConcurrencyModule) receiveFromChannel(L *lua.LState) int {
	name := L.CheckString(1)

	cm.mu.RLock()
	ch, exists := cm.channels[name]
	cm.mu.RUnlock()

	if !exists {
		L.RaiseError("Channel %s does not exist", name)
		return 0
	}

	select {
	case value := <-ch:
		L.Push(value)
	case <-time.After(time.Second):
		L.Push(lua.LNil)
	}

	return 1
}

func (cm *ConcurrencyModule) selectChannel(L *lua.LState) int {
	if L.GetTop() < 2 {
		L.RaiseError("select requires at least one channel")
		return 0
	}

	channels := make([]string, 0, L.GetTop())
	for i := 1; i <= L.GetTop(); i++ {
		channels = append(channels, L.CheckString(i))
	}

	cm.mu.RLock()
	defer cm.mu.RUnlock()

	cases := make([]chan lua.LValue, 0, len(channels))
	for _, name := range channels {
		if ch, exists := cm.channels[name]; exists {
			cases = append(cases, ch)
		}
	}

	if len(cases) == 0 {
		L.RaiseError("No valid channels provided")
		return 0
	}

	select {
	case value := <-cases[0]:
		L.Push(value)
		L.Push(lua.LString(channels[0]))
	default:
		L.Push(lua.LNil)
		L.Push(lua.LNil)
	}

	return 2
}

func (cm *ConcurrencyModule) waitForGoroutines(L *lua.LState) int {
	cm.wg.Wait()
	return 0
}
