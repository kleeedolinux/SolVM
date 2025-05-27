package modules

import (
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type Utils struct {
	L *lua.LState
}

func RegisterUtilsModule(L *lua.LState) {
	utils := &Utils{L: L}
	mod := L.NewTable()

	L.SetField(mod, "getfenv", L.NewFunction(utils.getfenv))
	L.SetField(mod, "setfenv", L.NewFunction(utils.setfenv))
	L.SetField(mod, "unpack", L.NewFunction(utils.unpack))
	L.SetField(mod, "split", L.NewFunction(utils.split))
	L.SetField(mod, "join", L.NewFunction(utils.join))
	L.SetField(mod, "escape", L.NewFunction(utils.escape))
	L.SetField(mod, "unescape", L.NewFunction(utils.unescape))

	L.SetGlobal("utils", mod)
}

func (u *Utils) getfenv(L *lua.LState) int {
	level := L.OptInt(1, 1)
	fn := L.GetGlobal("debug").(*lua.LTable).RawGetString("getinfo").(*lua.LFunction)
	L.Push(fn)
	L.Push(lua.LNumber(level))
	L.Push(lua.LString("f"))
	if err := L.PCall(2, 1, nil); err != nil {
		L.RaiseError("error in getfenv: %v", err)
	}
	info := L.Get(-1).(*lua.LTable)
	L.Push(info.RawGetString("func"))
	return 1
}

func (u *Utils) setfenv(L *lua.LState) int {
	fn := L.CheckFunction(1)
	env := L.CheckTable(2)
	L.SetFEnv(fn, env)
	L.Push(fn)
	return 1
}

func (u *Utils) unpack(L *lua.LState) int {
	tbl := L.CheckTable(1)
	start := L.OptInt(2, 1)
	end := L.OptInt(3, tbl.Len())

	for i := start; i <= end; i++ {
		L.Push(tbl.RawGetInt(i))
	}
	return end - start + 1
}

func (u *Utils) split(L *lua.LState) int {
	str := L.CheckString(1)
	sep := L.OptString(2, " ")
	result := L.NewTable()

	parts := strings.Split(str, sep)
	for i, part := range parts {
		L.RawSetInt(result, i+1, lua.LString(part))
	}

	L.Push(result)
	return 1
}

func (u *Utils) join(L *lua.LState) int {
	tbl := L.CheckTable(1)
	sep := L.OptString(2, " ")
	var parts []string

	tbl.ForEach(func(_, value lua.LValue) {
		parts = append(parts, value.String())
	})

	L.Push(lua.LString(strings.Join(parts, sep)))
	return 1
}

func (u *Utils) escape(L *lua.LState) int {
	str := L.CheckString(1)
	result := strings.ReplaceAll(str, "%", "%%")
	L.Push(lua.LString(result))
	return 1
}

func (u *Utils) unescape(L *lua.LState) int {
	str := L.CheckString(1)
	result := strings.ReplaceAll(str, "%%", "%")
	L.Push(lua.LString(result))
	return 1
}
