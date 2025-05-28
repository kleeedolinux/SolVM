package vm

import (
	"encoding/json"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type FunctionTrie struct {
	children map[rune]*FunctionTrie
	function lua.LGFunction
	id       int
}

type FunctionCache struct {
	mu          sync.RWMutex
	trie        *FunctionTrie
	jumpTable   map[int]lua.LGFunction
	nextID      int
	inlineCache map[string]*lua.LFunction
}

func NewFunctionCache() *FunctionCache {
	return &FunctionCache{
		trie: &FunctionTrie{
			children: make(map[rune]*FunctionTrie),
		},
		jumpTable:   make(map[int]lua.LGFunction),
		inlineCache: make(map[string]*lua.LFunction),
	}
}

func (fc *FunctionCache) Register(name string, fn lua.LGFunction) int {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	current := fc.trie
	for _, char := range name {
		if _, exists := current.children[char]; !exists {
			current.children[char] = &FunctionTrie{
				children: make(map[rune]*FunctionTrie),
			}
		}
		current = current.children[char]
	}

	if current.function == nil {
		current.id = fc.nextID
		fc.jumpTable[fc.nextID] = fn
		fc.nextID++
	}
	current.function = fn
	return current.id
}

func (fc *FunctionCache) Lookup(name string) (lua.LGFunction, int) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	current := fc.trie
	for _, char := range name {
		if next, exists := current.children[char]; exists {
			current = next
		} else {
			return nil, -1
		}
	}
	return current.function, current.id
}

func (fc *FunctionCache) GetByID(id int) lua.LGFunction {
	fc.mu.RLock()
	defer fc.mu.RUnlock()
	return fc.jumpTable[id]
}

func (fc *FunctionCache) CacheFunction(name string, fn *lua.LFunction) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.inlineCache[name] = fn
}

func (fc *FunctionCache) GetCachedFunction(name string) *lua.LFunction {
	fc.mu.RLock()
	defer fc.mu.RUnlock()
	return fc.inlineCache[name]
}

func (vm *SolVM) RegisterCustomFunctions() {
	vm.RegisterFunction("json_encode", jsonEncode)
	vm.RegisterFunction("json_decode", jsonDecode)
	vm.RegisterFunction("sleep", sleep)
	vm.importMod.Register()
	vm.concMod.Register()
	vm.monitor.Register()
	vm.httpMod.Register()
	vm.serverMod.Register()
	vm.fsMod.Register()
	vm.schedMod.Register()
	vm.netMod.Register()
	vm.debugMod.Register()
}

func (vm *SolVM) RegisterFunction(name string, fn lua.LGFunction) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	id := vm.functionCache.Register(name, fn)
	vm.state.SetGlobal(name, vm.state.NewFunction(func(L *lua.LState) int {
		if cached := vm.functionCache.GetCachedFunction(name); cached != nil {
			L.Push(cached)
			return 1
		}
		if cachedFn := vm.functionCache.GetByID(id); cachedFn != nil {
			return cachedFn(L)
		}
		return fn(L)
	}))
}

func jsonEncode(L *lua.LState) int {
	value := L.CheckAny(1)

	goValue := convertToGoValue(value)
	data, err := json.Marshal(goValue)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(string(data)))
	return 1
}

func jsonDecode(L *lua.LState) int {
	str := L.CheckString(1)

	var result interface{}
	if err := json.Unmarshal([]byte(str), &result); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(convertToLuaValue(L, result))
	return 1
}

func sleep(L *lua.LState) int {
	duration := float64(L.CheckNumber(1))
	time.Sleep(time.Duration(duration * float64(time.Second)))
	return 0
}

func convertToGoValue(value lua.LValue) interface{} {
	switch value.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return bool(value.(lua.LBool))
	case lua.LTNumber:
		return float64(value.(lua.LNumber))
	case lua.LTString:
		return string(value.(lua.LString))
	case lua.LTTable:
		table := value.(*lua.LTable)

		isArray := true
		maxIndex := 0
		table.ForEach(func(key, _ lua.LValue) {
			if key.Type() != lua.LTNumber {
				isArray = false
			} else {
				idx := int(key.(lua.LNumber))
				if idx > maxIndex {
					maxIndex = idx
				}
			}
		})

		if isArray {
			arr := make([]interface{}, maxIndex)
			table.ForEach(func(key, value lua.LValue) {
				idx := int(key.(lua.LNumber)) - 1
				arr[idx] = convertToGoValue(value)
			})
			return arr
		}

		obj := make(map[string]interface{})
		table.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				obj[key.String()] = convertToGoValue(value)
			}
		})
		return obj
	default:
		return nil
	}
}

func convertToLuaValue(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []interface{}:
		arr := L.NewTable()
		for i, item := range v {
			arr.RawSetInt(i+1, convertToLuaValue(L, item))
		}
		return arr
	case map[string]interface{}:
		tbl := L.NewTable()
		for k, v := range v {
			tbl.RawSetString(k, convertToLuaValue(L, v))
		}
		return tbl
	default:
		return lua.LNil
	}
}
