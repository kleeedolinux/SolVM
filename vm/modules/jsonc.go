package modules

import (
	"encoding/json"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func RegisterJSONCModule(L *lua.LState) {
	jsoncModule := L.NewTable()
	L.SetGlobal("jsonc", jsoncModule)

	L.SetField(jsoncModule, "encode", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckTable(1)
		data := make(map[string]interface{})
		table.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				data[key.String()] = luaValueToGo(value)
			}
		})
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			L.RaiseError("failed to encode JSONC: " + err.Error())
			return 0
		}
		L.Push(lua.LString(string(jsonData)))
		return 1
	}))

	L.SetField(jsoncModule, "decode", L.NewFunction(func(L *lua.LState) int {
		jsoncStr := L.CheckString(1)
		jsonStr := removeComments(jsoncStr)
		var data map[string]interface{}
		err := json.Unmarshal([]byte(jsonStr), &data)
		if err != nil {
			L.RaiseError("failed to decode JSONC: " + err.Error())
			return 0
		}
		L.Push(goValueToLua(L, data))
		return 1
	}))
}

func removeComments(input string) string {
	lines := strings.Split(input, "\n")
	var result []string
	inString := false
	escapeNext := false

	for _, line := range lines {
		var newLine strings.Builder
		for i := 0; i < len(line); i++ {
			if escapeNext {
				newLine.WriteByte(line[i])
				escapeNext = false
				continue
			}

			switch line[i] {
			case '\\':
				escapeNext = true
				newLine.WriteByte(line[i])
			case '"':
				inString = !inString
				newLine.WriteByte(line[i])
			case '/':
				if !inString && i+1 < len(line) {
					if line[i+1] == '/' {
						break
					} else if line[i+1] == '*' {
						i++
						for i+1 < len(line) {
							if line[i] == '*' && line[i+1] == '/' {
								i++
								break
							}
							i++
						}
						continue
					}
				}
				newLine.WriteByte(line[i])
			default:
				newLine.WriteByte(line[i])
			}
		}
		result = append(result, newLine.String())
	}
	return strings.Join(result, "\n")
}
