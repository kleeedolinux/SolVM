# SolVM

SolVM is a high-performance Lua runtime built in Go, based on gopher-lua. It provides concurrent execution, fault tolerance, and additional functionality for Lua scripts.

## Features

- Concurrent Lua script execution
- Fault tolerance with timeout mechanisms
- Custom functions for JSON encoding/decoding
- Thread-safe operations
- Built-in sleep functionality
- Module system for code organization
- Go-style concurrency primitives
- Memory leak detection and monitoring
- Enhanced error handling
- HTTP client functionality
- HTTP/HTTPS/WebSocket server support

## Installation

```bash
go get github.com/solvm
```

## Usage

Run Lua files using the command-line interface:

```bash
solvm [options] <lua-file>
```

Options:
- `-timeout`: Execution timeout (default: 5s)

Example:
```bash
solvm -timeout 10s examples/test.lua
```

## Server Features

SolVM provides HTTP, HTTPS, and WebSocket server functionality:

### Creating Servers
```lua
-- Create HTTP server
create_server("http_server", 8080, false)

-- Create HTTPS server
create_server("https_server", 8443, true, "cert.pem", "key.pem")
```

### HTTP Handlers
```lua
handle_http("http_server", "/", function(req)
    -- Access request details
    print("Method:", req.method)
    print("Path:", req.path)
    print("Query:", req.query)
    print("Headers:", req.headers)
    
    -- Return response
    return {
        status = 200,
        headers = {
            ["Content-Type"] = "application/json"
        },
        body = json_encode({
            message = "Hello from SolVM!"
        })
    }
end)
```

### WebSocket Handlers
```lua
handle_ws("http_server", "/ws", function(ws)
    -- Send message
    ws.send("Hello client!")
    
    -- Receive message
    local message = ws.receive()
    print("Received:", message)
end)
```

### Server Control
```lua
-- Start server
start_server("http_server")

-- Stop server
stop_server("http_server")
```

### Example Server
```lua
-- Create and configure server
create_server("my_server", 8080, false)

-- Handle HTTP requests
handle_http("my_server", "/", function(req)
    return {
        status = 200,
        body = "Hello World!"
    }
end)

-- Handle WebSocket connections
handle_ws("my_server", "/ws", function(ws)
    ws.send("Connected!")
    while true do
        local msg = ws.receive()
        if msg then
            ws.send("Echo: " .. msg)
        end
    end
end)

-- Start server
start_server("my_server")
```

## Monitoring and Error Handling

SolVM provides tools for monitoring memory usage and handling errors:

### Error Handling
```lua
on_error(function(err)
    print("Error occurred:", err)
end)
```

### Memory Monitoring
```lua
local mem_stats = check_memory()
print("Memory usage:")
print("  Allocation diff:", mem_stats.alloc_diff)
print("  Total allocation diff:", mem_stats.total_alloc_diff)
print("  System memory diff:", mem_stats.sys_diff)
print("  Number of goroutines:", mem_stats.goroutines)
```

### Goroutine Tracking
```lua
local goroutines = get_goroutines()
for id, name in pairs(goroutines) do
    print("Goroutine", id, ":", name)
end
```

## Concurrency Features

SolVM provides Go-style concurrency primitives:

### Goroutines
```lua
go(function()
    -- This runs in a separate goroutine
    print("Running in goroutine")
end)
```

### Channels
```lua
-- Create a channel
chan("my_channel", 5)  -- buffer size of 5

-- Send to channel
send("my_channel", "hello")

-- Receive from channel
local value = receive("my_channel")

-- Select from multiple channels
local value, channel = select("ch1", "ch2")
```

### Example
```lua
-- Create channels
chan("numbers", 5)
chan("results", 5)

-- Spawn a worker
go(function()
    while true do
        local num = receive("numbers")
        if num == nil then break end
        send("results", num * 2)
    end
end)

-- Send and receive
send("numbers", 42)
local result = receive("results")
```

## Module System

SolVM provides a module system for organizing and reusing code. Modules can be imported using the `import` function:

```lua
-- Import from local file
import("module_name")

-- Import from URL
import("https://raw.githubusercontent.com/user/repo/main/module.lua")
```

