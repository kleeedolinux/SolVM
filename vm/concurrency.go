package vm

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type Channel struct {
	ch     chan lua.LValue
	closed bool
	mu     sync.RWMutex
}

type ConcurrencyModule struct {
	vm       *SolVM
	channels map[string]*Channel
	mu       sync.RWMutex
	wg       sync.WaitGroup
	pool     *sync.Pool
	done     chan struct{}
}

func NewConcurrencyModule(vm *SolVM) *ConcurrencyModule {
	return &ConcurrencyModule{
		vm:       vm,
		channels: make(map[string]*Channel, 10),
		done:     make(chan struct{}),
		pool: &sync.Pool{
			New: func() interface{} {
				L := lua.NewState()
				L.SetMx(1000)
				return L
			},
		},
	}
}

func (cm *ConcurrencyModule) Register() {
	cm.vm.RegisterFunction("go", cm.goFunc)
	cm.vm.RegisterFunction("chan", cm.createChannel)
	cm.vm.RegisterFunction("send", cm.sendToChannel)
	cm.vm.RegisterFunction("receive", cm.receiveFromChannel)
	cm.vm.RegisterFunction("select", cm.selectChannel)
	cm.vm.RegisterFunction("wait", cm.waitForGoroutines)
	cm.vm.RegisterFunction("close_channel", cm.closeChannel)
}

func (cm *ConcurrencyModule) goFunc(L *lua.LState) int {
	fn := L.CheckFunction(1)

	if cm.vm.maxGoroutines > 0 {
		if runtime.NumGoroutine() >= cm.vm.maxGoroutines {
			L.RaiseError("maximum number of goroutines reached")
			return 0
		}
	}

	cm.wg.Add(1)
	go func() {
		defer cm.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				cm.vm.monitor.handleError(fmt.Errorf("goroutine panic: %v", r))
			}
		}()

		L2 := cm.pool.Get().(*lua.LState)
		defer func() {
			L2.Close()
			cm.pool.Put(L2)
		}()

		fn2 := L2.NewFunctionFromProto(fn.Proto)
		L2.Push(fn2)

		cm.copyGlobals(L, L2)

		if err := L2.PCall(0, 0, nil); err != nil {
			cm.vm.monitor.handleError(fmt.Errorf("goroutine error: %v", err))
		}
	}()

	return 0
}

func (cm *ConcurrencyModule) copyGlobals(src, dst *lua.LState) {
	globals := []string{
		"print", "sleep", "json_encode", "json_decode",
		"send", "receive", "select", "wait", "import",
		"on_error", "check_memory", "get_goroutines",
		"uuid", "random", "toml", "yaml", "jsonc",
		"text", "crypto", "dotenv", "datetime",
		"csv", "ft", "ini", "tar",
	}

	for _, name := range globals {
		if val := src.GetGlobal(name); val.Type() != lua.LTNil {
			dst.SetGlobal(name, val)
		}
	}
}

func (cm *ConcurrencyModule) createChannel(L *lua.LState) int {
	name := L.CheckString(1)
	bufferSize := L.OptInt(2, 0)

	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.channels[name]; exists {
		L.RaiseError("channel %s already exists", name)
		return 0
	}

	cm.channels[name] = &Channel{
		ch:     make(chan lua.LValue, bufferSize),
		closed: false,
	}
	return 0
}

func (cm *ConcurrencyModule) sendToChannel(L *lua.LState) int {
	name := L.CheckString(1)
	value := L.CheckAny(2)

	cm.mu.RLock()
	ch, exists := cm.channels[name]
	cm.mu.RUnlock()

	if !exists {
		L.RaiseError("channel %s does not exist", name)
		return 0
	}

	ch.mu.RLock()
	if ch.closed {
		ch.mu.RUnlock()
		L.RaiseError("channel %s is closed", name)
		return 0
	}
	ch.mu.RUnlock()

	select {
	case ch.ch <- value:
		L.Push(lua.LTrue)
	case <-time.After(time.Second):
		L.Push(lua.LFalse)
	case <-cm.done:
		L.Push(lua.LFalse)
	}

	return 1
}

func (cm *ConcurrencyModule) receiveFromChannel(L *lua.LState) int {
	name := L.CheckString(1)
	timeout := L.OptNumber(2, 1)

	cm.mu.RLock()
	ch, exists := cm.channels[name]
	cm.mu.RUnlock()

	if !exists {
		L.RaiseError("channel %s does not exist", name)
		return 0
	}

	ch.mu.RLock()
	if ch.closed {
		ch.mu.RUnlock()
		L.Push(lua.LNil)
		return 1
	}
	ch.mu.RUnlock()

	select {
	case value := <-ch.ch:
		L.Push(value)
	case <-time.After(time.Duration(float64(timeout) * float64(time.Second))):
		L.Push(lua.LNil)
	case <-cm.done:
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

	cases := make([]struct {
		ch     *Channel
		name   string
		closed bool
	}, 0, len(channels))

	for _, name := range channels {
		if ch, exists := cm.channels[name]; exists {
			ch.mu.RLock()
			closed := ch.closed
			ch.mu.RUnlock()
			cases = append(cases, struct {
				ch     *Channel
				name   string
				closed bool
			}{ch, name, closed})
		}
	}

	if len(cases) == 0 {
		L.RaiseError("no valid channels provided")
		return 0
	}

	selectCases := make([]reflect.SelectCase, len(cases)+1)
	for i, c := range cases {
		if c.closed {
			selectCases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(nil),
			}
		} else {
			selectCases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c.ch.ch),
			}
		}
	}

	selectCases[len(cases)] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(cm.done),
	}

	chosen, value, ok := reflect.Select(selectCases)
	if chosen == len(cases) {
		L.Push(lua.LNil)
		L.Push(lua.LNil)
		return 2
	}

	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString(cases[chosen].name))
		return 2
	}

	L.Push(value.Interface().(lua.LValue))
	L.Push(lua.LString(cases[chosen].name))
	return 2
}

func (cm *ConcurrencyModule) closeChannel(L *lua.LState) int {
	name := L.CheckString(1)

	cm.mu.Lock()
	ch, exists := cm.channels[name]
	if !exists {
		cm.mu.Unlock()
		L.RaiseError("channel %s does not exist", name)
		return 0
	}

	ch.mu.Lock()
	if ch.closed {
		ch.mu.Unlock()
		cm.mu.Unlock()
		L.RaiseError("channel %s is already closed", name)
		return 0
	}

	ch.closed = true
	close(ch.ch)
	ch.mu.Unlock()

	delete(cm.channels, name)
	cm.mu.Unlock()

	return 0
}

func (cm *ConcurrencyModule) waitForGoroutines(L *lua.LState) int {
	done := make(chan struct{})
	go func() {
		cm.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return 0
	case <-time.After(time.Second * 30):
		L.RaiseError("timeout waiting for goroutines")
		return 0
	}
}

func (cm *ConcurrencyModule) Close() {
	close(cm.done)
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for name, ch := range cm.channels {
		ch.mu.Lock()
		if !ch.closed {
			close(ch.ch)
			ch.closed = true
		}
		ch.mu.Unlock()
		delete(cm.channels, name)
	}
}
