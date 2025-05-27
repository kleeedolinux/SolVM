### SolVM: Beyond Speed – Crafting a Modern, Productive Lua Runtime

When you think about Lua, especially LuaJIT, raw execution speed often comes to mind. And it's true, LuaJIT can be blisteringly fast, sometimes even nipping at the heels of C performance for certain tasks. But SolVM isn't trying to win that specific race. Its philosophy is different. While performance is always good, SolVM prioritizes **developer productivity, effortless portability, and rich extensibility** for building modern applications with Lua.

The foundational decision that shapes SolVM's entire architecture is its construction in Go. This wasn't an arbitrary choice; it brings a constellation of benefits crucial for a contemporary runtime:

*   **Effortless Distribution:** One of Go's most lauded features is its ability to compile into a single, self-contained binary. This means SolVM can be built, and the result is *one file* that runs on Windows, Linux, or macOS (including ARM64 variants) without needing users to install a separate Lua interpreter, C compilers, or a web of dependencies. This is a game-changer for distributing CLI tools, automation scripts, or even small server applications built with SolVM. Just drop the binary, and it works.
*   **Native Integration with Modern Infrastructure:** The Go standard library is a treasure trove for today's development needs. It provides robust, battle-tested packages for HTTP clients and servers, file system interactions, JSON encoding/decoding, cryptography, and much more. SolVM leverages this by creating Lua-friendly wrappers around these Go functionalities. So, when your Lua script makes an HTTP request or parses a JSON file, it's actually Go code doing the heavy lifting, providing reliability and performance without requiring complex C bindings from the Lua side.
*   **High-Level Concurrency, Simplified:** Go's approach to concurrency with goroutines (lightweight, concurrently executing functions) and channels (for communication between goroutines) is widely acclaimed for its simplicity and power. SolVM exposes this model directly to Lua. This means you can write Lua scripts that perform multiple tasks "simultaneously" – like handling multiple network requests, running background tasks, or processing data in parallel – with a much more intuitive and less error-prone paradigm than traditional threading models often found in other Lua environments.
*   **Seamlessly Tapping into the Go Ecosystem:** The Go community has produced a vast number of high-quality open-source libraries. If SolVM needs to incorporate a new feature, say, support for a new data format or a specific cloud service API, and there's a Go library for it, integrating that library into SolVM is typically far more straightforward than wrestling with Lua's C API or FFI (Foreign Function Interface) to bridge with a C/C++ library. This agility allows SolVM to evolve and adapt to new technological landscapes more rapidly.

So, while LuaJIT might edge out SolVM in raw micro-benchmarks focusing on pure Lua loop execution or numerical computation, SolVM's strengths lie elsewhere. It's designed to be a runtime for *modern scripting needs*: building CLI tools that feel native, automating complex workflows, creating small to medium-sized web services, and generally making Lua a more versatile and less friction-filled tool in a developer's arsenal. The goal isn't to replace LuaJIT but to offer a compelling, hackable, and accessible alternative, built entirely in Go, making it inherently portable and easier for developers to understand, modify, and contribute to.

### The Curtain Rises: SolVM's Entry Point (`main.go`)

The journey of a SolVM execution begins in `main.go`. This Go program serves as the command-line interface (CLI) and the initial orchestrator for the entire runtime.

When a user types `solvm` in their terminal, perhaps followed by options and a script name, this `main.go` application springs to life. Its first order of business is to interpret the user's intentions. This is handled using Go's standard `flag` package. It meticulously parses command-line arguments to understand settings like an **execution timeout** (`-timeout`, defaulting to a sensible 5 seconds for typical scripts to prevent runaways), whether to enable **debug mode** (`-debug`) for more verbose output, or **trace mode** (`-trace`). Users can also specify resource constraints like the **memory limit** (`-memory-limit`, defaulting to 1024MB) and the **maximum number of goroutines** (`-max-goroutines`, default 1000) that SolVM is allowed to use. This configurability is key to tailoring the runtime environment to the specific needs of the script or application being run.

