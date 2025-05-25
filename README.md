So, SolVM. Its a Lua runtime, right? But built in Go, kinda from gopher-lua. The cool part? It runs your Lua scripts all speedy and at the same time, and it tries not to break. Got some extra tricks up its sleeve too.

## Whatcha get with it?

*   Run Lua scripts together, at the same time
*   Keeps scripts from going haywire (with timeouts!)
*   Easy JSON stuff (encode/decode)
*   Safe for all that multi-taskin'
*   Sleepy time for your scripts (`sleep()`)
*   Modules! To keep your code tidy
*   Concurrency like in Go (goroutines, channels!)
*   Spots memory leaks (sometimes)
*   Better error messages (we hope!)
*   An HTTP client built right in
*   Make your own HTTP, HTTPS, or WebSocket servers! How cool is that?

## Installin is easy:

```bash
go get solvm
```
(You need Go for this, obvs)

## How to Use It

To run a Lua file, you just do:

```bash
solvm [options] <your_lua_file.lua>
```

Options ya got:
*   `-timeout`: How long before it gives up (default is 5s)

Like, for instance:
```bash
solvm -timeout 10s examples/test.lua
```

## Server Magic

SolVM can be a server! HTTP, HTTPS, even WebSockets. Pretty neat, huh?

### Makin' Servers
```lua
-- Make an HTTP server
create_server("http_server", 8080, false)

-- Or an HTTPS one (needs certs)
create_server("https_server", 8443, true, "cert.pem", "key.pem")
```

### Handle HTTP stuff
```lua
handle_http("http_server", "/", function(req)
    -- See what came in
    print("Method:", req.method)
    print("Path:", req.path)
    print("Query:", req.query) -- stuff after the ?
    print("Headers:", req.headers)

    -- Send somethin' back
    return {
        status = 200,
        headers = {
            ["Content-Type"] = "application/json"
        },
        body = json_encode({
            message = "Hello from SolVM! How u doin?"
        })
    }
end)
```

### WebSocket Handlers (for real-time chatty apps!)
```lua
handle_ws("http_server", "/ws", function(ws)
    -- Say hi!
    ws.send("Hello client!")

    -- Listen for stuff
    local message = ws.receive()
    print("Got this:", message)
end)
```

### Control Your Server
```lua
-- Start it up!
start_server("http_server")

-- Or shut it down
stop_server("http_server")
```

### Full Server Example
```lua
-- Make and set up a server
create_server("my_server", 8080, false)

-- What to do for normal web pages
handle_http("my_server", "/", function(req)
    return {
        status = 200,
        body = "Hello World from my SolVM server!"
    }
end)

-- What to do for WebSockets
handle_ws("my_server", "/ws", function(ws)
    ws.send("Youre connected!")
    while true do
        local msg = ws.receive()
        if msg then
            ws.send("You said: " .. msg)
        else
            break -- important!
        end
    end
end)

-- Kick it off
start_server("my_server")
```

## Keeping an Eye on Things & Handlin' Errors

SolVM gives ya some tools for checkin memory and dealin with problems:

### Error Stuff
```lua
on_error(function(err)
    print("Oh no, an error:", err)
end)
```

### Memory Checkin
```lua
local mem_stats = check_memory()
print("Hows the memory doin:")
print("  Allocation diff:", mem_stats.alloc_diff)
print("  Total allocation diff:", mem_stats.total_alloc_diff)
print("  System memory diff:", mem_stats.sys_diff)
print("  Number of goroutines:", mem_stats.goroutines)
```

### See Your Goroutines
```lua
local goroutines = get_goroutines()
for id, name in pairs(goroutines) do
    print("Goroutine", id, "is", name)
end
```

## Doin' Things at the Same Time (Concurrency)

SolVM has some Go-like ways to do many things at once:

### Goroutines (like little workers)
```lua
go(function()
    -- This code runs on its own, kinda
    print("Runnin in a goroutine!")
end)
```

### Channels (for talkin' between goroutines)
```lua
-- Make a channel
chan("my_channel", 5)  -- can hold 5 things

-- Put somethin in
send("my_channel", "hello")

-- Get somethin out
local value = receive("my_channel")

-- Pick from a few channels
local value, channel_name_it_came_from = select("ch1", "ch2")
```

