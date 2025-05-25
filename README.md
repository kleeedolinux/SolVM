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
- UUID generation
- Random number generation
- TOML/YAML/JSONC support
- Text manipulation utilities
- Cryptographic functions (AES, DES, RC4, RSA, hashing)

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

## Additional Modules

### UUID Module
```lua
-- Generate UUID v4
local uuid1 = uuid.v4()
local uuid2 = uuid.v4_without_hyphens()

-- Validate UUID
local is_valid = uuid.is_valid(uuid1)
```

### Random Module
```lua
-- Generate random number between 0 and 1
local random_num = random.number()

-- Generate random integer between min and max
local random_int = random.int(1, 100)

-- Generate random string of specified length
local random_str = random.string(10)
```

### TOML Module
```lua
-- Encode table to TOML
local toml_str = toml.encode({
    title = "Example",
    owner = {
        name = "John Doe",
        age = 30
    }
})

-- Decode TOML to table
local decoded = toml.decode(toml_str)
```

### YAML Module
```lua
-- Encode table to YAML
local yaml_str = yaml.encode({
    name = "Example",
    version = "1.0",
    config = {
        debug = true
    }
})

-- Decode YAML to table
local decoded = yaml.decode(yaml_str)
```

### JSONC Module
```lua
-- Decode JSONC (JSON with comments)
local jsonc_str = [[
{
    // This is a comment
    "name": "Example",
    /* This is a block comment */
    "config": {
        "debug": true
    }
}
]]
local decoded = jsonc.decode(jsonc_str)

-- Encode table to JSONC
local encoded = jsonc.encode(decoded)
```

### Text Module
```lua
-- String manipulation
local text = "  Hello, World!  "
local trimmed = text.trim(text)
local lower = text.lower(text)
local upper = text.upper(text)
local title = text.title(text)

-- Split and join
local words = text.split(text.trim(text), " ")
local joined = text.join(words, ", ")

-- Replace
local replaced = text.replace(text, "World", "Lua")

-- String checks
local contains = text.contains(text, "Hello")
local starts = text.starts_with(text, "  H")
local ends = text.ends_with(text, "!  ")

-- Padding
local padded = text.pad_left("123", 5, "0")

-- Repeat
local repeated = text.repeat("abc", 3)
```

### Crypto Module
```lua
-- Hash functions
local data = "Hello, World!"
local md5_hash = crypto.md5(data)
local sha1_hash = crypto.sha1(data)
local sha256_hash = crypto.sha256(data)
local sha512_hash = crypto.sha512(data)

-- Base64 encoding/decoding
local encoded = crypto.base64_encode(data)
local decoded = crypto.base64_decode(encoded)

-- AES encryption/decryption
local key = "1234567890123456"  -- 16 bytes for AES-128
local iv = "1234567890123456"   -- 16 bytes for AES-128
local encrypted = crypto.aes_encrypt(data, key, iv)
local decrypted = crypto.aes_decrypt(encrypted, key, iv)

-- DES encryption/decryption
local des_key = "12345678"  -- 8 bytes for DES
local des_iv = "12345678"   -- 8 bytes for DES
local des_encrypted = crypto.des_encrypt(data, des_key, des_iv)
local des_decrypted = crypto.des_decrypt(des_encrypted, des_key, des_iv)

-- RC4 encryption/decryption
local rc4_key = "mysecretkey"
local rc4_encrypted = crypto.rc4_encrypt(data, rc4_key)
local rc4_decrypted = crypto.rc4_decrypt(rc4_encrypted, rc4_key)

-- RSA key generation
local rsa_keys = crypto.rsa_generate(2048)
print("Private Key:", rsa_keys.private)
print("Public Key:", rsa_keys.public)

-- Random bytes generation
local random = crypto.random_bytes(32)
```

The crypto module provides the following functions:

#### Hash Functions
- `crypto.md5(data)`: Calculate MD5 hash
- `crypto.sha1(data)`: Calculate SHA1 hash
- `crypto.sha256(data)`: Calculate SHA256 hash
- `crypto.sha512(data)`: Calculate SHA512 hash

#### Base64
- `crypto.base64_encode(data)`: Encode data to Base64
- `crypto.base64_decode(encoded)`: Decode Base64 data

#### AES Encryption
- `crypto.aes_encrypt(data, key, iv)`: Encrypt data using AES-CBC
- `crypto.aes_decrypt(encrypted, key, iv)`: Decrypt AES-CBC data

#### DES Encryption
- `crypto.des_encrypt(data, key, iv)`: Encrypt data using DES-CBC
- `crypto.des_decrypt(encrypted, key, iv)`: Decrypt DES-CBC data

#### RC4 Encryption
- `crypto.rc4_encrypt(data, key)`: Encrypt data using RC4
- `crypto.rc4_decrypt(encrypted, key)`: Decrypt RC4 data

#### RSA
- `crypto.rsa_generate([bits])`: Generate RSA key pair (default 2048 bits)

#### Random
- `crypto.random_bytes(length)`: Generate cryptographically secure random bytes

Note: All encrypted data is returned as Base64-encoded strings for safe transmission and storage. 