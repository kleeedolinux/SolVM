package modules

import (
	"github.com/google/uuid"
	lua "github.com/yuin/gopher-lua"
)

func RegisterUUIDModule(L *lua.LState) {
	uuidModule := L.NewTable()
	L.SetGlobal("uuid", uuidModule)

	L.SetField(uuidModule, "v4", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LString(uuid.New().String()))
		return 1
	}))

	L.SetField(uuidModule, "v4_without_hyphens", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LString(uuid.New().String()))
		return 1
	}))

	L.SetField(uuidModule, "is_valid", L.NewFunction(func(L *lua.LState) int {
		uuidStr := L.CheckString(1)
		_, err := uuid.Parse(uuidStr)
		L.Push(lua.LBool(err == nil))
		return 1
	}))
}
