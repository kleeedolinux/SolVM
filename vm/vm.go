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

type ScopeNode struct {
	Parent    *ScopeNode
	Children  []*ScopeNode
	Variables map[string]bool
	Functions map[string]*lua.LFunction
	Level     int
}

type CallNode struct {
	Function *lua.LFunction
	Calls    []*CallNode
	Visited  bool
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
	scopeTree     *ScopeNode
	callGraph     *CallNode
	scopeMu       sync.RWMutex
	functionCache *FunctionCache
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
		scopeTree: &ScopeNode{
			Variables: make(map[string]bool),
			Functions: make(map[string]*lua.LFunction),
			Level:     0,
		},
		functionCache: NewFunctionCache(),
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
	modules.RegisterTemplateModule(vm.state)
	modules.RegisterTableXModule(vm.state)
	modules.RegisterUtilsModule(vm.state)
	modules.RegisterTypesModule(vm.state)

	vm.concMod.Register()
}

func (vm *SolVM) analyzeScope(code string) error {
	vm.scopeMu.Lock()
	defer vm.scopeMu.Unlock()

	L := lua.NewState()
	defer L.Close()

	fn, err := L.LoadString(code)
	if err != nil {
		return err
	}

	vm.analyzeFunction(fn, vm.scopeTree)
	return nil
}

func (vm *SolVM) analyzeFunction(fn *lua.LFunction, parent *ScopeNode) {
	node := &ScopeNode{
		Parent:    parent,
		Children:  make([]*ScopeNode, 0),
		Variables: make(map[string]bool),
		Functions: make(map[string]*lua.LFunction),
		Level:     parent.Level + 1,
	}
	parent.Children = append(parent.Children, node)

	proto := fn.Proto
	for i := 0; i < int(proto.NumUpvalues); i++ {
		name := proto.DbgUpvalues[i]
		node.Variables[name] = true
	}

	for i := 0; i < len(proto.DbgLocals); i++ {
		name := proto.DbgLocals[i].Name
		node.Variables[name] = true
	}

	for i := 0; i < len(proto.Constants); i++ {
		if proto.Constants[i].Type() == lua.LTFunction {
			vm.analyzeFunction(proto.Constants[i].(*lua.LFunction), node)
		}
	}
}

func (vm *SolVM) buildCallGraph() {
	vm.scopeMu.Lock()
	defer vm.scopeMu.Unlock()

	vm.callGraph = &CallNode{
		Function: nil,
		Calls:    make([]*CallNode, 0),
	}

	vm.analyzeCallGraph(vm.scopeTree)
}

func (vm *SolVM) analyzeCallGraph(node *ScopeNode) {
	for _, fn := range node.Functions {
		callNode := &CallNode{
			Function: fn,
			Calls:    make([]*CallNode, 0),
		}
		vm.callGraph.Calls = append(vm.callGraph.Calls, callNode)
		vm.findFunctionCalls(fn, callNode)
	}

	for _, child := range node.Children {
		vm.analyzeCallGraph(child)
	}
}

func (vm *SolVM) findFunctionCalls(fn *lua.LFunction, node *CallNode) {
	if node.Visited {
		return
	}
	node.Visited = true

	proto := fn.Proto
	for i := 0; i < len(proto.Constants); i++ {
		if proto.Constants[i].Type() == lua.LTFunction {
			callNode := &CallNode{
				Function: proto.Constants[i].(*lua.LFunction),
				Calls:    make([]*CallNode, 0),
			}
			node.Calls = append(node.Calls, callNode)
			vm.findFunctionCalls(callNode.Function, callNode)
		}
	}
}

func (vm *SolVM) LoadString(code string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	if vm.memoryLimit > 0 {
		if err := vm.checkMemoryUsage(); err != nil {
			return err
		}
	}

	if err := vm.analyzeScope(code); err != nil {
		return err
	}

	vm.buildCallGraph()

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

func (vm *SolVM) RegisterModule(name string, module Module) {
	vm.moduleMu.Lock()
	defer vm.moduleMu.Unlock()
	vm.modules[name] = module
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
