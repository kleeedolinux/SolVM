package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kleeedolinux/solvm/vm"
)

func main() {
	timeout := flag.Duration("timeout", 5*time.Second, "Execution timeout")
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: solvm [options] <lua-file>")
		os.Exit(1)
	}

	file := flag.Arg(0)
	code, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	isServer := strings.Contains(string(code), "create_server") ||
		strings.Contains(string(code), "start_server")

	if isServer {
		*timeout = 0
	}

	vm := vm.NewSolVM(*timeout)
	defer vm.Close()

	vm.RegisterCustomFunctions()

	if err := vm.LoadString(string(code)); err != nil {
		fmt.Printf("Error executing Lua code: %v\n", err)
		os.Exit(1)
	}
}