### Example of Concurrency
```lua
-- Make some channels
chan("numbers", 5)
chan("results", 5)

-- Make a worker
go(function()
    while true do
        local num = receive("numbers")
        if num == nil then break end -- stop if channel closed or empty and done
        send("results", num * 2)
    end
end)

-- Send a number and get a result
send("numbers", 42)
local result = receive("results")
print("Result is:", result) -- should be 84
```

## Module System (for not makin a huge mess)

SolVM has a module system so you can split up your code. Use `import` to grab 'em:

```lua
-- Get a module from a local file
import("module_name")

-- Or from the internet!
import("https://raw.githubusercontent.com/user/repo/main/module.lua")
```

It looks for modules here:
1.  Same folder as your script
2.  A `modules` folder where your project is
3.  Web URLs (http:// or https://)

Example module (say, `modules/math_stuff.lua`):
```lua
local math_stuff = {}

function math_stuff.add(a, b)
    return a + b
end

function math_stuff.subtract(a, b)
    return a - b
end

return math_stuff
```

And usin' it:
```lua
-- Local one
import("math_stuff")
local sum = math_stuff.add(5, 3)
print(sum) -- prints 8

-- One from the web
import("https://raw.githubusercontent.com/someuser/somerepo/main/cool_utils.lua")
local remote_thing = cool_utils.do_something_cool()
```
FYI: Web modules get saved in memory so it dont download 'em over and over.

## Example Lua Script (a simple one)

```lua
local my_data = {
    name = "Tester",
    values = {10, 20, 30}
}

local json_version = json_encode(my_data)
print("As JSON:", json_version)

local back_to_lua = json_decode(json_version)
print("Decoded name:", back_to_lua.name)

print("Gonna sleep for 1 sec...")
sleep(1)
print("I'm awake!")
```

## Special Functions SolVM Gives Ya

SolVM adds some handy functions:

*   `json_encode(value)`: Lua to JSON string
*   `json_decode(string)`: JSON string to Lua
*   `sleep(seconds)`: Pauses for a bit
*   `import(module_name)`: Loads a Lua module
*   `go(function)`: Runs function in a new goroutine
*   `chan(name, [buffer_size])`: Makes a new channel
*   `send(channel, value)`: Sends to a channel
*   `receive(channel)`: Gets from a channel
*   `select(channel1, channel2, ...)`: Waits for any of several channels
*   `on_error(function)`: Sets what to do when an error happens
*   `check_memory()`: Shows memory stats
*   `get_goroutines()`: Lists active goroutines
*   `wait()`: Waits for all your goroutines to finish

### Networky Functions

SolVM can do some low-level network stuff too:

#### TCP Server
```lua
-- Make a TCP server listen on port 8080
local server_channel = tcp_listen(8080)

-- Handle new connections
go(function()
    while true do
        local conn = receive(server_channel) -- wait for a connection
        if conn then
            go(function() -- handle this connection in its own goroutine
                while true do
                    local data = conn.read()
                    if data then
                        print("TCP Server Got:", data)
                        conn.write("Server says echo: " .. data)
                    else
                        break -- connection closed or error
                    end
                end
                conn.close()
            end)
        else
            break -- server channel closed
        end
    end
end)
```

#### TCP Client
```lua
-- Connect to a TCP server
local conn = tcp_connect("localhost", 8080)

if conn then
    -- Send some data
    conn.write("Hello, TCP Server from SolVM!")

    -- Get a response
    local response = conn.read()
    print("TCP Client Got Response:", response)

    -- Close up
    conn.close()
else
    print("Couldnt connect to TCP server")
end
```

#### UDP
```lua
-- Send a UDP message (fire and forget!)
udp_sendto("localhost", 8080, "Hello, UDP world!")

-- Listen for UDP messages
local udp_listener_channel = udp_recvfrom(8080)
go(function()
    while true do
        local msg_table = receive(udp_listener_channel) -- {data="...", addr="..."}
        if msg_table then
            print("UDP Got from", msg_table.addr, ":", msg_table.data)
        else
            break -- channel closed
        end
    end
end)
```

#### DNS Lookup
```lua
-- Find IP for a hostname
local ip_addr = resolve_dns("example.com")
print("example.com IP is:", ip_addr)
```

Connection objects (from `tcp_listen` or `tcp_connect`) have:
*   `read()`: Reads data
*   `write(data)`: Writes data
*   `close()`: Closes it

### Scheduler Functions (do stuff later)

SolVM can schedule things to run:

```lua
-- Do this every 5 seconds
local interval_id = set_interval(function()
    print("Interval just happened!")
end, 5)

-- Do this once, after 3 seconds
local timeout_id = set_timeout(function()
    print("Timeout just happened!")
end, 3)

-- Run this based on a cron schedule (e.g., every minute)
local cron_id = cron("0 * * * * *", function()
    print("Cron job just ran!")
end)

-- To stop them:
-- clear_interval(interval_id)
-- clear_timeout(timeout_id)
-- clear_cron(cron_id)
```

Cron Format (kinda tricky):
```
.---------------- second (0 - 59)
|  .------------- minute (0 - 59)
|  |  .---------- hour (0 - 23)
|  |  |  .------- day of month (1 - 31)
|  |  |  |  .---- month (1 - 12)
|  |  |  |  |  .- day of week (0 - 6) (Sunday=0)
|  |  |  |  |  |
*  *  *  *  *  *
```
Examples:
*   `"0 * * * * *"` - Every minute, on the 0 second
*   `"0 0 * * * *"` - Every hour, at minute 0, second 0
*   `"0 0 0 * * *"` - Every day at midnight
*   `"0 0 12 * * *"`- Every day at noon
*   `"0 0 0 1 * *"` - First day of every month, at midnight
*   `"0 0 0 * * 0"` - Every Sunday at midnight

### Filesystem Stuff

SolVM can mess with files and folders:

```lua
-- Read a file's content
local stuff = read_file("my_stuff/notes.txt")
print(stuff)

-- Write to a file
write_file("my_stuff/output.txt", "This is new content!")

-- See what's in a directory
local dir_contents = list_dir("my_stuff/")
if dir_contents then
    for _, item in ipairs(dir_contents) do
        print("Name:", item.name)
        print("  Is it a folder?:", item.is_dir)
        print("  Size:", item.size)
        print("  Mode (permissions):", item.mode)
        print("  Last changed:", item.mod_time)
    end
end
```
`list_dir` gives you a table of files/folders with:
*   `name`: Name of the thing
*   `is_dir`: True if it's a directory, false if not
*   `size`: File size (bytes)
*   `mode`: File mode (permissions n stuff)
*   `mod_time`: When it was last changed

## HTTP Client (for fetching web pages)

SolVM has its own HTTP client to grab stuff from the web:

### Basic HTTP
```lua
-- GET request
local resp = http_get("https://api.someplace.com/info")
print("Status Code:", resp.status)
print("Body:", resp.body)

-- POST request (with JSON)
local my_data = { item = "thing", quantity = 2 }
local json_payload = json_encode(my_data)
local resp_post = http_post("https://api.someplace.com/submit", json_payload)

-- PUT request
local resp_put = http_put("https://api.someplace.com/update/123", json_payload)

-- DELETE request
local resp_delete = http_delete("https://api.someplace.com/remove/123")
```

### Custom HTTP Requests (more control)
```lua
-- Custom request with your own headers
local my_headers = {
    ["Authorization"] = "Bearer mysecrettoken",
    ["Content-Type"] = "application/json"
}
local resp_custom = http_request("GET", "https://api.someplace.com/secretdata", nil, my_headers)
```

### What You Get Back (Response)
All HTTP functions give you a table with:
*   `status`: The HTTP status code (like 200 for OK, 404 for not found)
*   `body`: The page content (as a string)
*   `headers`: A table of the response headers

### Error Handling for HTTP
```lua
on_error(function(err_msg)
    print("HTTP Request Broke:", err_msg)
end)
-- Or you can check resp.status for non-2xx codes
```

## Debug Functions (for when things go wrong)

SolVM got a few tools to help ya debug:

```lua
-- Watch a file. If it changes, run the function.
watch_file("my_script.lua", function()
    print("my_script.lua changed! Reloading...")
    reload_script() -- Tries to reload the main script
end)

-- Print where you are in the code right now
trace()
```
The debug functions are:
*   `watch_file(path_to_file, callback_function)`: Watches a file, calls callback on change.
*   `reload_script()`: Tries to reload the current main script.
*   `trace()`: Prints the current call stack (what functions called what).

Example use:
```lua
-- Watch your config and reload if it changes
watch_file("config.lua", function()
    print("Config file changed, reloading script...")
    reload_script()
end)

function a_deep_function()
    trace() -- This will show you how you got here
end

function another_function()
    a_deep_function()
end

another_function()
```

## License

MIT (So, do what ya want with it, pretty much)