SolVM also incorporates a user-friendly **update mechanism**. The `checkForUpdates` function makes a non-blocking HTTP GET request to the SolVM GitHub repository's release API. It fetches information about the latest release, specifically the `tag_name`, and compares it against the `VERSION` constant embedded in the SolVM binary at compile time. If a newer version is available, SolVM politely informs the user and suggests running `solvm update`. Should the user invoke this command, the `updateSolVM` function takes over. It intelligently determines the user's operating system (`runtime.GOOS` yields "windows", "linux", "darwin", etc.) and CPU architecture (`runtime.GOARCH` gives "amd64", "arm64"). With this information, it constructs the correct download URL for the platform-specific SolVM installer (e.g., `solvm-installer-linux-arm64`) from a dedicated installer release page on GitHub. The installer is then downloaded to a temporary location, made executable (on POSIX-like systems), and finally executed, allowing SolVM to seamlessly update itself.

Standard CLI courtesies are also present. Invoking `solvm -version` will display the copyright notice and the current SolVM version. If the command-line arguments are incorrect or if the user simply needs guidance, the `printUsage` function displays a comprehensive help message, detailing all available options and providing illustrative examples of how to use SolVM.

### The Interactive Realm: The SolVM Console (`runConsole`)

If SolVM is invoked without specifying a Lua script file, it gracefully transitions into its interactive console mode, often referred to as a REPL (Read-Eval-Print Loop). The `runConsole` function in `main.go` orchestrates this experience. Upon entry, it presents a welcoming banner displaying the SolVM version and copyright. It then enters a loop, prompting the user with `solvm> `.

Inside this loop, a `SolVM` instance (the core runtime object from the `vm` package, which we'll explore in detail) is created, configured with the default or user-specified settings. Crucially, `vm.RegisterCustomFunctions()` is called. This step injects all of SolVM's extended functionalities – its custom modules and built-in helper functions – into the Lua environment that the console will use.

The console then uses Go's `bufio.NewScanner` to read lines of input from the user. It's designed for ease of use:
*   Basic utility commands like `exit` or `quit` terminate the console session.
*   `help` displays a list of available console commands.
*   `clear` uses ANSI escape sequences (`\033[H\033[2J`) to clear the terminal screen, providing a cleaner workspace.
*   `version` re-displays the SolVM version information.
Any input that doesn't match these predefined commands is treated as Lua code. This code string is then passed to the `vm.LoadString(line)` method. The `SolVM` instance executes this Lua snippet, and if any errors occur during parsing or execution, they are caught and printed to the console, allowing for immediate feedback. This interactive mode is invaluable for experimenting with Lua, testing small code fragments, or quickly trying out SolVM's specific features.

### The Main Act: Executing a Lua Script

When a user provides a Lua script file (e.g., `solvm myscript.lua`), SolVM embarks on its primary mission: executing that script within its enriched environment.

The `main` function first validates that the provided file indeed has a `.lua` extension, a simple but effective sanity check. It then determines the **absolute path** of the script file using `filepath.Abs()`. The directory containing the script (`filepath.Dir(absPath)`) is then set as the `config.WorkingDir`. This working directory is significant because it serves as the base for resolving any relative paths that the Lua script might use when trying to `import` other local modules or access local files.

The content of the Lua script is then read into a byte slice using `os.ReadFile(absPath)`. A particularly insightful step follows: SolVM performs a quick, heuristic scan of the script's content, looking for keywords like `"create_server"` or `"start_server"`. If these are found, SolVM intelligently infers that the script is likely intended to run as a long-lived server process. In such cases, it automatically disables the execution `timeout` (by setting it to `0`), as servers typically need to run indefinitely, listening for incoming requests, rather than terminating after a fixed duration.

With the configuration finalized and the script code in hand, a `SolVM` instance is created: `vm := vm.NewSolVM(config)`. This `vm` object, defined in `vm/vm.go`, encapsulates the gopher-lua state (`LState`) and serves as the central hub for all SolVM's extended runtime capabilities. A `defer vm.Close()` statement ensures that resources held by the VM (like the Lua state and any background tasks) are cleaned up when the `main` function exits.

If `debug` mode was enabled via command-line flags, SolVM prints out the active configuration details, such as the memory limit, maximum goroutines, and the working directory, providing transparency into the runtime environment.

The cornerstone of SolVM's power is then invoked: `vm.RegisterCustomFunctions()`. This method meticulously registers all the custom functions and modules that SolVM provides, making them available to the Lua script. This includes foundational utilities like `json_encode` and `json_decode` (wrapping Go's JSON capabilities), a `sleep` function (using `time.Sleep`), and the specialized modules like `importMod` for module loading, `concMod` for concurrency, `monitor` for error handling and resource checking, `httpMod` for client-side HTTP, `serverMod` for building HTTP/WebSocket servers, `fsMod` for file system operations, `schedMod` for timed tasks, `netMod` for low-level networking, and `debugMod` for debugging utilities. Additionally, a suite of specific data format and utility modules located in `vm/modules/` (like `crypto`, `text`, `uuid`, `toml`, `yaml`, `csv`, `tar`, `template`, `tablex`, `types`, `utils`, etc.) are also loaded into the Lua global namespace or made accessible via `import()`.

