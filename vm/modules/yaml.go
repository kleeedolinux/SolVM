package modules

import (
	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v3"
)

func RegisterYAMLModule(L *lua.LState) {
	yamlModule := L.NewTable()
	L.SetGlobal("yaml", yamlModule)

	L.SetField(yamlModule, "encode", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckTable(1)
		data := make(map[string]interface{})
		table.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				data[key.String()] = luaValueToGo(value)
			}
		})
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			L.RaiseError("failed to encode YAML: " + err.Error())
			return 0
		}
		L.Push(lua.LString(string(yamlData)))
		return 1
	}))

	L.SetField(yamlModule, "decode", L.NewFunction(func(L *lua.LState) int {
		yamlStr := L.CheckString(1)
		var data map[string]interface{}
		err := yaml.Unmarshal([]byte(yamlStr), &data)
		if err != nil {
			L.RaiseError("failed to decode YAML: " + err.Error())
			return 0
		}
		L.Push(goValueToLua(L, data))
		return 1
	}))
}
