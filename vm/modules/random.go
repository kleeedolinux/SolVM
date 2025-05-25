package modules

import (
	"crypto/rand"
	"encoding/binary"
	"math"

	lua "github.com/yuin/gopher-lua"
)

func RegisterRandomModule(L *lua.LState) {
	randomModule := L.NewTable()
	L.SetGlobal("random", randomModule)

	L.SetField(randomModule, "number", L.NewFunction(func(L *lua.LState) int {
		var b [8]byte
		_, err := rand.Read(b[:])
		if err != nil {
			L.RaiseError("failed to generate random number: " + err.Error())
			return 0
		}
		L.Push(lua.LNumber(math.Float64frombits(binary.LittleEndian.Uint64(b[:]))))
		return 1
	}))

	L.SetField(randomModule, "int", L.NewFunction(func(L *lua.LState) int {
		min := L.CheckInt(1)
		max := L.CheckInt(2)
		if min >= max {
			L.RaiseError("min must be less than max")
			return 0
		}
		var b [8]byte
		_, err := rand.Read(b[:])
		if err != nil {
			L.RaiseError("failed to generate random number: " + err.Error())
			return 0
		}
		randNum := int(math.Float64frombits(binary.LittleEndian.Uint64(b[:])))
		result := min + (randNum % (max - min + 1))
		L.Push(lua.LNumber(result))
		return 1
	}))

	L.SetField(randomModule, "string", L.NewFunction(func(L *lua.LState) int {
		length := L.CheckInt(1)
		if length <= 0 {
			L.RaiseError("length must be positive")
			return 0
		}
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		result := make([]byte, length)
		for i := range result {
			var b [1]byte
			_, err := rand.Read(b[:])
			if err != nil {
				L.RaiseError("failed to generate random string: " + err.Error())
				return 0
			}
			result[i] = charset[int(b[0])%len(charset)]
		}
		L.Push(lua.LString(string(result)))
		return 1
	}))
}
