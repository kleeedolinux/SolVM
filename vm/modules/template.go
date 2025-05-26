package modules

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

func RegisterTemplateModule(L *lua.LState) {
	templateModule := L.NewTable()
	L.SetGlobal("template", templateModule)

	L.SetField(templateModule, "parse", L.NewFunction(func(L *lua.LState) int {
		templateStr := L.CheckString(1)
		tmpl, err := template.New("").Parse(templateStr)
		if err != nil {
			L.RaiseError("failed to parse template: " + err.Error())
			return 0
		}

		L.Push(L.NewFunction(func(L *lua.LState) int {
			if L.GetTop() < 1 {
				L.RaiseError("expected table argument")
				return 0
			}

			data := L.CheckTable(1)
			goData := make(map[string]interface{})
			data.ForEach(func(key, value lua.LValue) {
				if key.Type() == lua.LTString {
					goData[key.String()] = luaValueToGo(value)
				}
			})

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, goData); err != nil {
				L.RaiseError("failed to execute template: " + err.Error())
				return 0
			}

			L.Push(lua.LString(buf.String()))
			return 1
		}))

		return 1
	}))

	L.SetField(templateModule, "parse_file", L.NewFunction(func(L *lua.LState) int {
		filePath := L.CheckString(1)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			L.RaiseError("failed to read template file: " + err.Error())
			return 0
		}

		tmpl, err := template.New(filepath.Base(filePath)).Parse(string(content))
		if err != nil {
			L.RaiseError("failed to parse template: " + err.Error())
			return 0
		}

		L.Push(L.NewFunction(func(L *lua.LState) int {
			if L.GetTop() < 1 {
				L.RaiseError("expected table argument")
				return 0
			}

			data := L.CheckTable(1)
			goData := make(map[string]interface{})
			data.ForEach(func(key, value lua.LValue) {
				if key.Type() == lua.LTString {
					goData[key.String()] = luaValueToGo(value)
				}
			})

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, goData); err != nil {
				L.RaiseError("failed to execute template: " + err.Error())
				return 0
			}

			L.Push(lua.LString(buf.String()))
			return 1
		}))

		return 1
	}))

	L.SetField(templateModule, "parse_files", L.NewFunction(func(L *lua.LState) int {
		patterns := make([]string, L.GetTop())
		for i := 1; i <= L.GetTop(); i++ {
			patterns[i-1] = L.CheckString(i)
		}

		tmpl, err := template.ParseFiles(patterns...)
		if err != nil {
			L.RaiseError("failed to parse template files: " + err.Error())
			return 0
		}

		L.Push(L.NewFunction(func(L *lua.LState) int {
			if L.GetTop() < 1 {
				L.RaiseError("expected table argument")
				return 0
			}

			data := L.CheckTable(1)
			goData := make(map[string]interface{})
			data.ForEach(func(key, value lua.LValue) {
				if key.Type() == lua.LTString {
					goData[key.String()] = luaValueToGo(value)
				}
			})

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, goData); err != nil {
				L.RaiseError("failed to execute template: " + err.Error())
				return 0
			}

			L.Push(lua.LString(buf.String()))
			return 1
		}))

		return 1
	}))

	L.SetField(templateModule, "parse_glob", L.NewFunction(func(L *lua.LState) int {
		pattern := L.CheckString(1)
		tmpl, err := template.ParseGlob(pattern)
		if err != nil {
			L.RaiseError("failed to parse template glob: " + err.Error())
			return 0
		}

		L.Push(L.NewFunction(func(L *lua.LState) int {
			if L.GetTop() < 1 {
				L.RaiseError("expected table argument")
				return 0
			}

			data := L.CheckTable(1)
			goData := make(map[string]interface{})
			data.ForEach(func(key, value lua.LValue) {
				if key.Type() == lua.LTString {
					goData[key.String()] = luaValueToGo(value)
				}
			})

			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, goData); err != nil {
				L.RaiseError("failed to execute template: " + err.Error())
				return 0
			}

			L.Push(lua.LString(buf.String()))
			return 1
		}))

		return 1
	}))
}
