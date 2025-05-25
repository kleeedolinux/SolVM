package modules

import (
	"encoding/csv"
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func RegisterCSVModule(L *lua.LState) {
	csvModule := L.NewTable()
	L.SetGlobal("csv", csvModule)

	L.SetField(csvModule, "read", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		file, err := os.Open(path)
		if err != nil {
			L.RaiseError("failed to open CSV file: " + err.Error())
			return 0
		}
		defer file.Close()

		reader := csv.NewReader(file)
		records, err := reader.ReadAll()
		if err != nil {
			L.RaiseError("failed to read CSV: " + err.Error())
			return 0
		}

		result := L.NewTable()
		for i, record := range records {
			row := L.NewTable()
			for j, field := range record {
				L.RawSetInt(row, j+1, lua.LString(field))
			}
			L.RawSetInt(result, i+1, row)
		}

		L.Push(result)
		return 1
	}))

	L.SetField(csvModule, "write", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		table := L.CheckTable(2)

		file, err := os.Create(path)
		if err != nil {
			L.RaiseError("failed to create CSV file: " + err.Error())
			return 0
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		table.ForEach(func(key, value lua.LValue) {
			if row, ok := value.(*lua.LTable); ok {
				record := make([]string, 0)
				row.ForEach(func(_, field lua.LValue) {
					record = append(record, field.String())
				})
				writer.Write(record)
			}
		})

		if err := writer.Error(); err != nil {
			L.RaiseError("failed to write CSV: " + err.Error())
			return 0
		}

		return 0
	}))

	L.SetField(csvModule, "parse", L.NewFunction(func(L *lua.LState) int {
		data := L.CheckString(1)
		reader := csv.NewReader(strings.NewReader(data))
		records, err := reader.ReadAll()
		if err != nil {
			L.RaiseError("failed to parse CSV: " + err.Error())
			return 0
		}

		result := L.NewTable()
		for i, record := range records {
			row := L.NewTable()
			for j, field := range record {
				L.RawSetInt(row, j+1, lua.LString(field))
			}
			L.RawSetInt(result, i+1, row)
		}

		L.Push(result)
		return 1
	}))

	L.SetField(csvModule, "stringify", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckTable(1)
		var buffer strings.Builder
		writer := csv.NewWriter(&buffer)

		table.ForEach(func(key, value lua.LValue) {
			if row, ok := value.(*lua.LTable); ok {
				record := make([]string, 0)
				row.ForEach(func(_, field lua.LValue) {
					record = append(record, field.String())
				})
				writer.Write(record)
			}
		})

		writer.Flush()
		if err := writer.Error(); err != nil {
			L.RaiseError("failed to stringify CSV: " + err.Error())
			return 0
		}

		L.Push(lua.LString(buffer.String()))
		return 1
	}))
}
