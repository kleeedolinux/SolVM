package vm

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	lua "github.com/yuin/gopher-lua"
)

type SchedulerModule struct {
	vm        *SolVM
	intervals map[int]*time.Ticker
	timeouts  map[int]*time.Timer
	crons     map[int]cron.EntryID
	cron      *cron.Cron
	mu        sync.RWMutex
	nextID    int
}

func NewSchedulerModule(vm *SolVM) *SchedulerModule {
	return &SchedulerModule{
		vm:        vm,
		intervals: make(map[int]*time.Ticker),
		timeouts:  make(map[int]*time.Timer),
		crons:     make(map[int]cron.EntryID),
		cron:      cron.New(cron.WithSeconds()),
		nextID:    1,
	}
}

func (sm *SchedulerModule) Register() {
	sm.vm.RegisterFunction("set_interval", sm.setInterval)
	sm.vm.RegisterFunction("set_timeout", sm.setTimeout)
	sm.vm.RegisterFunction("cron", sm.setCron)
	sm.cron.Start()
}

func (sm *SchedulerModule) setInterval(L *lua.LState) int {
	fn := L.CheckFunction(1)
	seconds := float64(L.CheckNumber(2))

	sm.mu.Lock()
	id := sm.nextID
	sm.nextID++
	sm.mu.Unlock()

	ticker := time.NewTicker(time.Duration(seconds * float64(time.Second)))
	sm.intervals[id] = ticker

	go func() {
		for range ticker.C {
			L2 := lua.NewState()
			defer L2.Close()

			fn2 := L2.NewFunctionFromProto(fn.Proto)
			L2.Push(fn2)
			if err := L2.PCall(0, 0, nil); err != nil {
				sm.vm.monitor.handleError(fmt.Errorf("Interval function error: %v", err))
			}
		}
	}()

	L.Push(lua.LNumber(id))
	return 1
}

func (sm *SchedulerModule) setTimeout(L *lua.LState) int {
	fn := L.CheckFunction(1)
	seconds := float64(L.CheckNumber(2))

	sm.mu.Lock()
	id := sm.nextID
	sm.nextID++
	sm.mu.Unlock()

	timer := time.NewTimer(time.Duration(seconds * float64(time.Second)))
	sm.timeouts[id] = timer

	go func() {
		<-timer.C
		L2 := lua.NewState()
		defer L2.Close()

		fn2 := L2.NewFunctionFromProto(fn.Proto)
		L2.Push(fn2)
		if err := L2.PCall(0, 0, nil); err != nil {
			sm.vm.monitor.handleError(fmt.Errorf("Timeout function error: %v", err))
		}

		sm.mu.Lock()
		delete(sm.timeouts, id)
		sm.mu.Unlock()
	}()

	L.Push(lua.LNumber(id))
	return 1
}

func (sm *SchedulerModule) setCron(L *lua.LState) int {
	schedule := L.CheckString(1)
	fn := L.CheckFunction(2)

	sm.mu.Lock()
	id := sm.nextID
	sm.nextID++
	sm.mu.Unlock()

	entryID, err := sm.cron.AddFunc(schedule, func() {
		L2 := lua.NewState()
		defer L2.Close()

		fn2 := L2.NewFunctionFromProto(fn.Proto)
		L2.Push(fn2)
		if err := L2.PCall(0, 0, nil); err != nil {
			sm.vm.monitor.handleError(fmt.Errorf("Cron function error: %v", err))
		}
	})

	if err != nil {
		L.RaiseError("Invalid cron schedule: %v", err)
		return 0
	}

	sm.crons[id] = entryID
	L.Push(lua.LNumber(id))
	return 1
}

func (sm *SchedulerModule) ClearInterval(id int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if ticker, exists := sm.intervals[id]; exists {
		ticker.Stop()
		delete(sm.intervals, id)
	}
}

func (sm *SchedulerModule) ClearTimeout(id int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if timer, exists := sm.timeouts[id]; exists {
		timer.Stop()
		delete(sm.timeouts, id)
	}
}

func (sm *SchedulerModule) ClearCron(id int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if entryID, exists := sm.crons[id]; exists {
		sm.cron.Remove(entryID)
		delete(sm.crons, id)
	}
}
