package vm

import (
	"encoding/json"
	"time"

	lua "github.com/yuin/gopher-lua"
)

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
