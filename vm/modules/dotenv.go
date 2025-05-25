package modules

import (
	"bufio"
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func RegisterDotenvModule(L *lua.LState) {
	dotenvModule := L.NewTable()
	L.SetGlobal("dotenv", dotenvModule)

	L.SetField(dotenvModule, "load", L.NewFunction(func(L *lua.LState) int {
		path := L.OptString(1, ".env")
		file, err := os.Open(path)
		if err != nil {
			L.RaiseError("failed to open .env file: " + err.Error())
			return 0
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}

		if err := scanner.Err(); err != nil {
			L.RaiseError("failed to read .env file: " + err.Error())
			return 0
		}

		return 0
	}))

	L.SetField(dotenvModule, "get", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(1)
		defaultValue := L.OptString(2, "")
		value := os.Getenv(key)
		if value == "" {
			value = defaultValue
		}
		L.Push(lua.LString(value))
		return 1
	}))

	L.SetField(dotenvModule, "set", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(1)
		value := L.CheckString(2)
		os.Setenv(key, value)
		return 0
	}))
}
