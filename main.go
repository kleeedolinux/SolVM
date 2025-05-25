package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"solvm/vm"
)

func printUsage() {
	fmt.Println("SolVM - A Lua Virtual Machine with Enhanced Features")
	fmt.Println("\nUsage: solvm [options] <lua-file>")
	fmt.Println("\nOptions:")
	fmt.Println("  -timeout duration    Execution timeout (default 5s)")
	fmt.Println("  -debug              Enable debug mode")
	fmt.Println("  -trace              Enable trace mode")
	fmt.Println("  -memory-limit int   Memory limit in MB (default 1024)")
	fmt.Println("  -max-goroutines int Maximum number of goroutines (default 1000)")
	fmt.Println("  -version            Show version information")
	fmt.Println("\nExamples:")
	fmt.Println("  solvm script.lua")
	fmt.Println("  solvm -timeout 10s -debug script.lua")
	fmt.Println("  solvm -memory-limit 2048 server.lua")
}

func main() {
	timeout := flag.Duration("timeout", 5*time.Second, "Execution timeout")
	debug := flag.Bool("debug", false, "Enable debug mode")
	trace := flag.Bool("trace", false, "Enable trace mode")
	memoryLimit := flag.Int("memory-limit", 1024, "Memory limit in MB")
	maxGoroutines := flag.Int("max-goroutines", 1000, "Maximum number of goroutines")
	showVersion := flag.Bool("version", false, "Show version information")

	flag.Usage = printUsage
	flag.Parse()

	if *showVersion {
		fmt.Println("SolVM version 1.1.0")
		return
	}

	if flag.NArg() == 0 {
		printUsage()
		os.Exit(1)
	}

	file := flag.Arg(0)
	if !strings.HasSuffix(file, ".lua") {
		fmt.Printf("Error: File must have .lua extension: %s\n", file)
		os.Exit(1)
	}

	absPath, err := filepath.Abs(file)
	if err != nil {
		fmt.Printf("Error resolving file path: %v\n", err)
		os.Exit(1)
	}

	code, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	isServer := strings.Contains(string(code), "create_server") ||
		strings.Contains(string(code), "start_server")

	if isServer {
		*timeout = 0
		fmt.Println("Server mode detected: timeout disabled")
	}

	config := vm.Config{
		Timeout:       *timeout,
		Debug:         *debug,
		Trace:         *trace,
		MemoryLimit:   int64(*memoryLimit) * 1024 * 1024,
		MaxGoroutines: *maxGoroutines,
		WorkingDir:    filepath.Dir(absPath),
	}

	vm := vm.NewSolVM(config)
	defer vm.Close()

	if *debug {
		fmt.Printf("Debug mode enabled\n")
		fmt.Printf("Memory limit: %d MB\n", *memoryLimit)
		fmt.Printf("Max goroutines: %d\n", *maxGoroutines)
		fmt.Printf("Working directory: %s\n", config.WorkingDir)
	}

	vm.RegisterCustomFunctions()

	if err := vm.LoadString(string(code)); err != nil {
		fmt.Printf("Error executing Lua code: %v\n", err)
		os.Exit(1)
	}
}