Finally, the Lua script code itself is executed by calling `vm.LoadString(string(code))`. The `SolVM` instance processes this string, and the gopher-lua engine parses and runs the Lua bytecode. If any error occurs during this execution (be it a syntax error in Lua or a runtime error, perhaps from a misbehaving SolVM built-in function or an unhandled error in the Lua code itself), the error is captured. SolVM's `monitor` module's error handling mechanism is invoked, and the error is printed to the console, after which SolVM typically exits with a non-zero status code to indicate failure. If the script completes successfully, SolVM exits cleanly.

### Delving into the `vm` Package: The Heart of SolVM

The `vm` directory contains the Go code that truly defines SolVM's enhanced runtime. Let's explore its key components:

**`vm.go`: The Central `SolVM` Struct**
This file defines the `SolVM` struct, which is the main object representing a Lua runtime instance. It holds:
*   `state`: The actual `*lua.LState` from the gopher-lua library. This is the Lua virtual machine.
*   `timeout`: The execution timeout duration.
*   `ctx`, `cancel`: Go's `context.Context` and its cancel function, used to manage the execution lifetime and enforce timeouts.
*   `errorChan`: A channel used internally for asynchronous operations to report errors.
*   `debug`, `trace`, `memoryLimit`, `maxGoroutines`, `workingDir`: These store the configuration values.
*   Specialized module instances: `importMod`, `concMod`, `monitor`, `httpMod`, `serverMod`, `fsMod`, `schedMod`, `netMod`, `debugMod`. These are structs that encapsulate the logic for each major feature set.
*   `modules`: A map to keep track of dynamically registered modules.
*   `startMem`: Stores initial memory statistics for calculating usage.

The `NewSolVM(config Config)` constructor initializes the Lua state, sets up the context for timeout management, and critically, it creates instances of all its internal module handlers (like `NewImportModule`, `NewConcurrencyModule`, etc.). It then calls `vm.registerBuiltinModules()`, which is responsible for making the Go-powered functionalities from the `vm/modules/` subdirectory (like `crypto`, `toml`, `yaml`, `text`, `uuid`, etc.) available to the Lua environment by setting them as global tables or functions in the `LState`. The `RegisterCustomFunctions()` method further populates the Lua environment with top-level utility functions (e.g., `json_encode`, `sleep`) and calls the `Register()` method on each of its specialized module handlers (e.g., `concMod.Register()`, `httpMod.Register()`).

The `LoadString(code string)` method is the primary way to execute Lua code. It's mutex-protected (`vm.mu.Lock()`) to ensure thread safety if SolVM were to be used in more complex embedding scenarios (though the CLI primarily uses it serially for the main script). Before execution, if a `memoryLimit` is set, `checkMemoryUsage()` is called. This function reads current Go runtime memory statistics (`runtime.ReadMemStats()`) and compares the allocated memory (`m.Alloc`) against the `startMem.Alloc` (recorded when the `SolVM` instance was created) and the `memoryLimit`. If the limit is exceeded, it returns an error, preventing the script from consuming excessive memory. The actual Lua execution is done via `vm.state.DoString(code)`.

