package modules

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

func RegisterDatetimeModule(L *lua.LState) {
	datetimeModule := L.NewTable()
	L.SetGlobal("datetime", datetimeModule)

	L.SetField(datetimeModule, "now", L.NewFunction(func(L *lua.LState) int {
		now := time.Now()
		L.Push(lua.LNumber(now.Unix()))
		return 1
	}))

	L.SetField(datetimeModule, "format", L.NewFunction(func(L *lua.LState) int {
		timestamp := L.CheckNumber(1)
		format := L.OptString(2, time.RFC3339)
		t := time.Unix(int64(timestamp), 0)
		L.Push(lua.LString(t.Format(format)))
		return 1
	}))

	L.SetField(datetimeModule, "parse", L.NewFunction(func(L *lua.LState) int {
		timeStr := L.CheckString(1)
		format := L.OptString(2, time.RFC3339)
		t, err := time.Parse(format, timeStr)
		if err != nil {
			L.RaiseError("failed to parse time: " + err.Error())
			return 0
		}
		L.Push(lua.LNumber(t.Unix()))
		return 1
	}))

	L.SetField(datetimeModule, "add", L.NewFunction(func(L *lua.LState) int {
		timestamp := L.CheckNumber(1)
		duration := L.CheckString(2)
		d, err := time.ParseDuration(duration)
		if err != nil {
			L.RaiseError("invalid duration: " + err.Error())
			return 0
		}
		t := time.Unix(int64(timestamp), 0).Add(d)
		L.Push(lua.LNumber(t.Unix()))
		return 1
	}))

	L.SetField(datetimeModule, "diff", L.NewFunction(func(L *lua.LState) int {
		timestamp1 := L.CheckNumber(1)
		timestamp2 := L.CheckNumber(2)
		t1 := time.Unix(int64(timestamp1), 0)
		t2 := time.Unix(int64(timestamp2), 0)
		diff := t2.Sub(t1)
		L.Push(lua.LString(diff.String()))
		return 1
	}))

	L.SetField(datetimeModule, "sleep", L.NewFunction(func(L *lua.LState) int {
		duration := L.CheckString(1)
		d, err := time.ParseDuration(duration)
		if err != nil {
			L.RaiseError("invalid duration: " + err.Error())
			return 0
		}
		time.Sleep(d)
		return 0
	}))
}
