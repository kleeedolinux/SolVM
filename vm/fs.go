package vm

import (
	"os"

	lua "github.com/yuin/gopher-lua"
)

type FSModule struct {
	vm *SolVM
}

func NewFSModule(vm *SolVM) *FSModule {
	return &FSModule{vm: vm}
}

func (fm *FSModule) Register() {
	fm.vm.RegisterFunction("read_file", fm.readFile)
	fm.vm.RegisterFunction("write_file", fm.writeFile)
	fm.vm.RegisterFunction("list_dir", fm.listDir)
}

func (fm *FSModule) readFile(L *lua.LState) int {
	path := L.CheckString(1)

	data, err := os.ReadFile(path)
	if err != nil {
		L.RaiseError("Failed to read file: %v", err)
		return 0
	}

	L.Push(lua.LString(string(data)))
	return 1
}

func (fm *FSModule) writeFile(L *lua.LState) int {
	path := L.CheckString(1)
	data := L.CheckString(2)

	err := os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		L.RaiseError("Failed to write file: %v", err)
		return 0
	}

	return 0
}

func (fm *FSModule) listDir(L *lua.LState) int {
	path := L.CheckString(1)

	entries, err := os.ReadDir(path)
	if err != nil {
		L.RaiseError("Failed to list directory: %v", err)
		return 0
	}

	result := L.NewTable()
	for i, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileInfo := L.NewTable()
		fileInfo.RawSetString("name", lua.LString(entry.Name()))
		fileInfo.RawSetString("is_dir", lua.LBool(entry.IsDir()))
		fileInfo.RawSetString("size", lua.LNumber(info.Size()))
		fileInfo.RawSetString("mode", lua.LString(info.Mode().String()))
		fileInfo.RawSetString("mod_time", lua.LString(info.ModTime().String()))

		result.RawSetInt(i+1, fileInfo)
	}

	L.Push(result)
	return 1
}