The `ExecuteAsync` method is designed for scenarios where Lua code might need to run without blocking the main Go thread, wrapping `LoadString` in a goroutine and using `errorChan` and the `context` for completion or timeout signaling.

**`concurrency.go`: Goroutines and Channels for Lua**
This is where SolVM brings Go-style concurrency to Lua.
*   **`Channel` struct:** Represents a buffered or unbuffered channel, holding a Go channel (`chan lua.LValue`) and a flag for whether it's closed.
*   **`ConcurrencyModule` struct:** Manages all channels (`channels map[string]*Channel`), uses a `sync.WaitGroup` (`wg`) to keep track of active goroutines (for the `wait()` function), and employs a `sync.Pool` (`pool`) for `lua.LState` objects. Reusing `LState`s from a pool is an optimization to reduce the overhead of creating new Lua states for every goroutine, as state creation can be somewhat expensive.
*   **`goFunc(L *lua.LState) int`:** This is the Lua-callable `go()` function. When a Lua function is passed to `go()`, `goFunc` increments the `WaitGroup`, checks against `maxGoroutines`, and then launches a new Go goroutine. Inside this Go goroutine:
    1.  It defers `wg.Done()` to signal completion.
    2.  It includes a `recover()` to catch panics within the goroutine and report them via `vm.monitor.handleError()`.
    3.  It gets an `LState` from the `pool` (or creates a new one if the pool is empty).
    4.  Crucially, it calls `cm.copyGlobals(L, L2)`. This function iterates through a predefined list of essential SolVM global functions (like `print`, `sleep`, `send`, `receive`, module names like `json`, `crypto`, etc.) and copies them from the parent Lua state (`L`) to the new goroutine's Lua state (`L2`). This ensures that the Lua code running in the goroutine has access to the same SolVM built-in functionalities.
    5.  The Lua function is then executed in this new state (`L2.PCall(0, 0, nil)`).
    6.  After execution (or error), the `LState` is closed and returned to the `pool`.
*   **`createChannel`, `sendToChannel`, `receiveFromChannel`, `closeChannel`:** These functions manage the lifecycle and operations on named channels, using mutexes (`cm.mu`) for thread-safe access to the `channels` map. `send` and `receive` use Go's `select` statement with a timeout and a check against `cm.done` (a channel closed when SolVM shuts down) to prevent indefinite blocking.
*   **`selectChannel`:** Implements Lua's `select(...)` functionality. It takes multiple channel names, builds a slice of `reflect.SelectCase` based on these channels, and uses Go's `reflect.Select` to wait for one of them to become ready. This allows Lua code to multiplex over several channels efficiently.
*   **`waitForGoroutines`:** This Lua-callable `wait()` simply calls `cm.wg.Wait()`, blocking until all goroutines launched via `go()` have completed. A timeout is included to prevent indefinite hangs.

