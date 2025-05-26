package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"solvm/vm"
)

const VERSION = "1.3.0"
const COPYRIGHT = "SolVM (c) 2025"
const GITHUB_API_URL = "https://api.github.com/repos/kleeedolinux/SolVM/releases/latest"

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func checkForUpdates() {
	resp, err := http.Get(GITHUB_API_URL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var release Release
	if err := json.Unmarshal(body, &release); err != nil {
		return
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	if latestVersion != VERSION {
		fmt.Printf("New version available: %s (current: %s)\n", latestVersion, VERSION)
		fmt.Println("Run 'solvm update' to update to the latest version")
	}
}

func updateSolVM() error {
	var installerURL string
	var installerName string

	switch runtime.GOOS {
	case "windows":
		installerName = "solvm-installer-windows-amd64.exe"
	case "linux":
		if runtime.GOARCH == "arm64" {
			installerName = "solvm-installer-linux-arm64"
		} else {
			installerName = "solvm-installer-linux-amd64"
		}
	case "darwin":
		if runtime.GOARCH == "arm64" {
			installerName = "solvm-installer-darwin-arm64"
		} else {
			installerName = "solvm-installer-darwin-amd64"
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	installerURL = fmt.Sprintf("https://github.com/kleeedolinux/SolVM-installer/releases/download/v1.0.0/%s", installerName)

	resp, err := http.Get(installerURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	installerPath := filepath.Join(os.TempDir(), installerName)
	file, err := os.Create(installerPath)
	if err != nil {
		return err
	}
	defer os.Remove(installerPath)

	_, err = io.Copy(file, resp.Body)
	file.Close()
	if err != nil {
		return err
	}

	if runtime.GOOS != "windows" {
		if err := os.Chmod(installerPath, 0755); err != nil {
			return err
		}
	}

	cmd := exec.Command(installerPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

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
	fmt.Println("  -update             Update to the latest version")
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
	update := flag.Bool("update", false, "Update to the latest version")

	flag.Usage = printUsage
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s v%s\n", COPYRIGHT, VERSION)
		return
	}

	if *update {
		if err := updateSolVM(); err != nil {
			fmt.Printf("Error updating SolVM: %v\n", err)
			os.Exit(1)
		}
		return
	}

	checkForUpdates()

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
