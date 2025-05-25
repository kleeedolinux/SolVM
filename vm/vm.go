package vm

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"solvm/vm/modules"

	lua "github.com/yuin/gopher-lua"
)

type Module interface {
	Register()
}

type Config struct {
	Timeout       time.Duration
	Debug         bool
	Trace         bool
	MemoryLimit   int64
	MaxGoroutines int
	WorkingDir    string
}

type SolVM struct {
	state         *lua.LState
	mu            sync.RWMutex
	timeout       time.Duration
	ctx           context.Context
	cancel        context.CancelFunc
	errorChan     chan error
	importMod     *ImportModule
	concMod       *ConcurrencyModule
	monitor       *MonitorModule
	httpMod       *HTTPModule
	serverMod     *ServerModule
	fsMod         *FSModule
	schedMod      *SchedulerModule
	netMod        *NetworkModule
	debugMod      *DebugModule
	debug         bool
	trace         bool
	memoryLimit   int64
	maxGoroutines int
	workingDir    string
	modules       map[string]Module
	moduleMu      sync.RWMutex
	startMem      runtime.MemStats
}

func NewSolVM(config Config) *SolVM {
	var ctx context.Context
	var cancel context.CancelFunc

	if config.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), config.Timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	L := lua.NewState()
	vm := &SolVM{
		state:         L,
		timeout:       config.Timeout,
		ctx:           ctx,
		cancel:        cancel,
		errorChan:     make(chan error, 1),
		debug:         config.Debug,
		trace:         config.Trace,
		memoryLimit:   config.MemoryLimit,
		maxGoroutines: config.MaxGoroutines,
		workingDir:    config.WorkingDir,
		modules:       make(map[string]Module),
	}

	runtime.ReadMemStats(&vm.startMem)
	vm.initializeModules()
	vm.registerBuiltinModules()

	return vm
}

func (vm *SolVM) initializeModules() {
	vm.importMod = NewImportModule(vm)
	vm.concMod = NewConcurrencyModule(vm)
	vm.monitor = NewMonitorModule(vm)
	vm.httpMod = NewHTTPModule(vm)
	vm.serverMod = NewServerModule(vm)
	vm.fsMod = NewFSModule(vm)
	vm.schedMod = NewSchedulerModule(vm)
	vm.netMod = NewNetworkModule(vm)
	vm.debugMod = NewDebugModule(vm)
}

func (vm *SolVM) registerBuiltinModules() {
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
	modules.RegisterTARModule(vm.state)
}

func (vm *SolVM) LoadString(code string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	if vm.memoryLimit > 0 {
		if err := vm.checkMemoryUsage(); err != nil {
			return err
		}
	}

	err := vm.state.DoString(code)
	if err != nil {
		vm.monitor.handleError(err)
	}
	return err
}

func (vm *SolVM) ExecuteAsync(code string) error {
	if vm.maxGoroutines > 0 {
		if runtime.NumGoroutine() >= vm.maxGoroutines {
			return fmt.Errorf("maximum number of goroutines reached")
		}
	}

	go func() {
		if err := vm.LoadString(code); err != nil {
			select {
			case vm.errorChan <- err:
			default:
				vm.monitor.handleError(fmt.Errorf("error channel full: %v", err))
			}
			return
		}
		select {
		case vm.errorChan <- nil:
		default:
		}
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

func (vm *SolVM) RegisterModule(name string, module Module) {
	vm.moduleMu.Lock()
	defer vm.moduleMu.Unlock()
	vm.modules[name] = module
	module.Register()
}

func (vm *SolVM) GetModule(name string) (Module, bool) {
	vm.moduleMu.RLock()
	defer vm.moduleMu.RUnlock()
	module, exists := vm.modules[name]
	return module, exists
}

func (vm *SolVM) checkMemoryUsage() error {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	allocated := m.Alloc - vm.startMem.Alloc
	if allocated > uint64(vm.memoryLimit) {
		return fmt.Errorf("memory limit exceeded: %d > %d", allocated, vm.memoryLimit)
	}
	return nil
}

func (vm *SolVM) Close() {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.cancel()
	vm.state.Close()
}