**`import.go`: Sophisticated Module Loading**
The `ImportModule` is responsible for Lua's `import()` function. It's significantly more advanced than Lua's default `require`.
*   **Caching:** It maintains an in-memory cache (`cache map[string]*ModuleCache`) for module code. This `ModuleCache` stores the code, timestamp, and size. When a module is imported, SolVM first checks this cache. If found and not stale (though staleness check isn't explicitly shown in this snippet, it's a common pattern for caches), the cached code is used, avoiding disk I/O or network requests. The cache has a `maxCacheSize` and an LRU-like eviction policy (`evictOldestCache`) to manage its size.
*   **Loading Strategies:**
    *   **Local Files:** If the path ends in `.lua`, it tries to load it directly or from a configured `modulesDir`.
    *   **Folders:** If the path ends with `/`, `importFolder` is called. It reads all `.lua` files in that directory (relative to `modulesDir`), executes each, and makes their returned tables available under a namespace derived from the folder and file names.
    *   **ZIP Archives:** If the path ends in `.zip`, `importFromZip` handles it.
        *   If the path is a URL (`isURL` checks for `http://` or `https://`), `downloadZip` fetches the ZIP file using the `httpClient` (configured with timeouts and connection pooling).
        *   It then uses Go's `archive/zip` to read the contents. `.lua` files within the ZIP are executed, and their results are typically namespaced.
    *   **GitHub Repositories:** If the path looks like a GitHub URL (e.g., `github.com/owner/repo`), `importFromGitHub` is invoked.
        *   `parseGitHubURL` extracts the owner and repository name.
        *   `getGitHubDownloadURL` attempts to find the latest release's ZIP download URL from the GitHub API. If no release is found, it defaults to downloading the `main` branch as a ZIP.
        *   The ZIP is then downloaded and processed by `importFromZip`.
*   **Security/Resource Limits:** When loading modules (especially from URLs or files), `maxModuleSize` is enforced to prevent excessively large files from being loaded.
*   **`metadata(tbl)` function:** When a module calls `metadata({...})`, this function in `import.go` captures that table. It finds the currently executing module's name (derived from `debug.getinfo().Source`) and stores the provided metadata table associated with that module, typically within `package.loaded[moduleName]`. This allows metadata like version, author, dependencies, etc., to be programmatically accessible later.
*   **Preventing Re-import:** A `loaded map[string]bool` tracks already imported modules to avoid redundant execution.

**`http.go` & `server.go`: Web Capabilities**
*   **`HTTPModule` (`http.go`):** Provides Lua functions like `http_get`, `http_post`, `http_put`, `http_delete`, and a generic `http_request`. These functions use a shared Go `http.Client` (configured with a timeout) to make the actual HTTP requests. Responses (status code, headers, body) are converted into Lua tables for the script to use. Error handling is piped through `vm.monitor.handleError`.
*   **`ServerModule` (`server.go`):** Allows Lua scripts to create and manage web servers.
    *   `createServer(serverID, port, isHTTPS, [certFile, keyFile])`: Creates an `http.Server` instance in Go, configured for HTTP or HTTPS (loading TLS certificates if specified). These servers are stored in a map (`sm.servers`) keyed by `serverID`.
    *   `startServer(serverID)`: Starts the specified server in a new Go goroutine (`server.ListenAndServe()` or `server.ListenAndServeTLS()`).
    *   `stopServer(serverID)`: Gracefully shuts down a server.
    *   `handleHTTP(serverID, path, handlerFunc)`: Registers a Lua function (`handlerFunc`) to handle HTTP requests for a given `path` on a specific server. When a request comes in, SolVM converts the Go `http.Request` details (method, path, query, headers) into a Lua table and calls the Lua handler. The Lua handler is expected to return a Lua table representing the response (status, headers, body), which SolVM then uses to construct the HTTP response.
    *   `handleWebSocket(serverID, path, handlerFunc)`: Similar to `handleHTTP`, but for WebSocket connections. It uses the `gorilla/websocket` library to upgrade HTTP connections to WebSockets. The Lua handler function receives a Lua table representing the WebSocket connection, with methods like `send(message)` and `receive()` (which internally call `conn.WriteMessage` and `conn.ReadMessage` on the Go WebSocket connection object).

**Other Core `vm` Components:**
*   **`fs.go` (`FSModule`):** Provides `read_file`, `write_file`, and `list_dir`. These are thin wrappers around Go's `os` package functions (`os.ReadFile`, `os.WriteFile`, `os.ReadDir`), making file system interaction straightforward from Lua.
*   **`scheduler.go` (`SchedulerModule`):** Enables timed and scheduled execution of Lua functions.
    *   `set_interval(func, seconds)`: Uses `time.NewTicker` in Go to repeatedly call the Lua function.
    *   `set_timeout(func, seconds)`: Uses `time.NewTimer` to call the Lua function once after a delay.
    *   `cron(schedule_string, func)`: Uses the `robfig/cron/v3` Go library to schedule Lua functions based on cron expressions (e.g., `"0 * * * *"` for hourly execution).
    It manages these timers and cron jobs, allowing them to be cleared, and uses a `sync.Pool` for `LState`s for the callbacks.
