package modules

import (
	lua "github.com/yuin/gopher-lua"
)

type Types struct {
	L *lua.LState
}

func RegisterTypesModule(L *lua.LState) {
	types := &Types{L: L}
	mod := L.NewTable()

	L.SetField(mod, "is_callable", L.NewFunction(types.isCallable))
	L.SetField(mod, "is_integer", L.NewFunction(types.isInteger))
	L.SetField(mod, "is_number", L.NewFunction(types.isNumber))
	L.SetField(mod, "is_string", L.NewFunction(types.isString))
	L.SetField(mod, "is_table", L.NewFunction(types.isTable))
	L.SetField(mod, "is_function", L.NewFunction(types.isFunction))
	L.SetField(mod, "is_boolean", L.NewFunction(types.isBoolean))
	L.SetField(mod, "is_nil", L.NewFunction(types.isNil))
	L.SetField(mod, "type", L.NewFunction(types.typeOf))

	L.SetGlobal("types", mod)
}

func (t *Types) isCallable(L *lua.LState) int {
	val := L.Get(1)
	L.Push(lua.LBool(val.Type() == lua.LTFunction || (val.Type() == lua.LTTable && L.GetMetatable(val) != lua.LNil)))
	return 1
}

func (t *Types) isInteger(L *lua.LState) int {
	val := L.Get(1)
	if val.Type() != lua.LTNumber {
		L.Push(lua.LFalse)
		return 1
	}
	num := val.(lua.LNumber)
	L.Push(lua.LBool(float64(num) == float64(int64(num))))
	return 1
}

func (t *Types) isNumber(L *lua.LState) int {
	val := L.Get(1)
	L.Push(lua.LBool(val.Type() == lua.LTNumber))
	return 1
}

func (t *Types) isString(L *lua.LState) int {
	val := L.Get(1)
	L.Push(lua.LBool(val.Type() == lua.LTString))
	return 1
}

func (t *Types) isTable(L *lua.LState) int {
	val := L.Get(1)
	L.Push(lua.LBool(val.Type() == lua.LTTable))
	return 1
}

func (t *Types) isFunction(L *lua.LState) int {
	val := L.Get(1)
	L.Push(lua.LBool(val.Type() == lua.LTFunction))
	return 1
}

func (t *Types) isBoolean(L *lua.LState) int {
	val := L.Get(1)
	L.Push(lua.LBool(val.Type() == lua.LTBool))
	return 1
}

func (t *Types) isNil(L *lua.LState) int {
	val := L.Get(1)
	L.Push(lua.LBool(val.Type() == lua.LTNil))
	return 1
}

func (t *Types) typeOf(L *lua.LState) int {
	val := L.Get(1)
	L.Push(lua.LString(val.Type().String()))
	return 1
}
