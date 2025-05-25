package vm

import (
	"context"
	"sync"
	"time"

	"solvm/vm/modules"

	lua "github.com/yuin/gopher-lua"
)

type SolVM struct {
	state     *lua.LState
	mu        sync.RWMutex
	timeout   time.Duration
	ctx       context.Context
	cancel    context.CancelFunc
	errorChan chan error
	importMod *ImportModule
	concMod   *ConcurrencyModule
	monitor   *MonitorModule
	httpMod   *HTTPModule
	serverMod *ServerModule
	fsMod     *FSModule
	schedMod  *SchedulerModule
	netMod    *NetworkModule
	debugMod  *DebugModule
}

func NewSolVM(timeout time.Duration) *SolVM {
	var ctx context.Context
	var cancel context.CancelFunc

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	vm := &SolVM{
		state:     lua.NewState(),
		timeout:   timeout,
		ctx:       ctx,
		cancel:    cancel,
		errorChan: make(chan error, 1),
	}
	vm.importMod = NewImportModule(vm)
	vm.concMod = NewConcurrencyModule(vm)
	vm.monitor = NewMonitorModule(vm)
	vm.httpMod = NewHTTPModule(vm)
	vm.serverMod = NewServerModule(vm)
	vm.fsMod = NewFSModule(vm)
	vm.schedMod = NewSchedulerModule(vm)
	vm.netMod = NewNetworkModule(vm)
	vm.debugMod = NewDebugModule(vm)

	modules.RegisterUUIDModule(vm.state)
	modules.RegisterRandomModule(vm.state)
	modules.RegisterTOMLModule(vm.state)
	modules.RegisterYAMLModule(vm.state)
	modules.RegisterJSONCModule(vm.state)
	modules.RegisterTextModule(vm.state)
	modules.RegisterCryptoModule(vm.state)
	modules.RegisterDotenvModule(vm.state)
	modules.RegisterDatetimeModule(vm.state)
	modules.RegisterCSVModule(vm.state)
	modules.RegisterFTModule(vm.state)
	modules.RegisterINIModule(vm.state)

	return vm
}

func (vm *SolVM) LoadString(code string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	err := vm.state.DoString(code)
	if err != nil {
		vm.monitor.handleError(err)
	}
	return err
}

func (vm *SolVM) ExecuteAsync(code string) error {
	go func() {
		if err := vm.LoadString(code); err != nil {
			vm.errorChan <- err
			return
		}
		vm.errorChan <- nil
	}()

	select {
	case err := <-vm.errorChan:
		return err
	case <-vm.ctx.Done():
		if vm.timeout > 0 {
			err := vm.ctx.Err()
			vm.monitor.handleError(err)
			return err
		}
		return nil
	}
}

func (vm *SolVM) RegisterFunction(name string, fn lua.LGFunction) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.state.SetGlobal(name, vm.state.NewFunction(fn))
}

func (vm *SolVM) Close() {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.cancel()
	vm.state.Close()
}
