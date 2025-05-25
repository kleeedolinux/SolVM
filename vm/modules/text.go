package modules

import (
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func RegisterTextModule(L *lua.LState) {
	textModule := L.NewTable()
	L.SetGlobal("text", textModule)

	L.SetField(textModule, "trim", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		L.Push(lua.LString(strings.TrimSpace(str)))
		return 1
	}))

	L.SetField(textModule, "lower", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		L.Push(lua.LString(strings.ToLower(str)))
		return 1
	}))

	L.SetField(textModule, "upper", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		L.Push(lua.LString(strings.ToUpper(str)))
		return 1
	}))

	L.SetField(textModule, "title", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		L.Push(lua.LString(strings.Title(strings.ToLower(str))))
		return 1
	}))

	L.SetField(textModule, "split", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		sep := L.CheckString(2)
		parts := strings.Split(str, sep)
		table := L.NewTable()
		for i, part := range parts {
			L.RawSetInt(table, i+1, lua.LString(part))
		}
		L.Push(table)
		return 1
	}))

	L.SetField(textModule, "join", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckTable(1)
		sep := L.CheckString(2)
		var parts []string
		table.ForEach(func(_, value lua.LValue) {
			parts = append(parts, value.String())
		})
		L.Push(lua.LString(strings.Join(parts, sep)))
		return 1
	}))

	L.SetField(textModule, "replace", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		old := L.CheckString(2)
		new := L.CheckString(3)
		L.Push(lua.LString(strings.Replace(str, old, new, -1)))
		return 1
	}))

	L.SetField(textModule, "contains", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		substr := L.CheckString(2)
		L.Push(lua.LBool(strings.Contains(str, substr)))
		return 1
	}))

	L.SetField(textModule, "starts_with", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		prefix := L.CheckString(2)
		L.Push(lua.LBool(strings.HasPrefix(str, prefix)))
		return 1
	}))

	L.SetField(textModule, "ends_with", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		suffix := L.CheckString(2)
		L.Push(lua.LBool(strings.HasSuffix(str, suffix)))
		return 1
	}))

	L.SetField(textModule, "pad_left", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		length := L.CheckInt(2)
		pad := L.CheckString(3)
		if len(pad) == 0 {
			L.RaiseError("pad string cannot be empty")
			return 0
		}
		for len(str) < length {
			str = pad + str
		}
		L.Push(lua.LString(str))
		return 1
	}))

	L.SetField(textModule, "repeat_str", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		count := L.CheckInt(2)
		if count < 0 {
			L.RaiseError("count cannot be negative")
			return 0
		}
		L.Push(lua.LString(strings.Repeat(str, count)))
		return 1
	}))
}
