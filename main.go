package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"solvm/vm"
)

const VERSION = "1.2.0"
const COPYRIGHT = "SolVM (c) 2025"

func printUsage() {
	fmt.Printf("%s - A Lua Virtual Machine with Enhanced Features\n", COPYRIGHT)
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
	fmt.Println("\nRunning without arguments starts the SolVM console")
}

func runConsole(config vm.Config) {
	fmt.Printf("%s v%s\n", COPYRIGHT, VERSION)
	fmt.Println("Type 'exit' or 'quit' to exit")
	fmt.Println("Type 'help' for available commands")
	fmt.Println()

	vm := vm.NewSolVM(config)
	defer vm.Close()
	vm.RegisterCustomFunctions()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("solvm> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		switch line {
		case "exit", "quit":
			return
		case "help":
			fmt.Println("Available commands:")
			fmt.Println("  exit, quit  - Exit the console")
			fmt.Println("  help       - Show this help message")
			fmt.Println("  clear      - Clear the screen")
			fmt.Println("  version    - Show version information")
			fmt.Println("  Any other input is treated as Lua code")
		case "clear":
			fmt.Print("\033[H\033[2J")
		case "version":
			fmt.Printf("%s v%s\n", COPYRIGHT, VERSION)
		default:
			if err := vm.LoadString(line); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
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
		fmt.Printf("%s v%s\n", COPYRIGHT, VERSION)
		return
	}

	config := vm.Config{
		Timeout:       *timeout,
		Debug:         *debug,
		Trace:         *trace,
		MemoryLimit:   int64(*memoryLimit) * 1024 * 1024,
		MaxGoroutines: *maxGoroutines,
	}

	if flag.NArg() == 0 {
		config.WorkingDir, _ = os.Getwd()
		runConsole(config)
		return
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

	config.WorkingDir = filepath.Dir(absPath)

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