Modules are searched in the following locations:
1. Same directory as the current script
2. `modules` directory in the project root
3. Remote URLs (http:// or https://)

Example module (`modules/math_utils.lua`):
```lua
local math_utils = {}

function math_utils.add(a, b)
    return a + b
end

function math_utils.subtract(a, b)
    return a - b
end

return math_utils
```

Using the module:
```lua
-- Local import
import("math_utils")
local result = math_utils.add(5, 3)
print(result) -- Output: 8

-- Remote import
import("https://raw.githubusercontent.com/username/repo/main/utils.lua")
local remote_result = utils.some_function()
```

Note: Remote modules are cached in memory to prevent multiple downloads of the same module.

## Example Lua Script

```lua
local data = {
    name = "test",
    values = {1, 2, 3}
}

local json_str = json_encode(data)
print("Encoded:", json_str)

local decoded = json_decode(json_str)
print("Decoded name:", decoded.name)

print("Sleeping for 1 second...")
sleep(1)
print("Awake!")
```

## Custom Functions

SolVM provides several custom functions:

- `json_encode(value)`: Converts Lua values to JSON string
- `json_decode(string)`: Converts JSON string to Lua values
- `sleep(seconds)`: Pauses execution for specified seconds
- `import(module_name)`: Imports a Lua module
- `go(function)`: Runs a function in a new goroutine
- `chan(name, [buffer_size])`: Creates a new channel
- `send(channel, value)`: Sends a value to a channel
- `receive(channel)`: Receives a value from a channel
- `select(channel1, channel2, ...)`: Selects from multiple channels
- `on_error(function)`: Registers an error handler
- `check_memory()`: Returns memory usage statistics
- `get_goroutines()`: Returns list of active goroutines
- `wait()`: Waits for all goroutines to complete

### Network Functions

SolVM provides low-level network operations:

#### TCP Server
```lua
-- Create TCP server
local server = tcp_listen(8080)

-- Handle incoming connections
go(function()
    while true do
        local conn = receive(server)
        if conn then
            go(function()
                while true do
                    local data = conn.read()
                    if data then
                        print("Received:", data)
                        conn.write("Echo: " .. data)
                    else
                        break
                    end
                end
                conn.close()
            end)
        end
    end
end)
```

#### TCP Client
```lua
-- Connect to TCP server
local conn = tcp_connect("localhost", 8080)

-- Send data
conn.write("Hello, Server!")

-- Receive response
local response = conn.read()
print("Response:", response)

-- Close connection
conn.close()
```

#### UDP
```lua
-- Send UDP message
udp_sendto("localhost", 8080, "Hello, UDP!")

-- Receive UDP messages
local udp = udp_recvfrom(8080)
while true do
    local msg = udp.receive()
    if msg then
        print("Received from", msg.addr, ":", msg.data)
    end
end
```

#### DNS Resolution
```lua
-- Resolve hostname to IP
local ip = resolve_dns("example.com")
print("IP address:", ip)
```

Connection objects provide the following methods:
- `read()`: Read data from the connection
- `write(data)`: Write data to the connection
- `close()`: Close the connection

### Scheduler Functions

SolVM provides scheduling capabilities:

```lua
-- Set an interval (runs every 5 seconds)
local interval_id = set_interval(function()
    print("Interval triggered!")
end, 5)

-- Set a timeout (runs after 3 seconds)
local timeout_id = set_timeout(function()
    print("Timeout triggered!")
end, 3)

-- Set a cron job (runs every minute)
local cron_id = cron("0 * * * * *", function()
    print("Cron job triggered!")
end)
```

Cron Schedule Format:
```
┌─────────────── second (0 - 59)
│ ┌───────────── minute (0 - 59)
│ │ ┌─────────── hour (0 - 23)
│ │ │ ┌───────── day of month (1 - 31)
│ │ │ │ ┌─────── month (1 - 12)
│ │ │ │ │ ┌───── day of week (0 - 6)
│ │ │ │ │ │
* * * * * *
```

Examples:
- `"0 * * * * *"` - Every minute
- `"0 0 * * * *"` - Every hour
- `"0 0 0 * * *"` - Every day at midnight
- `"0 0 12 * * *"` - Every day at noon
- `"0 0 0 1 * *"` - First day of every month
- `"0 0 0 * * 0"` - Every Sunday

### Filesystem Functions

SolVM provides filesystem operations:

```lua
-- Read file content
local content = read_file("path/to/file.txt")
print(content)

-- Write content to file
write_file("path/to/file.txt", "Hello, World!")

-- List directory contents
local files = list_dir("path/to/directory")
for _, file in ipairs(files) do
    print("Name:", file.name)
    print("Is directory:", file.is_dir)
    print("Size:", file.size)
    print("Mode:", file.mode)
    print("Modified:", file.mod_time)
end
```

The `list_dir` function returns a table of file information with the following fields:
- `name`: File or directory name
- `is_dir`: Boolean indicating if it's a directory
- `size`: File size in bytes
- `mode`: File permissions and mode
- `mod_time`: Last modification time

## HTTP Client

SolVM provides HTTP client functionality for making HTTP requests:

### Basic HTTP Methods
```lua
-- GET request
local response = http_get("https://api.example.com/data")
print("Status:", response.status)
print("Body:", response.body)

-- POST request with JSON
local data = {
    name = "test",
    value = 42
}
local json_data = json_encode(data)
local response = http_post("https://api.example.com/data", json_data)

-- PUT request
local response = http_put("https://api.example.com/data", json_data)

-- DELETE request
local response = http_delete("https://api.example.com/data")
```

### Custom HTTP Requests
```lua
-- Custom request with headers
local headers = {
    ["Authorization"] = "Bearer token",
    ["Content-Type"] = "application/json"
}
local response = http_request("GET", "https://api.example.com/data", nil, headers)
```

### Response Format
All HTTP functions return a table with the following fields:
- `status`: HTTP status code
- `body`: Response body as string
- `headers`: Table containing response headers

### Error Handling
```lua
on_error(function(err)
    print("HTTP Error:", err)
end)
```

## Debug Functions

SolVM provides debugging capabilities:

```lua
-- Watch a file for changes and reload when modified
watch_file("script.lua", function()
    print("Script modified, reloading...")
    reload_script()
end)

-- Print current stack trace
trace()
```

The debug functions provide:
- `watch_file(path, callback)`: Watch a file for changes and call the callback when modified
- `reload_script()`: Reload the current script
- `trace()`: Print the current stack trace

Example usage:
```lua
-- Watch and auto-reload
watch_file("config.lua", function()
    print("Config changed, reloading...")
    reload_script()
end)

-- Debug stack trace
function some_function()
    trace()  -- Will print the current call stack
end
```

## License

MIT 