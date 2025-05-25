package modules

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func RegisterTARModule(L *lua.LState) {
	tarModule := L.NewTable()
	L.SetGlobal("tar", tarModule)

	L.SetField(tarModule, "create", L.NewFunction(func(L *lua.LState) int {
		archivePath := L.CheckString(1)
		sourcePath := L.CheckString(2)
		compress := L.OptBool(3, false)

		file, err := os.Create(archivePath)
		if err != nil {
			L.RaiseError("failed to create archive: " + err.Error())
			return 0
		}
		defer file.Close()

		var writer io.Writer = file
		if compress {
			gzipWriter := gzip.NewWriter(file)
			defer gzipWriter.Close()
			writer = gzipWriter
		}

		tarWriter := tar.NewWriter(writer)
		defer tarWriter.Close()

		err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(sourcePath, path)
			if err != nil {
				return err
			}
			header.Name = relPath

			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}

			if !info.IsDir() {
				file, err := os.Open(path)
				if err != nil {
					return err
				}
				defer file.Close()

				_, err = io.Copy(tarWriter, file)
				if err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			L.RaiseError("failed to create archive: " + err.Error())
			return 0
		}

		return 0
	}))

	L.SetField(tarModule, "extract", L.NewFunction(func(L *lua.LState) int {
		archivePath := L.CheckString(1)
		destPath := L.CheckString(2)

		file, err := os.Open(archivePath)
		if err != nil {
			L.RaiseError("failed to open archive: " + err.Error())
			return 0
		}
		defer file.Close()

		var reader io.Reader = file
		if strings.HasSuffix(archivePath, ".gz") || strings.HasSuffix(archivePath, ".tgz") {
			gzipReader, err := gzip.NewReader(file)
			if err != nil {
				L.RaiseError("failed to read gzip archive: " + err.Error())
				return 0
			}
			defer gzipReader.Close()
			reader = gzipReader
		}

		tarReader := tar.NewReader(reader)

		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				L.RaiseError("failed to read archive: " + err.Error())
				return 0
			}

			target := filepath.Join(destPath, header.Name)

			switch header.Typeflag {
			case tar.TypeDir:
				if err := os.MkdirAll(target, 0755); err != nil {
					L.RaiseError("failed to create directory: " + err.Error())
					return 0
				}
			case tar.TypeReg:
				dir := filepath.Dir(target)
				if err := os.MkdirAll(dir, 0755); err != nil {
					L.RaiseError("failed to create directory: " + err.Error())
					return 0
				}

				file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
				if err != nil {
					L.RaiseError("failed to create file: " + err.Error())
					return 0
				}

				if _, err := io.Copy(file, tarReader); err != nil {
					file.Close()
					L.RaiseError("failed to extract file: " + err.Error())
					return 0
				}
				file.Close()
			}
		}

		return 0
	}))

	L.SetField(tarModule, "list", L.NewFunction(func(L *lua.LState) int {
		archivePath := L.CheckString(1)

		file, err := os.Open(archivePath)
		if err != nil {
			L.RaiseError("failed to open archive: " + err.Error())
			return 0
		}
		defer file.Close()

		var reader io.Reader = file
		if strings.HasSuffix(archivePath, ".gz") || strings.HasSuffix(archivePath, ".tgz") {
			gzipReader, err := gzip.NewReader(file)
			if err != nil {
				L.RaiseError("failed to read gzip archive: " + err.Error())
				return 0
			}
			defer gzipReader.Close()
			reader = gzipReader
		}

		tarReader := tar.NewReader(reader)
		result := L.NewTable()

		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				L.RaiseError("failed to read archive: " + err.Error())
				return 0
			}

			fileInfo := L.NewTable()
			L.SetField(fileInfo, "name", lua.LString(header.Name))
			L.SetField(fileInfo, "size", lua.LNumber(header.Size))
			L.SetField(fileInfo, "mode", lua.LNumber(header.Mode))
			L.SetField(fileInfo, "type", lua.LString(string(header.Typeflag)))
			L.SetField(fileInfo, "mod_time", lua.LNumber(header.ModTime.Unix()))

			L.RawSetInt(result, result.Len()+1, fileInfo)
		}

		L.Push(result)
		return 1
	}))
}