*   **`network.go` (`NetworkModule`):** Handles lower-level networking beyond HTTP.
    *   `tcp_listen(port)` and `tcp_connect(host, port)`: Create TCP listeners and client connections using Go's `net` package. Accepted/created connections are represented as Lua tables with `read`, `write`, and `close` methods that map to the underlying Go connection operations.
    *   `udp_sendto(addr, port, message)` and `udp_recvfrom(port)`: Provide UDP send and receive capabilities. `udp_recvfrom` returns a Lua table with `receive` and `close` methods.
    *   `resolve_dns(hostname)`: Uses `net.LookupIP` to perform DNS resolution.
*   **`debug.go` (`DebugModule`):**
    *   `watch_file(filePath, callbackFunc)`: Monitors a file for changes using a `time.Ticker` to periodically check `os.Stat().ModTime()`. If a change is detected, the Lua `callbackFunc` is executed in a new Lua state. This is the backbone of hot reloading.
    *   `reload_script()`: Intended to be called from a `watch_file` callback. It would re-read the main script file (path likely stored in a global like `_SCRIPT_PATH`) and execute it in a fresh Lua state (though the snippet shows `L2.DoString`, implying the *current* state might be what's intended for re-execution, or a new state is created and then the old one potentially discarded by the caller).
    *   `trace()`: Uses Lua's `debug.traceback` to get and print the current call stack.
*   **`monitor.go` (`MonitorModule`):**
    *   `on_error(handlerFunc)`: Allows Lua scripts to register global error handler functions. When `vm.monitor.handleError(err)` is called (from anywhere in SolVM, including panics in goroutines), it iterates through these registered Lua handlers and calls them with the error message.
    *   `check_memory()`: Provides detailed memory usage statistics (alloc diff, total alloc diff, system memory, GC count, number of goroutines) by using `runtime.ReadMemStats()`.
    *   `get_goroutines()`: Intended to return a list/map of active (SolVM-managed) goroutines, possibly by tracking them in `goroutineMap` when they are created via `go()`.

**Specialized Modules in `vm/modules/`:**
These files typically define a `Register<ModuleName>Module(L *lua.LState)` function. This function creates a new Lua table, populates it with functions that wrap Go logic, and then sets this table as a global in the Lua state (e.g., `L.SetGlobal("crypto", cryptoModule)`).

*   **`crypto.go`:** Implements functions like `crypto.md5()`, `crypto.sha256()`, `crypto.base64_encode()`, `crypto.aes_encrypt()`, `crypto.rsa_generate()`, etc., by calling corresponding functions from Go's `crypto/*` packages. For encryption, it often handles padding (like PKCS7) as well.
*   **Data Format Modules (`csv.go`, `ini.go`, `jsonc.go`, `toml.go`, `yaml.go`):**
    Each of these provides `encode` and `decode` (or `read`/`write`, `parse`/`stringify`) functions for their respective formats.
    *   `encode`: Takes a Lua table, converts it to a Go `map[string]interface{}` or `[]interface{}` (via helper functions like `luaValueToGo`), and then uses the relevant Go library (e.g., `encoding/json`, `github.com/BurntSushi/toml`, `gopkg.in/yaml.v3`, `encoding/csv`) to serialize it into a string.
    *   `decode`: Takes a string in the format, uses the Go library to unmarshal it into a Go `map[string]interface{}`, and then converts this back into a Lua table (via `goValueToLua`).
    *   `jsonc.go` is special because it includes `removeComments` logic to strip JavaScript-style comments from a JSONC string before parsing it as regular JSON.
*   **`datetime.go`:** Provides `datetime.now()`, `datetime.format()`, `datetime.parse()`, `datetime.add()` (for adding durations), and `datetime.diff()`. These leverage Go's `time` package for robust date/time handling.
*   **`dotenv.go`:** `dotenv.load(path)` reads a `.env` file line by line, splits `KEY=VALUE` pairs, and uses `os.Setenv()` to make them available as environment variables. `dotenv.get(key, default)` retrieves them using `os.Getenv()`.
*   **`ft.go` (File Transfer):** `ft.download(url, path)` uses `http.Get` to fetch a file and `io.Copy` to save it. `ft.copy` and `ft.move` use `os.Open`, `os.Create`, `io.Copy`, and `os.Rename`.
*   **`random.go`:** `random.number()`, `random.int(min, max)`, and `random.string(length)` use Go's `crypto/rand` for cryptographically secure random data generation, which is generally preferred over `math/rand` for many use cases.
*   **`tar.go`:** `tar.create(archivePath, sourcePath, [compress])` uses `archive/tar` and optionally `compress/gzip` to create TAR archives. It walks the `sourcePath` (`filepath.Walk`) to add files and directories. `tar.extract` reads a TAR archive (handling GZip decompression if needed) and recreates the file structure. `tar.list` iterates through archive entries to list contents.
*   **`template.go`:** `template.parse(string)`, `template.parse_file(path)`, `template.parse_files(paths...)`, and `template.parse_glob(pattern)` use Go's `html/template` (or `text/template`) package to parse template definitions. They return a Lua function. When this Lua function is called with a data table, the Go template is executed with that data (after converting the Lua table to a Go map), and the rendered string is returned.
*   **`text.go`:** Offers a variety of string utilities like `text.trim()`, `text.lower()`, `text.split()`, `text.join()`, `text.replace()`, `text.contains()`, `text.pad_left()`, etc., mostly by calling equivalent functions from Go's `strings` package.
*   **`types.go`:** Provides functions like `types.is_callable()`, `types.is_integer()`, `types.is_number()`, etc. These offer more specific type checks than Lua's built-in `type()` function, for instance, distinguishing between an integer and a floating-point number if the Lua number representation allows.
*   **`utils.go`:** Contains general helper functions. For example, `utils.getfenv()` and `utils.setfenv()` interact with Lua's function environments (though their implementation might be simplified or rely on `debug.getinfo` for `getfenv`). `utils.unpack` mimics Lua's `table.unpack`. `utils.split` and `utils.join` are string utilities. `utils.escape` and `utils.unescape` might provide basic string escaping for specific contexts.
*   **`uuid.go`:** `uuid.v4()` and `uuid.v4_without_hyphens()` use a Go UUID library (like `github.com/google/uuid`) to generate Version 4 UUIDs. `uuid.is_valid()` parses a string to check if it's a valid UUID.
*   **`tablex.go`:** This is a more substantial utility module focused on advanced table operations.
    *   It includes functions for deep copying (`deepcopy`), comparing tables structurally (`compare`), transforming tables (`mapTable`, `filter`, `reduce`), restructuring (`flatten`, `slice`, `partition`, `rotate`, `shuffle`), and extracting data (`keys`, `values`).
    *   `pretty` provides a way to pretty-print Lua tables for debugging.
    *   `load` and `loadfile` can execute Lua code that returns a table.
    *   It introduces concepts like `map_new`, `set_new`, `ordered_map_new` which might return tables with specific metatables to simulate these data structures (though full implementation would require more extensive metatable programming).
    *   `array2d_*` functions provide utilities for working with 2D arrays (tables of tables), including creation, get/set, map, filter, and transpose. The `cols` count is often stored in the metatable of the 2D array for consistency.
    *   `permute` and `combinations` generate permutations and combinations of elements in a table.

This detailed breakdown illustrates that SolVM is far more than a simple Lua interpreter. It's a comprehensive runtime environment where Go's strengths in systems programming, networking, concurrency, and its rich standard library are systematically and thoughtfully exposed to Lua scripts. This fusion aims to make Lua a more powerful, productive, and enjoyable language for a broad spectrum of modern development tasks, without sacrificing the core simplicity and hackability that Lua is known for.