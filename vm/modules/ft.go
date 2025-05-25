package modules

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

func RegisterFTModule(L *lua.LState) {
	ftModule := L.NewTable()
	L.SetGlobal("ft", ftModule)

	L.SetField(ftModule, "download", L.NewFunction(func(L *lua.LState) int {
		url := L.CheckString(1)
		path := L.CheckString(2)

		resp, err := http.Get(url)
		if err != nil {
			L.RaiseError("failed to download file: " + err.Error())
			return 0
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			L.RaiseError("failed to download file: " + resp.Status)
			return 0
		}

		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			L.RaiseError("failed to create directory: " + err.Error())
			return 0
		}

		file, err := os.Create(path)
		if err != nil {
			L.RaiseError("failed to create file: " + err.Error())
			return 0
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			L.RaiseError("failed to save file: " + err.Error())
			return 0
		}

		return 0
	}))

	L.SetField(ftModule, "upload", L.NewFunction(func(L *lua.LState) int {
		path := L.CheckString(1)
		url := L.CheckString(2)

		file, err := os.Open(path)
		if err != nil {
			L.RaiseError("failed to open file: " + err.Error())
			return 0
		}
		defer file.Close()

		resp, err := http.Post(url, "application/octet-stream", file)
		if err != nil {
			L.RaiseError("failed to upload file: " + err.Error())
			return 0
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			L.RaiseError("failed to upload file: " + resp.Status)
			return 0
		}

		return 0
	}))

	L.SetField(ftModule, "copy", L.NewFunction(func(L *lua.LState) int {
		src := L.CheckString(1)
		dst := L.CheckString(2)

		srcFile, err := os.Open(src)
		if err != nil {
			L.RaiseError("failed to open source file: " + err.Error())
			return 0
		}
		defer srcFile.Close()

		dir := filepath.Dir(dst)
		if err := os.MkdirAll(dir, 0755); err != nil {
			L.RaiseError("failed to create directory: " + err.Error())
			return 0
		}

		dstFile, err := os.Create(dst)
		if err != nil {
			L.RaiseError("failed to create destination file: " + err.Error())
			return 0
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			L.RaiseError("failed to copy file: " + err.Error())
			return 0
		}

		return 0
	}))

	L.SetField(ftModule, "move", L.NewFunction(func(L *lua.LState) int {
		src := L.CheckString(1)
		dst := L.CheckString(2)

		dir := filepath.Dir(dst)
		if err := os.MkdirAll(dir, 0755); err != nil {
			L.RaiseError("failed to create directory: " + err.Error())
			return 0
		}

		if err := os.Rename(src, dst); err != nil {
			L.RaiseError("failed to move file: " + err.Error())
			return 0
		}

		return 0
	}))
}
