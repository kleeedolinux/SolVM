package modules

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func RegisterTemplateModule(L *lua.LState) {
	templateModule := L.NewTable()
	L.SetGlobal("template", templateModule)

	L.SetField(templateModule, "parse", L.NewFunction(func(L *lua.LState) int {
		templateStr := L.CheckString(1)
		if strings.TrimSpace(templateStr) == "" {
			L.RaiseError("template string cannot be empty")
			return 0
		}

		tmpl := template.New("")
		tmpl = tmpl.Funcs(template.FuncMap{
			"safe": func(s string) template.HTML {
				return template.HTML(s)
			},
		})

		tmpl, err := tmpl.Parse(templateStr)
		if err != nil {
			L.RaiseError(fmt.Sprintf("failed to parse template: %v", err))
			return 0
		}

		L.Push(L.NewFunction(func(L *lua.LState) int {
			defer func() {
				if err := recover(); err != nil {
					L.RaiseError(fmt.Sprintf("template execution panic: %v", err))
				}
			}()

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
				L.RaiseError(fmt.Sprintf("failed to execute template: %v", err))
				return 0
			}

			L.Push(lua.LString(buf.String()))
			return 1
		}))

		return 1
	}))

	L.SetField(templateModule, "parse_file", L.NewFunction(func(L *lua.LState) int {
		filePath := L.CheckString(1)

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			L.RaiseError(fmt.Sprintf("template file not found: %s", filePath))
			return 0
		}

		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			L.RaiseError(fmt.Sprintf("failed to read template file: %v", err))
			return 0
		}

		if strings.TrimSpace(string(content)) == "" {
			L.RaiseError(fmt.Sprintf("template file %s is empty", filePath))
			return 0
		}

		tmpl := template.New(filepath.Base(filePath))
		tmpl = tmpl.Funcs(template.FuncMap{
			"safe": func(s string) template.HTML {
				return template.HTML(s)
			},
		})

		tmpl, err = tmpl.Parse(string(content))
		if err != nil {
			L.RaiseError(fmt.Sprintf("failed to parse template: %v", err))
			return 0
		}

		L.Push(L.NewFunction(func(L *lua.LState) int {
			defer func() {
				if err := recover(); err != nil {
					L.RaiseError(fmt.Sprintf("template execution panic: %v", err))
				}
			}()

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
				L.RaiseError(fmt.Sprintf("failed to execute template: %v", err))
				return 0
			}

			L.Push(lua.LString(buf.String()))
			return 1
		}))

		return 1
	}))

	L.SetField(templateModule, "parse_files", L.NewFunction(func(L *lua.LState) int {
		if L.GetTop() < 1 {
			L.RaiseError("expected at least one template file path")
			return 0
		}

		patterns := make([]string, L.GetTop())
		for i := 1; i <= L.GetTop(); i++ {
			patterns[i-1] = L.CheckString(i)
			if _, err := os.Stat(patterns[i-1]); os.IsNotExist(err) {
				L.RaiseError(fmt.Sprintf("template file not found: %s", patterns[i-1]))
				return 0
			}
		}

		tmpl := template.New(filepath.Base(patterns[0]))
		tmpl = tmpl.Funcs(template.FuncMap{
			"safe": func(s string) template.HTML {
				return template.HTML(s)
			},
		})

		tmpl, err := tmpl.ParseFiles(patterns...)
		if err != nil {
			L.RaiseError(fmt.Sprintf("failed to parse template files: %v", err))
			return 0
		}

		L.Push(L.NewFunction(func(L *lua.LState) int {
			defer func() {
				if err := recover(); err != nil {
					L.RaiseError(fmt.Sprintf("template execution panic: %v", err))
				}
			}()

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
				L.RaiseError(fmt.Sprintf("failed to execute template: %v", err))
				return 0
			}

			L.Push(lua.LString(buf.String()))
			return 1
		}))

		return 1
	}))

	L.SetField(templateModule, "parse_glob", L.NewFunction(func(L *lua.LState) int {
		pattern := L.CheckString(1)
		if strings.TrimSpace(pattern) == "" {
			L.RaiseError("glob pattern cannot be empty")
			return 0
		}

		files, err := filepath.Glob(pattern)
		if err != nil {
			L.RaiseError(fmt.Sprintf("invalid glob pattern: %v", err))
			return 0
		}

		if len(files) == 0 {
			L.RaiseError(fmt.Sprintf("no files found matching pattern: %s", pattern))
			return 0
		}

		tmpl := template.New(filepath.Base(files[0]))
		tmpl = tmpl.Funcs(template.FuncMap{
			"safe": func(s string) template.HTML {
				return template.HTML(s)
			},
		})

		tmpl, err = tmpl.ParseFiles(files...)
		if err != nil {
			L.RaiseError(fmt.Sprintf("failed to parse template files: %v", err))
			return 0
		}

		L.Push(L.NewFunction(func(L *lua.LState) int {
			defer func() {
				if err := recover(); err != nil {
					L.RaiseError(fmt.Sprintf("template execution panic: %v", err))
				}
			}()

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
				L.RaiseError(fmt.Sprintf("failed to execute template: %v", err))
				return 0
			}

			L.Push(lua.LString(buf.String()))
			return 1
		}))

		return 1
	}))
}
