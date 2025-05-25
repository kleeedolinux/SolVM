package modules

import (
	"bufio"
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func RegisterINIModule(L *lua.LState) {
	iniModule := L.NewTable()
	L.SetGlobal("ini", iniModule)

	L.SetField(iniModule, "read", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		file, err := os.Open(path)
		if err != nil {
			L.RaiseError("failed to open INI file: " + err.Error())
			return 0
		}
		defer file.Close()

		result := L.NewTable()
		var currentSection *lua.LTable

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
				continue
			}

			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				sectionName := line[1 : len(line)-1]
				currentSection = L.NewTable()
				L.SetField(result, sectionName, currentSection)
			} else if currentSection != nil {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					L.SetField(currentSection, key, lua.LString(value))
				}
			}
		}

		if err := scanner.Err(); err != nil {
			L.RaiseError("failed to read INI file: " + err.Error())
			return 0
		}

		L.Push(result)
		return 1
	}))

	L.SetField(iniModule, "write", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		table := L.CheckTable(2)

		file, err := os.Create(path)
		if err != nil {
			L.RaiseError("failed to create INI file: " + err.Error())
			return 0
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		defer writer.Flush()

		table.ForEach(func(key, value lua.LValue) {
			if section, ok := value.(*lua.LTable); ok {
				writer.WriteString("[" + key.String() + "]\n")
				section.ForEach(func(sectionKey, sectionValue lua.LValue) {
					writer.WriteString(sectionKey.String() + "=" + sectionValue.String() + "\n")
				})
				writer.WriteString("\n")
			}
		})

		return 0
	}))

	L.SetField(iniModule, "parse", L.NewFunction(func(L *lua.LState) int {
		data := L.CheckString(1)
		result := L.NewTable()
		var currentSection *lua.LTable

		scanner := bufio.NewScanner(strings.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
				continue
			}

			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				sectionName := line[1 : len(line)-1]
				currentSection = L.NewTable()
				L.SetField(result, sectionName, currentSection)
			} else if currentSection != nil {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					L.SetField(currentSection, key, lua.LString(value))
				}
			}
		}

		L.Push(result)
		return 1
	}))

	L.SetField(iniModule, "stringify", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckTable(1)
		var buffer strings.Builder

		table.ForEach(func(key, value lua.LValue) {
			if section, ok := value.(*lua.LTable); ok {
				buffer.WriteString("[" + key.String() + "]\n")
				section.ForEach(func(sectionKey, sectionValue lua.LValue) {
					buffer.WriteString(sectionKey.String() + "=" + sectionValue.String() + "\n")
				})
				buffer.WriteString("\n")
			}
		})

		L.Push(lua.LString(buffer.String()))
		return 1
	}))
}
