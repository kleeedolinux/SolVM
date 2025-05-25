package modules

import (
	"bytes"

	"github.com/BurntSushi/toml"
	lua "github.com/yuin/gopher-lua"
)

func RegisterTOMLModule(L *lua.LState) {
	tomlModule := L.NewTable()
	L.SetGlobal("toml", tomlModule)

	L.SetField(tomlModule, "encode", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckTable(1)
		data := make(map[string]interface{})
		table.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				data[key.String()] = luaValueToGo(value)
			}
		})
		var buf bytes.Buffer
		err := toml.NewEncoder(&buf).Encode(data)
		if err != nil {
			L.RaiseError("failed to encode TOML: " + err.Error())
			return 0
		}
		L.Push(lua.LString(buf.String()))
		return 1
	}))

	L.SetField(tomlModule, "decode", L.NewFunction(func(L *lua.LState) int {
		tomlStr := L.CheckString(1)
		var data map[string]interface{}
		_, err := toml.Decode(tomlStr, &data)
		if err != nil {
			L.RaiseError("failed to decode TOML: " + err.Error())
			return 0
		}
		L.Push(goValueToLua(L, data))
		return 1
	}))
}

func luaValueToGo(value lua.LValue) interface{} {
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
		result := make(map[string]interface{})
		table.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				result[key.String()] = luaValueToGo(value)
			}
		})
		return result
	default:
		return nil
	}
}

func goValueToLua(L *lua.LState, value interface{}) lua.LValue {
	switch v := value.(type) {
	case nil:
		return lua.LNil
	case bool:
		return lua.LBool(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case map[string]interface{}:
		table := L.NewTable()
		for key, val := range v {
			L.SetField(table, key, goValueToLua(L, val))
		}
		return table
	default:
		return lua.LNil
	}
}
