
# SolVM: A Comprehensive Guide

SolVM is a versatile scripting environment powered by Lua, designed to offer an extensive suite of built-in functionalities. It empowers developers with tools for concurrent programming, robust networking capabilities, cryptographic operations, efficient file system management, seamless data format handling, and much more. This guide aims to walk you through the core features and modules available in SolVM, helping you leverage its full potential.

## Table of Contents

1.  [Core SolVM Concepts](#core-solvm-concepts)
    *   [Concurrent Programming: Goroutines and Channels](#concurrent-programming-goroutines-and-channels)
    *   [Graceful Error Management](#graceful-error-management)
    *   [Fundamental Utilities](#fundamental-utilities)
2.  [Exploring SolVM's Modules and Built-in Features](#exploring-solvms-modules-and-built-in-features)
    *   [Cryptography Services (`crypto`)](#cryptography-services-crypto)
    *   [Debugging and Application Insights](#debugging-and-application-insights)
    *   [Network Operations](#network-operations)
        *   [Domain Name System (DNS) Resolution](#domain-name-system-dns-resolution)
        *   [Interacting with Web Services: HTTP Client](#interacting-with-web-services-http-client)
        *   [Building Web Applications: HTTP Server](#building-web-applications-http-server)
        *   [Reliable Network Streams: TCP Communication](#reliable-network-streams-tcp-communication)
        *   [Fast Datagrams: UDP Communication](#fast-datagrams-udp-communication)
    *   [Managing Files and Directories](#managing-files-and-directories)
    *   [Dynamic Development: Hot Reloading](#dynamic-development-hot-reloading)
    *   [Structuring Your Code: Importing Modules](#structuring-your-code-importing-modules)
        *   [Defining Module Metadata](#defining-module-metadata)
    *   [Working with Various Data Formats](#working-with-various-data-formats)
        *   [JSON Handling (`json`)](#json-handling-json)
        *   [Universally Unique Identifiers (`uuid`)](#universally-unique-identifiers-uuid)
        *   [Generating Random Data (`random`)](#generating-random-data-random)
        *   [TOML Configuration Files (`toml`)](#toml-configuration-files-toml)
        *   [YAML Data Serialization (`yaml`)](#yaml-data-serialization-yaml)
        *   [JSON with Comments (`jsonc`)](#json-with-comments-jsonc)
        *   [Comma-Separated Values (`csv`)](#comma-separated-values-csv)
        *   [INI Configuration Files (`ini`)](#ini-configuration-files-ini)
    *   [Advanced Text Manipulation (`text`)](#advanced-text-manipulation-text)
    *   [Handling Environment Variables (`dotenv`)](#handling-environment-variables-dotenv)
    *   [Precise Date and Time Operations (`datetime`)](#precise-date-and-time-operations-datetime)
    *   [Streamlining File Transfers (`ft`)](#streamlining-file-transfers-ft)
    *   [Creating and Managing Archives (`tar`)](#creating-and-managing-archives-tar)
    *   [Interacting with the Operating System (`os`)](#interacting-with-the-operating-system-os)
    *   [Template Engine (`template`)](#template-engine-template)
    *   [General Purpose Utilities (`utils`)](#general-purpose-utilities-utils)
    *   [Advanced Type System (`types`)](#advanced-type-system-types)
3.  [Creating Your Own Modules: An Example (`math_utils.lua`)](#creating-your-own-modules-an-example-math_utilslua)

---

## 1. Core SolVM Concepts

At the heart of SolVM are powerful primitives that form the foundation for building complex applications.

### Concurrent Programming: Goroutines and Channels

SolVM embraces concurrent programming through a system inspired by Go's goroutines and channels, enabling you to perform multiple tasks seemingly simultaneously and communicate between them efficiently.

To manage communication between these concurrent tasks, you first create a **channel**. This is done using the `chan(name, [buffer_size])` function. For instance, `chan("data_stream", 5)` would establish a channel named "data_stream" capable of holding up to 5 items before a send operation blocks. If you omit the `buffer_size` or set it to zero, the channel becomes unbuffered, meaning a `send` operation will wait until a `receive` operation is ready for that specific item, ensuring synchronized data exchange.

To execute a piece of code concurrently, you encapsulate it within a function and then pass this function to the `go(function)` command. For example, `go(process_data_task)` would launch the `process_data_task` function in a new, lightweight execution thread, often referred to as a goroutine. This allows the main script to continue its execution without waiting for `process_data_task` to complete.

Once goroutines and channels are in place, data is exchanged using `send(channel_name, value)` to transmit a `value` to the channel identified by `channel_name`, and `receive(channel_name)` to retrieve a value from it. The `receive` function is blocking; it will pause the goroutine's execution until a value is available on the channel. A common pattern is for `receive` to return `nil` when a channel has been closed by the sender and all buffered items have been consumed, signaling to the receiver that no more data will arrive.

For scenarios where a goroutine needs to monitor multiple channels and react to the first one that becomes ready, SolVM provides the `select(channel_name1, channel_name2, ...)` function. This powerful construct pauses execution until an operation (typically a receive, but can also be a send if channels are used bidirectionally in more advanced patterns) can proceed on one of the listed channels. It then returns two values: the `value` itself and the `channel_name` from which it originated. If multiple channels are ready simultaneously, `select` makes a pseudo-random choice.

To ensure that your main script or a parent goroutine doesn't terminate prematurely before all its spawned concurrent tasks have finished their work, SolVM offers the `wait()` function. Calling `wait()` will block the current execution flow until all goroutines initiated with `go` have completed. This is crucial for orderly application shutdown and to prevent data loss or incomplete operations.

**Illustrative Example of Concurrency:**
```lua
-- Create a buffered channel for numbers and another for results
chan("numbers_pipeline", 5)
chan("results_pipeline", 5)

-- Spawn a worker goroutine to process numbers
go(function()
    print("Worker goroutine started.")
    while true do
        local num = receive("numbers_pipeline")
        if num == nil then
            print("Numbers channel closed, worker exiting.")
            break -- Exit loop if channel is closed and empty
        end
        print("Worker received:", num)
        sleep(0.5)  -- Simulate some processing time
        local result = num * num
        send("results_pipeline", result)
        print("Worker sent result:", result)
    end
end)

-- Main flow: Send numbers to the worker
print("Main: Sending numbers...")
for i = 1, 3 do
    print("Main: Sending:", i)
    send("numbers_pipeline", i)
end

-- Important: Close the numbers channel to signal the worker that no more data is coming
-- This allows the worker's receive("numbers_pipeline") to eventually return nil.
-- If not closed, the worker might block indefinitely on receive.
-- close("numbers_pipeline") -- Assuming 'close' function exists as per Go patterns.
                          -- If not, an alternative signaling mechanism (e.g., a special 'done' value)
                          -- or a fixed number of receives would be needed.
                          -- For this example's flow, let's assume the worker processes a known number of items.


-- Main flow: Receive and print results
print("Main: Receiving results...")
for i = 1, 3 do
    local result = receive("results_pipeline")
    print("Main: Received result:", result)
end

-- Using select for readiness on multiple channels
chan("event_ch_A")
chan("event_ch_B")

go(function()
    sleep(0.7)
    send("event_ch_A", "Event A occurred!")
end)

go(function()
    sleep(0.4)
    send("event_ch_B", "Event B fired!")
end)

print("Main: Waiting for an event using select...")
local event_data, source_channel = select("event_ch_A", "event_ch_B")
print("Main: Select received '"..tostring(event_data).."' from channel '"..tostring(source_channel).."'")

-- Wait for all spawned goroutines to complete their execution before the script ends
print("Main: Waiting for all goroutines to finish...")
wait()
print("Main: All operations complete. Exiting.")
```

### Graceful Error Management

Robust applications require a solid error handling strategy. SolVM provides a mechanism to define a global error handler using the `on_error(function(err))` function. When you register a callback function with `on_error`, this function will be invoked if an unhandled error or exception occurs anywhere in your script, including within any goroutines. The error object or message (`err`) is passed as an argument to your handler, allowing you to log the error, perform cleanup operations, or notify an administrator.

**Example of Global Error Handling:**
```lua
on_error(function(runtime_error)
    print("Critical Error Captured: " .. tostring(runtime_error))
    -- In a real application, you might log this to a file or an external monitoring service.
end)

go(function()
    print("Goroutine attempting a risky operation...")
    local a = 10
    local b = 0
    local c = a / b -- This will cause a division by zero error
    print("This line will not be reached in the goroutine.")
end)

print("Main script continues while goroutine might encounter an error.")
sleep(0.1) -- Give goroutine a chance to run and error out

-- If the error occurs in the main thread (not in a goroutine),
-- SolVM's top-level behavior might also be influenced by on_error,
-- or it might terminate the script directly after invoking the handler.
-- error("This is a deliberate error in the main thread.")

wait() -- Wait for goroutines to finish (or error out)
print("Script finished or error handled.")
```

### Fundamental Utilities

SolVM includes essential utility functions for everyday scripting tasks:

The `print(...)` function is your go-to for outputting information. It takes one or more arguments and displays their string representations to the standard output console, separated by tabs, followed by a newline.

For introducing delays or pacing operations, the `sleep(seconds)` function is available. It pauses the execution of the current goroutine (or the main script thread) for the specified duration in `seconds`. The duration can be a fractional number, allowing for sub-second precision, for example, `sleep(0.25)` would pause for 250 milliseconds.

---

## 2. Exploring SolVM's Modules and Built-in Features

SolVM is equipped with a rich set of modules that extend Lua's core capabilities, covering areas from cryptography to network communication and data parsing.

### Cryptography Services (`crypto`)

The `crypto` module is your toolkit for a wide range of cryptographic operations, essential for securing data and verifying integrity.

For generating **cryptographic hashes**, the module provides several industry-standard algorithms. You can compute an MD5 hash using `crypto.md5(data)`, a SHA1 hash with `crypto.sha1(data)`, a SHA256 hash via `crypto.sha256(data)`, and a SHA512 hash using `crypto.sha512(data)`. Each of these functions takes the input `data` (typically a string) and returns its corresponding hash value as a hexadecimal string.

When it comes to **encoding binary data into a text format**, Base64 is a common choice. The `crypto` module offers `crypto.base64_encode(data)` to convert your input `data` into a Base64 string, and `crypto.base64_decode(encoded_data)` to revert a Base64 string back to its original form.

For **symmetric encryption**, SolVM supports several algorithms.
**AES (Advanced Encryption Standard)** is available through `crypto.aes_encrypt(data, key, iv)` for encryption and `crypto.aes_decrypt(encrypted_data, key, iv)` for decryption. These functions require the `data` to be encrypted/decrypted, a secret `key` (e.g., 16 bytes for AES-128, 24 bytes for AES-192, or 32 bytes for AES-256), and an `iv` (Initialization Vector, typically 16 bytes for AES, crucial for security in modes like CBC).
Similarly, **DES (Data Encryption Standard)** is supported with `crypto.des_encrypt(data, key, iv)` and `crypto.des_decrypt(encrypted_data, key, iv)`. For DES, the `key` should be 8 bytes, and the `iv` is also 8 bytes.
The **RC4 stream cipher** is provided via `crypto.rc4_encrypt(data, key)` for encryption and `crypto.rc4_decrypt(rc4_encrypted_data, key)` for decryption. RC4 only requires the `data` and a secret `key`.

For **asymmetric cryptography**, the module includes functionality for **RSA key pair generation**. Calling `crypto.rsa_generate(bits)` (e.g., `crypto.rsa_generate(2048)` for a 2048-bit key) will return a table containing the `private` and `public` keys in PEM format.

Lastly, if you need a source of **cryptographically secure random bytes**, you can use `crypto.random_bytes(count)` to generate a string containing `count` random bytes.

**Cryptography Showcase:**
```lua
local data_to_secure = "This is highly confidential information!"
print("Original Data:", data_to_secure)

-- Hash functions
print("MD5 Hash:", crypto.md5(data_to_secure))
print("SHA256 Hash:", crypto.sha256(data_to_secure))

-- Base64 encoding/decoding
local base64_encoded = crypto.base64_encode(data_to_secure)
print("Base64 Encoded:", base64_encoded)
local base64_decoded = crypto.base64_decode(base64_encoded)
print("Base64 Decoded:", base64_decoded)

-- AES encryption/decryption
local aes_key = "0123456789abcdef"  -- 16 bytes for AES-128
local aes_iv = "fedcba9876543210"   -- 16 bytes for IV
local aes_encrypted = crypto.aes_encrypt(data_to_secure, aes_key, aes_iv)
print("AES Encrypted (hex preview):", string.sub(aes_encrypted:gsub(".", function(c) return string.format("%02x", string.byte(c)) end), 1, 32) .. "...")
local aes_decrypted = crypto.aes_decrypt(aes_encrypted, aes_key, aes_iv)
print("AES Decrypted:", aes_decrypted)

-- RSA Key Generation
local rsa_keys = crypto.rsa_generate(2048)
print("RSA Private Key (first 60 chars):", string.sub(rsa_keys.private, 1, 60) .. "...")
print("RSA Public Key (first 60 chars):", string.sub(rsa_keys.public, 1, 60) .. "...")

-- Random Bytes
local secure_random_data = crypto.random_bytes(16)
print("Secure Random Bytes (16, hex):", secure_random_data:gsub(".", function(c) return string.format("%02x", string.byte(c)) end))
```

### Debugging and Application Insights

SolVM includes features to aid in debugging your scripts and understanding their runtime behavior.

To understand the call sequence leading to a particular point in your code, you can use the `trace()` function. When called, `trace()` prints the current goroutine's stack trace to the console, showing the chain of function calls.

For development workflows where you're actively modifying a script, the `watch_file(filepath, callback_function)` function is invaluable. It monitors the specified `filepath` for any changes. When a change is detected, the provided `callback_function` is executed. This is often used in conjunction with `reload_script()`, which, when called within the callback, instructs SolVM to reload and re-execute the current script, allowing for a live-reloading development experience.

To monitor the resource consumption of your SolVM application, especially concerning memory and concurrency, the `check_memory()` function provides vital statistics. It returns a table containing information such as `alloc_diff` (difference in allocated memory since the last check or start), `total_alloc_diff` (total memory allocated since start), `sys_diff` (system memory difference), and `goroutines` (the current number of active goroutines). This is particularly useful for identifying potential memory leaks or understanding the concurrency profile of your application over time.

If you need to inspect the currently running goroutines, the `get_goroutines()` function returns a table where keys are goroutine IDs and values are their names or entry points. This can help in understanding what concurrent tasks are active at any given moment.

**Debugging Utilities Example:**
```lua
print("Debug demonstration script started.")

-- Function to illustrate stack trace
function level_one()
    function level_two()
        function level_three()
            print("Current execution stack trace:")
            trace() -- Show the call stack
        end
        level_three()
    end
    level_two()
end

-- Call the function to see the trace
level_one()

-- Example of watching a file (assuming this script is named 'my_debug_script.lua')
-- In a real scenario, you'd watch the script file itself.
-- For this example, let's simulate it.
-- watch_file("my_debug_script.lua", function()
--     print("my_debug_script.lua has changed! Reloading script...")
--     reload_script()
-- end)

-- Memory and Goroutine check
local initial_mem_stats = check_memory()
print("Initial Memory - Goroutines:", initial_mem_stats.goroutines, "AllocDiff:", initial_mem_stats.alloc_diff)

go(function()
    print("Spawned a new goroutine for memory check.")
    local temp_table = {}
    for i = 1, 500 do temp_table[i] = "some data " .. i end
    sleep(0.2) -- Keep goroutine alive for a bit
    print("Goroutine finishing.")
end)

sleep(0.3) -- Allow time for goroutine to run and potentially affect memory

local mid_mem_stats = check_memory()
print("Mid Memory - Goroutines:", mid_mem_stats.goroutines, "AllocDiff:", mid_mem_stats.alloc_diff)

local active_goroutines = get_goroutines()
print("Currently active goroutines:")
for id, name in pairs(active_goroutines) do
    print("  ID:", id, "Name/Entry:", name)
end

print("Script will now loop. Modify 'my_debug_script.lua' to see watch_file/reload_script in action if enabled.")
-- To see watch_file, you'd need to uncomment it and save the file.
-- For now, we'll just wait for other goroutines.
wait()
print("Debug demonstration finished.")
```

### Network Operations

SolVM provides comprehensive support for various networking protocols.

#### Domain Name System (DNS) Resolution

To translate human-readable domain names into IP addresses, SolVM offers the `resolve_dns(domain_name)` function. You pass it a `domain_name` string (e.g., "google.com"), and it returns the corresponding IP address as a string. If the domain cannot be resolved, it may return `nil` or an error indicator.

**DNS Lookup Example:**
```lua
print("Performing DNS lookups...")
local domains_to_check = { "google.com", "github.com", "nonexistentdomain123abc.com" }

for _, domain in ipairs(domains_to_check) do
    local ip_address = resolve_dns(domain)
    if ip_address then
        print(string.format("Domain: %-30s IP Address: %s", domain, ip_address))
    else
        print(string.format("Domain: %-30s Could not resolve.", domain))
    end
end
-- The script would typically need a loop or wait() to keep running if other async operations were present.
```

#### Interacting with Web Services: HTTP Client

SolVM allows your scripts to act as HTTP clients, making requests to web servers and APIs.

For making a simple **GET request**, you can use `http_get(url)`. This function takes the target `url` as a string and returns a table representing the response. This response table typically includes fields like `status` (the HTTP status code, e.g., 200), `headers` (a table of response headers), and `body` (the content of the response).

To send data with a **POST request**, particularly JSON data, you would use `http_post(url, data_body, [content_type])`. You provide the `url`, the `data_body` (e.g., a JSON string), and optionally the `content_type` (which defaults to "application/json" if `data_body` looks like JSON, but can be specified, e.g., "application/x-www-form-urlencoded"). Like `http_get`, it returns a response table.

For more control over the HTTP request, including custom methods (PUT, DELETE, etc.) and headers, SolVM provides a generic `http_request(method, url, body, headers_table)` function. Here, `method` is the HTTP method string (e.g., "GET", "POST", "PUT"), `url` is the target, `body` is the request payload (can be an empty string for methods like GET), and `headers_table` is a Lua table where keys are header names and values are header values.

It's good practice to set up an error handler using `on_error` when dealing with network requests, as they can fail for various reasons (network issues, server errors, invalid URLs).

**HTTP Client Operations:**
```lua
-- Ensure an error handler is set for network operations
on_error(function(err)
    print("HTTP Client Error:", err)
end)

-- Test GET request
print("Performing GET request to httpbin.org/get...")
local get_response = http_get("https://httpbin.org/get")
if get_response then
    print("GET Response Status:", get_response.status)
    print("GET Response Body (first 100 chars):", string.sub(get_response.body, 1, 100))
else
    print("GET request failed.")
end

-- Test POST request with JSON data
print("\nPerforming POST request to httpbin.org/post...")
local post_payload = { message = "Hello from SolVM!", value = 123 }
local post_payload_json = json_encode(post_payload) -- Assumes json_encode is available

local post_response = http_post("https://httpbin.org/post", post_payload_json)
if post_response then
    print("POST Response Status:", post_response.status)
    -- httpbin.org/post echoes the data sent
    print("POST Response Body (first 100 chars):", string.sub(post_response.body, 1, 100))
else
    print("POST request failed.")
end

-- Test custom request with specific headers
print("\nPerforming Custom GET request with headers to httpbin.org/headers...")
local custom_headers = {
    ["X-SolVM-Version"] = "1.0",
    ["Authorization"] = "Bearer my-secret-token"
}
local custom_req_response = http_request("GET", "https://httpbin.org/headers", "", custom_headers)
if custom_req_response then
    print("Custom Request Response Status:", custom_req_response.status)
    print("Custom Request Response Body (first 150 chars):", string.sub(custom_req_response.body, 1, 150))
else
    print("Custom request failed.")
end

-- Test error handling with an invalid URL
print("\nTesting request to an invalid URL...")
local error_test_response = http_get("http://this-is-not-a-real-domain-at-all.xyz")
if error_test_response == nil then
    print("Correctly handled error for invalid URL (response was nil as expected).")
end
```

#### Building Web Applications: HTTP Server

SolVM can also host HTTP servers, allowing you to build web applications and APIs.

First, you create a server instance using `create_server(server_name, port, use_tls)`. You provide a unique `server_name` for reference, the `port` number it should listen on, and a boolean `use_tls` indicating whether it should use HTTPS (true) or HTTP (false). If `use_tls` is true, you'll typically need to configure certificate and key files.

Once the server is created, you define handlers for different HTTP paths using `handle_http(server_name, path_pattern, handler_function)`. The `server_name` refers to the server you created, `path_pattern` is the URL path (e.g., "/", "/api/users"), and `handler_function` is a Lua function that will be executed when a request matches the path. This handler function receives a `request` table as an argument, containing details like `request.method`, `request.path`, `request.query` (parsed query parameters), `request.headers`, and `request.body`. The handler function must return a response table, which should include `status` (HTTP status code), `headers` (a table of response headers), and `body` (the response content as a string).

For real-time, bidirectional communication, SolVM supports WebSockets. You can set up a WebSocket endpoint using `handle_ws(server_name, path_pattern, ws_handler_function)`. The `ws_handler_function` receives a `websocket_connection` object when a client connects. This object has methods like `websocket_connection.send(message)` to send data to the client and `websocket_connection.receive()` to read messages from the client (which blocks until a message arrives or returns `nil` if the connection closes).

After defining all your handlers, you start the server using `start_server(server_name)`. This will begin listening for incoming connections. Typically, your script will then enter an infinite loop (e.g., `while true do sleep(1) end`) to keep the server running.

**HTTP & WebSocket Server Example:**
```lua
-- Create an HTTP server instance named "my_app_server" on port 8080, not using TLS (HTTP)
create_server("my_app_server", 8080, false)
print("HTTP server 'my_app_server' created, configured for port 8080.")

-- Handle HTTP requests to the root path "/"
handle_http("my_app_server", "/", function(req_details)
    print("Received HTTP request on '/':")
    print("  Method:", req_details.method)
    print("  Path:", req_details.path)
    print("  Query Params:", json_encode(req_details.query)) -- Assuming json_encode for pretty print

    local response_body = {
        greeting = "Welcome to the SolVM HTTP Server!",
        your_request = {
            method = req_details.method,
            path = req_details.path,
            query = req_details.query
        }
    }
    
    -- Return the response
    return {
        status = 200, -- OK
        headers = {
            ["Content-Type"] = "application/json",
            ["X-Powered-By"] = "SolVM"
        },
        body = json_encode(response_body) -- Send JSON back
    }
end)

-- Handle WebSocket connections on the "/livefeed" path
handle_ws("my_app_server", "/livefeed", function(ws_conn)
    print("New WebSocket connection established on /livefeed.")
    
    -- Send an initial welcome message to the connected client
    ws_conn.send(json_encode({ type = "system_message", content = "Connection successful! Welcome to the live feed." }))
    
    -- Loop to receive and echo messages
    while true do
        local client_message_json = ws_conn.receive()
        if client_message_json == nil then
            -- nil indicates the client disconnected or an error occurred
            print("WebSocket client disconnected from /livefeed.")
            break
        end
        
        print("Received WebSocket message on /livefeed:", client_message_json)
        
        -- Echo the message back, perhaps wrapped or processed
        local client_message = json_decode(client_message_json) -- Assuming json_decode
        ws_conn.send(json_encode({ type = "echo_response", original_message = client_message }))
    end
    
    print("WebSocket handler for /livefeed finished for this connection.")
end)

-- Start the server to begin listening for connections
print("Starting 'my_app_server' on port 8080...")
start_server("my_app_server")

-- Keep the script running to serve requests
print("Server is running. Access at http://localhost:8080/ and ws://localhost:8080/livefeed")
while true do
    sleep(60) -- Sleep for a while, server runs in background threads
    print("Server heartbeat...")
end
```

#### Reliable Network Streams: TCP Communication

For establishing reliable, connection-oriented network communication, SolVM provides TCP (Transmission Control Protocol) functionalities.

To create a **TCP server** that listens for incoming connections, you use `tcp_listen(port)`. This function binds to the specified `port` on all available network interfaces and returns a special channel-like object (let's call it `server_listener_channel`). You can then `receive` from this `server_listener_channel`; each successful receive operation yields a `connection` object representing a new client connection.

Each `connection` object obtained from the server (or created by a client) has methods for communication:
*   `connection.read()`: Reads data sent by the remote peer. It blocks until data is available or the connection is closed (returning `nil`).
*   `connection.write(data)`: Sends `data` (a string) to the remote peer.
*   `connection.close()`: Closes the TCP connection.
*   `connection.remote_addr`: A property providing the remote address of the connected peer.

Typically, for each new connection accepted by the server, you'd spawn a new goroutine to handle communication with that client independently.

To act as a **TCP client** and connect to a TCP server, you use `tcp_connect(host, port)`. This function attempts to establish a connection to the specified `host` (IP address or hostname) and `port`. If successful, it returns a `connection` object, similar to the one received by the server, which you can then use to `read` and `write` data.

**TCP Server and Client Demonstration:**
```lua
local tcp_port = 8088

-- Create a TCP server
print("Attempting to start TCP server on port " .. tcp_port)
local server_listener_channel = tcp_listen(tcp_port)
if not server_listener_channel then
    print("Failed to start TCP server. Port might be in use.")
    return -- Exit if server cannot start
end
print("TCP server is now listening on port " .. tcp_port)

-- Goroutine to handle incoming client connections
go(function()
    print("Server: Connection acceptor goroutine started.")
    while true do
        local client_conn = receive(server_listener_channel) -- Wait for a new connection
        if client_conn then
            print("Server: Accepted new connection from: " .. client_conn.remote_addr)
            
            -- Spawn another goroutine to handle this specific client
            go(function(conn)
                print("Server: Handler goroutine started for client " .. conn.remote_addr)
                while true do
                    local received_data = conn.read() -- Read data from the client
                    if received_data then
                        print("Server: Received from " .. conn.remote_addr .. ": '" .. received_data .. "'")
                        local response_data = "Server echoes: " .. received_data
                        conn.write(response_data) -- Send a response back
                        print("Server: Sent to " .. conn.remote_addr .. ": '" .. response_data .. "'")
                    else
                        -- read() returned nil, meaning connection was closed by client or error
                        print("Server: Connection with " .. conn.remote_addr .. " closed or lost.")
                        break -- Exit the loop for this client
                    end
                end
                conn.close() -- Ensure connection is closed from server side too
                print("Server: Handler goroutine finished for client " .. conn.remote_addr)
            end, client_conn) -- Pass the client_conn to the new goroutine
        else
            print("Server: Listener channel closed or error. Acceptor goroutine exiting.")
            break -- Exit if the listener channel itself has an issue
        end
    end
end)

-- Goroutine to act as a TCP client
go(function()
    sleep(1) -- Give the server a moment to fully start up
    
    print("Client: Attempting to connect to localhost:" .. tcp_port)
    local server_conn = tcp_connect("localhost", tcp_port)
    if not server_conn then
        print("Client: Failed to connect to the server.")
        return
    end
    print("Client: Successfully connected to the server.")
    
    local messages_to_send = {"Hello Server, this is Client!", "How are you doing today?", "Final message."}
    for _, msg in ipairs(messages_to_send) do
        print("Client: Sending: '" .. msg .. "'")
        server_conn.write(msg)
        local server_response = server_conn.read()
        if server_response then
            print("Client: Received from server: '" .. server_response .. "'")
        else
            print("Client: Did not receive response or server closed connection.")
            break
        end
        sleep(0.2) -- Small delay between messages
    end
    
    server_conn.close()
    print("Client: Connection closed. Client goroutine finished.")
end)

-- Keep the main script alive so the server and client can run
print("TCP demonstration setup complete. Running for 10 seconds...")
sleep(10)
-- close(server_listener_channel) -- If a close function exists for the listener to gracefully shutdown.
wait()
print("TCP demonstration finished.")
```

#### Fast Datagrams: UDP Communication

For connectionless, datagram-based communication, SolVM supports UDP (User Datagram Protocol). UDP is often used for applications where speed is preferred over guaranteed delivery, like gaming or streaming.

To **receive UDP packets**, you use `udp_recvfrom(port)`. This function binds to the specified `port` on all available interfaces and returns a `udp_socket` object. This object has a `receive()` method. Calling `udp_socket.receive()` blocks until a UDP packet arrives on that port. It then returns a table containing the received `data` (as a string), the sender's IP `addr` (address), and the sender's `port`.

To **send UDP packets**, you use `udp_sendto(destination_addr, destination_port, data_payload)`. You specify the `destination_addr` (IP address or hostname, or a broadcast address like "255.255.255.255"), the `destination_port` on the remote machine, and the `data_payload` string you want to send.

Since UDP is connectionless, you typically have one goroutine listening for incoming packets and potentially other goroutines sending packets as needed.

**UDP Sender and Receiver Example:**
```lua
local udp_listen_port = 8089
print("UDP communication test started.")

-- Create a UDP socket to receive packets on the specified port
print("Setting up UDP receiver on port " .. udp_listen_port)
local udp_socket = udp_recvfrom(udp_listen_port)
if not udp_socket then
    print("Failed to create UDP receiver socket.")
    return
end
print("UDP receiver is listening on port " .. udp_listen_port)

-- Goroutine to handle incoming UDP messages
go(function()
    print("UDP Receiver: Goroutine started, waiting for messages.")
    while true do
        local incoming_message = udp_socket.receive() -- Blocks until a packet is received
        if incoming_message then
            print(string.format("UDP Receiver: Got message from %s:%d - Data: '%s'",
                incoming_message.addr, incoming_message.port, incoming_message.data))
            
            -- Example: Send an echo response back to the sender
            local echo_response = "UDP Echo: " .. incoming_message.data
            print(string.format("UDP Receiver: Sending echo '%s' to %s:%d",
                echo_response, incoming_message.addr, incoming_message.port))
            udp_sendto(incoming_message.addr, incoming_message.port, echo_response)
        else
            -- This part might not be reached if receive() always blocks or errors out.
            -- Behavior depends on SolVM's implementation if the socket is closed.
            print("UDP Receiver: receive() returned nil. Exiting listener goroutine.")
            break
        end
    end
end)

-- Goroutine to send some test UDP messages
go(function()
    sleep(1) -- Wait a bit for the receiver goroutine to be ready
    
    local target_host = "localhost" -- Send to ourself for this test
    local target_port = udp_listen_port
    
    local messages = {
        "Hello from UDP sender!",
        "This is a test packet.",
        "SolVM UDP is working!"
    }
    
    for _, msg_payload in ipairs(messages) do
        print(string.format("UDP Sender: Sending '%s' to %s:%d", msg_payload, target_host, target_port))
        udp_sendto(target_host, target_port, msg_payload)
        sleep(0.5) -- Wait a bit to see responses if any
    end

    -- Example of sending a broadcast message (network configuration permitting)
    -- print(string.format("UDP Sender: Sending broadcast message to 255.255.255.255:%d", target_port))
    -- udp_sendto("255.255.255.255", target_port, "This is a broadcast message!")
    
    print("UDP Sender: All test messages sent. Sender goroutine finished.")
end)

-- Keep the script running to allow UDP communication
print("UDP sender and receiver are active. Running for 10 seconds...")
sleep(10)
-- udp_socket.close() -- If a close method exists for the UDP socket to unbind the port.
wait()
print("UDP communication test finished.")
```

### Managing Files and Directories

SolVM provides functions for interacting with the file system, often grouped conceptually under a file system module (even if not explicitly namespaced as `fs.` in the examples).

To **read the entire content of a file** into a string, you use `read_file(filepath)`. You provide the `filepath` to the desired file, and it returns the file's content. If the file doesn't exist or cannot be read, it will likely return `nil` or raise an error.

To **write data to a file**, overwriting it if it exists or creating it if it doesn't, you use `write_file(filepath, content)`. The `filepath` specifies the target file, and `content` is the string data to be written.

For **listing the contents of a directory**, SolVM offers `list_dir(directory_path)`. This function takes the `directory_path` and returns a table (an array) of items found within that directory. Each item in the table is itself a table, typically containing properties like `name` (the file or subdirectory name), `size` (in bytes, for files), and `is_dir` (a boolean indicating if it's a directory).

**File System Operations Example:**
```lua
local config_file_path = "my_app_config.json"
local output_file_path = "generated_output.txt"
local current_directory = "."

-- Write a sample config file (using json_encode for structure)
local sample_config_data = { settingA = "value1", settingB = 123, enabled = true }
local config_json_content = json_encode(sample_config_data) -- Assuming json_encode
print("Writing sample config to:", config_file_path)
write_file(config_file_path, config_json_content)

-- Read the config file
print("Reading config file:", config_file_path)
local read_config_content = read_file(config_file_path)
if read_config_content then
    print("Content of config file:")
    print(read_config_content)
    -- You would typically parse this, e.g., using json_decode(read_config_content)
    local parsed_data = json_decode(read_config_content)
    print("Parsed settingA from config:", parsed_data.settingA)
else
    print("Failed to read config file:", config_file_path)
end

-- Write to an output file
local text_to_write = "This is a line of text generated by SolVM script at " .. os.date() -- Assuming os.date()
print("Writing to output file:", output_file_path)
write_file(output_file_path, text_to_write)
print("Content written. Check '"..output_file_path.."'")

-- List contents of the current directory
print("\nListing contents of directory:", current_directory)
local directory_items = list_dir(current_directory)
if directory_items then
    print("Found " .. #directory_items .. " items:")
    for _, item_info in ipairs(directory_items) do
        local item_type = item_info.is_dir and "Directory" or "File"
        local item_details = string.format("  - Name: %-30s Type: %-10s", item_info.name, item_type)
        if not item_info.is_dir then
            item_details = item_details .. string.format(" Size: %d bytes", item_info.size)
        end
        print(item_details)
    end
else
    print("Failed to list directory contents for:", current_directory)
end
```

### Dynamic Development: Hot Reloading

SolVM supports a dynamic development workflow through its hot-reloading capabilities. This allows you to modify your script's code and have the changes applied almost instantly without restarting the entire application.

The core of this feature lies in two functions: `watch_file(filepath, callback_function)` and `reload_script()`.
You use `watch_file` to monitor your main script file (or any relevant configuration files) for modifications. When a change is detected in the `filepath`, the provided `callback_function` is executed.
Inside this callback, you typically call `reload_script()`. This function instructs SolVM to discard the current state of the script (though some global state or external connections might persist or need careful handling) and re-execute it from the beginning using the newly modified file content.

Additionally, for tasks that need to run periodically, SolVM provides `set_interval(callback_function, interval_seconds)`. This function repeatedly executes the `callback_function` every `interval_seconds`. This is useful for tasks like periodic state logging, health checks, or refreshing data, and it can be combined with hot reloading to update the behavior of these periodic tasks when the script is reloaded.

**Hot Reloading in Action:**
```lua
-- This script demonstrates hot-reloading.
-- To test, save this as 'hot_reload_example.lua', run it, then modify and save the file.
print("Hot-reload demonstration script initiated. Version 1.0")

local app_config = {
    greeting_message = "Hello from SolVM - Initial Version!",
    update_counter = 0
}

-- Function to display the current application state
function display_current_state()
    print("\n----- Current Application State -----")
    print("  Greeting:", app_config.greeting_message)
    print("  Update Counter:", app_config.update_counter)
    print("-----------------------------------")
    app_config.update_counter = app_config.update_counter + 1
end

-- Watch this script file for changes
-- IMPORTANT: The filename here should match the actual name of this script file.
local this_script_filename = "hot_reload_example.lua" -- Adjust if your filename is different
watch_file(this_script_filename, function()
    print("\n*********************************************")
    print("'"..this_script_filename.."' has been modified. Reloading script...")
    print("*********************************************\n")
    reload_script() -- The magic happens here!
end)
print("Now watching '"..this_script_filename.."' for changes.")

-- Set up a periodic task to print the state every 3 seconds
set_interval(function()
    display_current_state()
end, 3) -- Interval in seconds

print("Script is running. Modify and save this file to trigger a hot reload.")
print("The 'Update Counter' will reset upon reload, but 'Greeting' might change if you edit it.")

-- Keep the script alive to observe hot-reloading and interval timer
while true do
    sleep(1)
    -- The main loop here just keeps the script from exiting.
    -- All action happens via watch_file and set_interval.
end
```

### Structuring Your Code: Importing Modules

As your SolVM projects grow, organizing code into reusable modules becomes essential. SolVM supports importing modules in several ways using the `import("module_name")` function:

1. Individual Lua files: `module_name.lua`
2. Module folders: `module_name/`
3. ZIP archives: `module_name.zip` (local or remote)

When you call `import("module_name")`, SolVM looks for:
1. A file named `module_name.lua` in the modules directory
2. A folder named `module_name/` containing multiple Lua files
3. A ZIP file named `module_name.zip` (local or remote URL)

For individual files, it executes the file and returns the value (typically a table containing functions and data). For folders and ZIP files, it creates a namespace containing all the modules.

#### Defining Module Metadata

SolVM modules can declare metadata about themselves using the `metadata()` function, typically called at the beginning of the module file. This function takes a table as an argument, allowing you to specify various properties for the module.

The available metadata fields are:
*   `name` (string): The official name of the module.
*   `version` (string): The version of the module (e.g., "1.0.0", "0.2.1-beta").
*   `author` (string): The name or organization of the module's author(s).
*   `description` (string): A brief description of what the module does.
*   `repository` (string): A URL to the module's source code repository (e.g., a GitHub link).
*   `license` (string): The license under which the module is distributed (e.g., "MIT", "Apache-2.0").
*   `dependencies` (table, optional): A table of strings, where each string specifies a dependency on another module and its version requirement (e.g., `{"other_module >= 1.0.0", "shared_lib == 2.1.x"}`).

**General Metadata Syntax:**
```lua
metadata({
    name = "your_module_name",
    version = "1.0.0",
    author = "Your Name or Organization",
    description = "A concise description of what this module provides.",
    repository = "https://github.com/yourusername/your-repo",
    license = "MIT",
    dependencies = {
        "another_solvm_module >= 1.2.0",
        "utility_library == 0.5.x"
    }
})

-- ... rest of your module code ...
```
This metadata can be used by SolVM or associated tooling for package management, documentation generation, or dependency resolution in the future.

**Using Individual Module Imports:**
```lua
-- Assume math_utils.lua contains the module definition shown in Section 3
local math_lib = import("math_utils")

if not math_lib then
    print("Error: Could not import 'math_utils' module.")
    return
end

print("Successfully imported 'math_utils' module.")
-- If math_utils.lua defined metadata, SolVM might make it accessible, e.g. math_lib.metadata

local num1 = 15
local num2 = 7

local sum_result = math_lib.add(num1, num2)
print(string.format("%d + %d = %d", num1, num2, sum_result))
```

**Using Folder Imports:**
```lua
-- Assume a 'modules/utils/' folder exists with string.lua, math.lua, etc.
-- e.g., modules/utils/string.lua:
-- metadata({ name = "utils.string", version = "1.0" })
-- local M = {}
-- function M.trim(s) return s:match("^%s*(.-)%s*$") end
-- return M

import("utils/") -- Assuming this looks in a 'modules' directory or relative path

-- Access modules through the folder namespace
-- The exact access method might depend on SolVM's import implementation for folders.
-- If 'utils/' is imported as a table named 'utils':
if utils and utils.string then
    local str = utils.string.trim("  hello  ")
    print("Trimmed string:", str)
end
if utils and utils.math then
    local sum = utils.math.add(1, 2) -- Assuming utils.math.add exists
    print("Sum from utils.math:", sum)
end
-- If utils.array.map exists:
-- local arr = utils.array.map({1, 2, 3}, function(x) return x * 2 end)
```

**Using ZIP Imports:**
```lua
-- Import from a remote ZIP file
-- Assuming the ZIP contains modules that might also have metadata
-- import("https://example.com/modules.zip")

-- Access modules from the remote ZIP (namespace might be 'modules' or based on ZIP name)
-- local math_zip = modules.math.add(1, 2)
-- local str_zip = modules.string.trim("  hello  ")

-- Import from a local ZIP file
-- import("local_modules.zip")

-- Access modules from the local ZIP
-- local utils_zip = local_modules.utils.format("test")
```

**Module File Example (math.lua):**
(This example will be expanded in Section 3 to show `metadata` in use)
```lua
-- File: math_utils.lua (content shown in Section 3)
-- It would start with a metadata() call.

local math = {}

function math.add(a, b)
    return a + b
end

-- ... other functions ...

return math
```

**Using JSON with Imported Modules (Conceptual):**
```lua
-- Assuming 'utils/math.lua' provides subtract and is imported via import("utils/")
local calculation_task = {
    type = "subtraction",
    operand1 = 50,
    operand2 = 18
}

local task_json = json_encode(calculation_task)
print("Calculation task (JSON):", task_json)

local decoded_task = json_decode(task_json)
if decoded_task and decoded_task.type == "subtraction" and utils and utils.math and utils.math.subtract then
    local subtraction_result = utils.math.subtract(decoded_task.operand1, decoded_task.operand2)
    print(string.format("%d - %d = %d", decoded_task.operand1, decoded_task.operand2, subtraction_result))
else
    print("Could not perform subtraction from decoded task. Check module import and function availability.")
end
```

### Working with Various Data Formats

SolVM provides built-in support for parsing and encoding several common data formats, simplifying data interchange and configuration management.

#### JSON Handling (`json`)

JavaScript Object Notation (JSON) is a widely used lightweight data-interchange format. SolVM includes functions to work with JSON:
*   `json_encode(lua_table)`: Converts a Lua table into a JSON string.
*   `json_decode(json_string)`: Parses a JSON string and converts it back into a Lua table.

**JSON Example:**
```lua
local my_data_structure = {
    userName = "SolVM_User",
    score = 9876,
    isActive = true,
    tags = {"lua", "scripting", "vm"}
}

-- Encode Lua table to JSON string
local json_output_string = json_encode(my_data_structure)
print("Encoded JSON String:", json_output_string)

-- Decode JSON string back to Lua table
local decoded_lua_table = json_decode(json_output_string)
if decoded_lua_table then
    print("Decoded User Name:", decoded_lua_table.userName)
    print("Decoded Tags (first):", decoded_lua_table.tags[1])
else
    print("Failed to decode JSON string.")
end
```

#### Universally Unique Identifiers (`uuid`)

The `uuid` module helps in generating and validating UUIDs, which are useful for creating unique identifiers across systems.
*   `uuid.v4()`: Generates a standard Version 4 UUID (randomly generated) including hyphens.
*   `uuid.v4_without_hyphens()`: Generates a Version 4 UUID but omits the hyphens.
*   `uuid.is_valid(uuid_string)`: Checks if the given `uuid_string` is a valid UUID format and returns `true` or `false`.

#### Generating Random Data (`random`)

The `random` module provides utilities for generating various types of random data beyond simple numbers.
*   `random.number()`: Generates a random floating-point number, typically between 0.0 and 1.0.
*   `random.int(min_val, max_val)`: Generates a random integer between `min_val` and `max_val` (inclusive).
*   `random.string(length)`: Generates a random alphanumeric string of the specified `length`.

#### TOML Configuration Files (`toml`)

Tom's Obvious, Minimal Language (TOML) is a configuration file format designed to be easy to read.
*   `toml.encode(lua_table)`: Converts a Lua table into a TOML formatted string.
*   `toml.decode(toml_string)`: Parses a TOML string into a Lua table.

#### YAML Data Serialization (`yaml`)

YAML (YAML Ain't Markup Language) is another human-readable data serialization standard.
*   `yaml.encode(lua_table)`: Converts a Lua table into a YAML formatted string.
*   `yaml.decode(yaml_string)`: Parses a YAML string into a Lua table.

#### JSON with Comments (`jsonc`)

JSONC is a superset of JSON that allows for comments within the data structure, making configuration files more descriptive.
*   `jsonc.decode(jsonc_string)`: Parses a JSONC string (stripping comments) into a Lua table.
*   `jsonc.encode(lua_table)`: Converts a Lua table into a standard JSON string (comments are not added during encoding).

**Data Format Modules Showcase (UUID, Random, TOML, YAML, JSONC):**
```lua
-- UUID Module
local new_uuid_std = uuid.v4()
local new_uuid_compact = uuid.v4_without_hyphens()
print("Generated UUID v4 (standard):", new_uuid_std)
print("Generated UUID v4 (compact):", new_uuid_compact)
print("Is '"..new_uuid_std.."' a valid UUID?", uuid.is_valid(new_uuid_std))
print("Is 'not-a-uuid' a valid UUID?", uuid.is_valid("not-a-uuid"))

-- Random Module
local rand_float = random.number()
local rand_integer = random.int(10, 20)
local rand_alpha_str = random.string(12)
print("\nRandom Float (0-1):", rand_float)
print("Random Integer (10-20):", rand_integer)
print("Random Alphanumeric String (length 12):", rand_alpha_str)

-- TOML Module
local toml_config_data = {
    application_name = "SolVM Suite",
    version_info = { major = 1, minor = 2, patch = 0 },
    features_enabled = { "concurrency", "networking", "crypto" }
}
local toml_output_string = toml.encode(toml_config_data)
print("\nTOML Encoded Output:\n" .. toml_output_string)
local decoded_from_toml = toml.decode(toml_output_string)
print("Decoded TOML - Application Name:", decoded_from_toml.application_name)

-- YAML Module
local yaml_settings_data = {
    user = "test_user",
    preferences = { theme = "dark", notifications = false },
    roles = { "editor", "viewer" }
}
local yaml_output_string = yaml.encode(yaml_settings_data)
print("\nYAML Encoded Output:\n" .. yaml_output_string)
local decoded_from_yaml = yaml.decode(yaml_output_string)
print("Decoded YAML - User Theme:", decoded_from_yaml.preferences.theme)

-- JSONC Module
local jsonc_input_string = [[
{
    // This is a JSONC example with comments
    "projectName": "SolVM Project X", // Project identifier
    "debugMode": true, /* Multi-line
                           comment */
    "port": 8080
}
]]
local decoded_from_jsonc = jsonc.decode(jsonc_input_string)
print("\nDecoded JSONC - Project Name:", decoded_from_jsonc.projectName)
print("Decoded JSONC - Debug Mode:", decoded_from_jsonc.debugMode)
local re_encoded_json = jsonc.encode(decoded_from_jsonc) -- Encodes as standard JSON
print("Re-encoded from JSONC (as JSON):", re_encoded_json)
```

#### Comma-Separated Values (`csv`)

The `csv` module provides functions for working with CSV files, a common plain text format for tabular data.
*   `csv.write(filepath, data_table)`: Writes data from a Lua `data_table` (typically a table of tables, where the first inner table can be headers) to a CSV file specified by `filepath`.
*   `csv.read(filepath)`: Reads data from a CSV file at `filepath` and returns it as a Lua table of tables.

#### INI Configuration Files (`ini`)

The `ini` module handles INI files, a simple text-based format for configuration parameters, often organized into sections.
*   `ini.write(filepath, config_table)`: Writes a Lua `config_table` (where keys can be section names, and their values are tables of key-value pairs) to an INI file at `filepath`.
*   `ini.read(filepath)`: Reads an INI file from `filepath` and parses it into a Lua table structure.

**CSV and INI Modules Example:**
```lua
-- CSV Module
local employee_data = {
    {"ID", "Name", "Department", "Salary"},
    {"E101", "Alice Wonderland", "Engineering", "85000"},
    {"E102", "Bob The Builder", "Operations", "72000"},
    {"E103", "Charlie Brown", "HR", "68000"}
}
local csv_filename = "employee_records.csv"
csv.write(csv_filename, employee_data)
print(string.format("\nCSV data written to '%s'", csv_filename))

local loaded_csv_data = csv.read(csv_filename)
if loaded_csv_data and loaded_csv_data[3] then
    print("Loaded CSV Data (Employee 2 Name):", loaded_csv_data[3][2]) -- Accessing "Bob The Builder"
else
    print("Failed to load or parse CSV data correctly.")
end


-- INI Module
local server_configuration = {
    main_server = {
        host = "app.example.com",
        port = 443,
        protocol = "https"
    },
    database_connection = {
        type = "postgresql",
        user = "app_user",
        timeout_seconds = 30
    }
}
local ini_filename = "server_setup.ini"
ini.write(ini_filename, server_configuration)
print(string.format("\nINI configuration written to '%s'", ini_filename))

local loaded_ini_config = ini.read(ini_filename)
if loaded_ini_config and loaded_ini_config.main_server then
    print("Loaded INI - Main Server Host:", loaded_ini_config.main_server.host)
    print("Loaded INI - Database Timeout:", loaded_ini_config.database_connection.timeout_seconds)
else
    print("Failed to load or parse INI data correctly.")
end
```

### Advanced Text Manipulation (`text`)

The `text` module offers a suite of functions for common string manipulation tasks, extending Lua's built-in string library.
*   `text.trim(str)`: Removes leading and trailing whitespace from `str`.
*   `text.lower(str)`: Converts `str` to lowercase.
*   `text.upper(str)`: Converts `str` to uppercase.
*   `text.title(str)`: Converts `str` to title case (e.g., "hello world" becomes "Hello World").
*   `text.split(str, separator)`: Splits `str` into a table of substrings using `separator`.
*   `text.join(table_of_strings, separator)`: Joins elements of `table_of_strings` into a single string using `separator`.
*   `text.replace(str, old_substring, new_substring)`: Replaces all occurrences of `old_substring` with `new_substring` in `str`.
*   `text.contains(str, substring)`: Returns `true` if `str` contains `substring`, `false` otherwise.
*   `text.starts_with(str, prefix)`: Returns `true` if `str` starts with `prefix`.
*   `text.ends_with(str, suffix)`: Returns `true` if `str` ends with `suffix`.
*   `text.pad_left(str, length, pad_char)`: Pads `str` on the left with `pad_char` until it reaches `length`.
*   `text.repeat_str(str, count)`: Repeats `str` `count` times.

**Text Utilities Showcase:**
```lua
local sample_phrase = "   SolVM Text Utilities Demo!   "
print("\nOriginal Phrase: '" .. sample_phrase .. "'")

print("Trimmed: '" .. text.trim(sample_phrase) .. "'")
print("Lowercase: '" .. text.lower(sample_phrase) .. "'")
print("Uppercase: '" .. text.upper(sample_phrase) .. "'")
print("Title Case: '" .. text.title(text.trim(sample_phrase)) .. "'") -- Trim first for better title casing

local words_array = text.split(text.trim(sample_phrase), " ")
print("Split into words (joined with ','): " .. text.join(words_array, ", "))

local replaced_phrase = text.replace(sample_phrase, "Demo", "Showcase")
print("Replaced 'Demo' with 'Showcase': '" .. replaced_phrase .. "'")

print("Contains 'SolVM':", text.contains(sample_phrase, "SolVM"))
print("Starts with '   SolVM':", text.starts_with(sample_phrase, "   SolVM"))
print("Ends with 'Demo!   ':", text.ends_with(sample_phrase, "Demo!   "))

local padded_number = text.pad_left("7", 3, "0")
print("Padded '7' to length 3 with '0': '" .. padded_number .. "'") -- "007"

local repeated_pattern = text.repeat_str("xo", 4)
print("Repeated 'xo' 4 times: '" .. repeated_pattern .. "'") -- "xoxoxoxo"
```

### Handling Environment Variables (`dotenv`)

The `dotenv` module facilitates loading configuration variables from a `.env` file into the script's environment. This is a common practice for managing application settings, especially secrets, separately from code.
*   `dotenv.load(filepath)`: Loads variables from the specified `filepath` (typically ".env" in the script's root directory).
*   `dotenv.get(variable_name, default_value)`: Retrieves the value of an environment variable named `variable_name`. If the variable is not set, it returns the `default_value`.

**`.env` File Example (not Lua, just content of a file named `.env`):**
```
API_SECRET_KEY=yourActualSecretKeyGoesHere
DATABASE_URL=postgres://user:pass@host:port/dbname
FEATURE_FLAG_X=true
```

**Dotenv Module Usage:**
```lua
-- Assume a .env file exists with:
-- API_KEY=mysecret123abc
-- DEBUG_MODE=true

-- Load variables from '.env' file (if it exists in the current directory)
dotenv.load(".env") -- Or specify a path: dotenv.load("/path/to/your/.env")
print("\nLoaded .env file (if present).")

local retrieved_api_key = dotenv.get("API_KEY", "fallback_default_key")
print("Retrieved API Key:", retrieved_api_key)

local non_existent_var = dotenv.get("NON_EXISTENT_VARIABLE", "this_is_the_default")
print("Retrieved Non-Existent Variable (shows default):", non_existent_var)

local debug_setting_str = dotenv.get("DEBUG_MODE", "false")
-- Environment variables are typically strings, so convert if needed
local is_debug_mode = (debug_setting_str == "true")
print("Is Debug Mode Active (from .env):", is_debug_mode)
```

### Precise Date and Time Operations (`datetime`)

The `datetime` module provides robust functions for working with dates and times.
*   `datetime.now()`: Returns the current timestamp (often as a Unix timestamp or a specific SolVM datetime object).
*   `datetime.format(timestamp, layout_string)`: Formats the given `timestamp` into a human-readable string according to the `layout_string`. The layout string often follows Go's time formatting conventions (e.g., "2006-01-02 15:04:05").
*   `datetime.add(timestamp, duration_string)`: Adds a duration (specified by `duration_string`, e.g., "24h", "5m", "10s") to the given `timestamp` and returns the new timestamp.
*   `datetime.diff(timestamp1, timestamp2)`: Calculates the difference between two timestamps, returning a duration (which might be a string representation or a numerical value in seconds).

**Datetime Module Example:**
```lua
print("\nDatetime module demonstration:")
local current_time_ts = datetime.now()
print("Current Timestamp (raw):", current_time_ts) -- The raw format depends on SolVM's implementation

-- Format the current time
-- The format string "2006-01-02 15:04:05" is a common Go-style reference date.
local formatted_current_time = datetime.format(current_time_ts, "2006-01-02 15:04:05 MST")
print("Formatted Current Time:", formatted_current_time)

-- Add 3 hours and 30 minutes to the current time
local future_time_ts = datetime.add(current_time_ts, "3h30m")
local formatted_future_time = datetime.format(future_time_ts, "Mon, 02 Jan 2006 'at' 03:04 PM")
print("Time in 3h 30m:", formatted_future_time)

-- Calculate the difference between now and the future time
local time_difference = datetime.diff(current_time_ts, future_time_ts)
print("Difference between now and future time:", time_difference) -- Output format of diff depends on SolVM
```

### Streamlining File Transfers (`ft`)

The `ft` (file transfer) module appears to offer utilities for common file operations that involve moving or copying files, potentially including remote transfers.
*   `ft.download(url, destination_path)`: Downloads a file from the given `url` and saves it to `destination_path`.
*   `ft.copy(source_path, destination_path)`: Copies a file from `source_path` to `destination_path`.
*   `ft.move(source_path, destination_path)`: Moves (renames) a file from `source_path` to `destination_path`.

**File Transfer Module Example:**
*(This example assumes `ft.download` works with a public URL and that directories like `archive` exist or are created by the functions. Error handling would be crucial in a real script.)*
```lua
print("\nFile Transfer (ft) module demonstration:")

-- For this example to run, you'd need a live URL for a small file.
-- Let's use placeholder paths for local operations.
-- os.execute("mkdir -p archive") -- Ensure 'archive' directory exists (using os module)
-- write_file("source_file_for_ft.txt", "This is a test file for ft operations.") -- Create a source file

-- Placeholder for download
-- print("Attempting to download a file...")
-- ft.download("https://www.example.com/robots.txt", "downloaded_robots.txt")
-- print("Download attempted. Check for 'downloaded_robots.txt'.")

-- Copy a local file
print("\nCopying 'source_file_for_ft.txt' to 'copied_file.txt'...")
-- write_file("source_file_for_ft.txt", "Content for copy and move.") -- Create if not present
-- ft.copy("source_file_for_ft.txt", "copied_file.txt")
-- print("Copy operation complete (simulated).")

-- Move a local file
print("\nMoving 'copied_file.txt' to 'archive/moved_file.txt'...")
-- os.execute("mkdir -p archive") -- Ensure archive dir exists
-- ft.move("copied_file.txt", "archive/moved_file.txt")
-- print("Move operation complete (simulated).")

-- print("\nCheck your file system for 'downloaded_robots.txt', 'copied_file.txt' (if not moved), and 'archive/moved_file.txt'.")
-- print("Note: Actual download requires a valid URL and network access.")
-- print("Note: For local copy/move to work, source files and destination directories must be valid.")
-- Since the provided snippet for `ft` doesn't include file creation,
-- this part is more conceptual. Actual usage would depend on `write_file` or `os.execute`
-- to set up the scenario if files aren't downloaded.
```

### Creating and Managing Archives (`tar`)

The `tar` module provides functionalities for working with Tape Archive (.tar) files, which are commonly used for bundling multiple files and directories into a single archive file. It also supports GZip compression for these archives.

*   `tar.create(archive_filepath, source_path_or_table, [compress])`: Creates a TAR archive.
    *   `archive_filepath`: The path where the .tar (or .tar.gz) file will be created (e.g., "my_backup.tar" or "my_backup.tar.gz").
    *   `source_path_or_table`: Either a single string path to a file or directory to be archived, or a table of string paths.
    *   `compress` (boolean, optional): If `true`, the resulting archive will be GZip compressed (typically resulting in a .tar.gz extension). Defaults to `false`.
*   `tar.list(archive_filepath)`: Lists the contents of a TAR archive specified by `archive_filepath`. It returns a table of items, where each item is a table detailing a file or directory within the archive (e.g., with `name`, `size`, `type` properties).
*   `tar.extract(archive_filepath, destination_directory)`: Extracts the contents of the `archive_filepath` into the `destination_directory`. The `destination_directory` will be created if it doesn't exist.

**TAR Module Operations Example:**
```lua
print("\nTAR module demonstration:")

-- Setup: Create a directory structure to archive
os.execute("rm -rf test_archive_dir extracted_archive_dir_uncompressed extracted_archive_dir_compressed my_archive.tar my_archive_c.tar.gz") -- Clean up
os.execute("mkdir -p test_archive_dir/inner_subdir")
write_file("test_archive_dir/main_file.txt", "This is the main file content.")
write_file("test_archive_dir/inner_subdir/nested_file.log", "Log entry 1\nLog entry 2")
print("Created 'test_archive_dir' with some files for archiving.")

local uncompressed_tar_path = "my_archive.tar"
local compressed_tar_path = "my_archive_c.tar.gz"
local extraction_path_uncompressed = "extracted_archive_dir_uncompressed"
local extraction_path_compressed = "extracted_archive_dir_compressed"

-- Create an uncompressed TAR archive
print(string.format("\nCreating uncompressed archive '%s' from 'test_archive_dir'...", uncompressed_tar_path))
tar.create(uncompressed_tar_path, "test_archive_dir")
print("Uncompressed archive created.")

-- Create a compressed TAR archive (tar.gz)
print(string.format("\nCreating compressed archive '%s' from 'test_archive_dir'...", compressed_tar_path))
tar.create(compressed_tar_path, "test_archive_dir", true) -- true for compression
print("Compressed archive created.")

-- List contents of the uncompressed archive
print(string.format("\nListing contents of '%s':", uncompressed_tar_path))
local uncompressed_files_list = tar.list(uncompressed_tar_path)
if uncompressed_files_list then
    for _, file_info in ipairs(uncompressed_files_list) do
        print(string.format("  - %s (Type: %s, Size: %s)", file_info.name, file_info.type, tostring(file_info.size)))
    end
else
    print("Could not list archive contents.")
end

-- Extract the uncompressed archive
os.execute("rm -rf " .. extraction_path_uncompressed) -- Clean extraction dir
print(string.format("\nExtracting '%s' to '%s'...", uncompressed_tar_path, extraction_path_uncompressed))
tar.extract(uncompressed_tar_path, extraction_path_uncompressed)
print("Extraction complete. Check the '"..extraction_path_uncompressed.."' directory.")

-- List contents of the compressed archive
print(string.format("\nListing contents of compressed archive '%s':", compressed_tar_path))
local compressed_files_list = tar.list(compressed_tar_path)
if compressed_files_list then
    for _, file_info in ipairs(compressed_files_list) do
        print(string.format("  - %s (Type: %s, Size: %s)", file_info.name, file_info.type, tostring(file_info.size)))
    end
else
    print("Could not list compressed archive contents.")
end

-- Extract the compressed archive
os.execute("rm -rf " .. extraction_path_compressed) -- Clean extraction dir
print(string.format("\nExtracting compressed archive '%s' to '%s'...", compressed_tar_path, extraction_path_compressed))
tar.extract(compressed_tar_path, extraction_path_compressed)
print("Compressed extraction complete. Check the '"..extraction_path_compressed.."' directory.")

-- Cleanup (optional)
-- os.execute("rm -rf test_archive_dir my_archive.tar my_archive_c.tar.gz " .. extraction_path_uncompressed .. " " .. extraction_path_compressed)
print("\nTAR module demonstration finished.")
```

### Interacting with the Operating System (`os`)

The `os` module provides a way to execute shell commands on the underlying operating system. This is a powerful feature but should be used with caution, especially with user-provided input, due to security implications.
*   `os.execute(command_string)`: Executes the `command_string` in the system's shell. The return value typically indicates success or failure (e.g., an exit code), but this can vary based on the OS and SolVM's implementation.

**OS Command Execution Example:**
```lua
print("\nOS command execution demonstration:")

-- Example: List files in the current directory (cross-platform might vary)
-- For Unix-like systems (Linux, macOS):
print("\nAttempting to list directory contents using 'ls -l' (Unix-like):")
local ls_result = os.execute("ls -l")
-- print("Result of 'ls -l':", ls_result) -- Output depends on SolVM's os.execute behavior

-- For Windows:
-- print("\nAttempting to list directory contents using 'dir' (Windows):")
-- local dir_result = os.execute("dir")
-- print("Result of 'dir':", dir_result)

-- Example: Create a directory
local new_dir_name = "temp_os_created_dir"
os.execute("rm -rf " .. new_dir_name) -- Clean up if exists from previous run
print("\nAttempting to create directory: " .. new_dir_name)
local mkdir_result = os.execute("mkdir " .. new_dir_name)
if mkdir_result then -- Assuming true or 0 for success
    print("Directory '" .. new_dir_name .. "' creation attempt successful (or command ran).")
    -- You might want to check if the directory actually exists using list_dir or another os.execute
else
    print("Directory '" .. new_dir_name .. "' creation attempt failed (or command indicated error).")
end

-- Example: Remove the created directory
print("\nAttempting to remove directory: " .. new_dir_name)
local rmdir_command = "rm -rf " .. new_dir_name -- Unix-like
-- For Windows, it might be: "rmdir /s /q " .. new_dir_name
local rmdir_result = os.execute(rmdir_command)
if rmdir_result then
    print("Directory '" .. new_dir_name .. "' removal attempt successful (or command ran).")
else
    print("Directory '" .. new_dir_name .. "' removal attempt failed (or command indicated error).")
end

print("\nWarning: os.execute can be dangerous if used with untrusted input.")
```

### Template Engine (`template`)

SolVM includes a powerful template engine, allowing you to generate dynamic text-based content (such as HTML, configuration files, emails, etc.) by embedding logic and data placeholders within template definitions. This engine is typically based on Go's `text/template` or `html/template` packages, offering a rich and familiar syntax.

The `template` module provides several functions to load and parse templates:

*   `template.parse(template_string)`: Parses a template directly from a provided string.
*   `template.parse_file(filepath)`: Parses a template from a single specified file.
*   `template.parse_files(filepath1, filepath2, ...)`: Parses one or more template files. This is useful for defining and using layouts, partials, or a collection of named templates. Typically, the templates can then refer to each other (e.g., a layout template including a content template).
*   `template.parse_glob(pattern)`: Parses all template files matching a given glob pattern (e.g., "templates/*.html").

Each of these parsing functions returns a compiled `template_object`. This object is then invoked as a function, passing it a Lua table containing the data to be rendered into the template. The `template_object` call returns the final rendered string.

Common template syntax features include:
*   Data access: `{{.FieldName}}` or `{{.MapKey}}` to access fields of the data table.
*   Control structures: `{{if .Condition}}...{{else if .AnotherCondition}}...{{else}}...{{end}}` for conditional rendering.
*   Range loops: `{{range .Items}}...{{.}} (or {{.ItemField}})...{{end}}` to iterate over arrays or slices.
*   Pipelines and functions: `{{.Value | printf "%.2f"}}`
*   Named template definitions and inclusions: `{{define "myPartial"}}...{{end}}` and `{{template "myPartial" .}}`.

**Template Engine Usage Example:**

```lua
print("\nTemplate Engine demonstration:")

-- Example 1: Parse template from string
local string_template_definition = [[
<!DOCTYPE html>
<html>
<head>
    <title>{{.pageTitle}}</title>
    <style> body { font-family: sans-serif; } </style>
</head>
<body>
    <h1>{{.pageTitle}}</h1>
    <p>{{.mainContent}}</p>
    {{if .items}}
    <h2>Items:</h2>
    <ul>
    {{range .items}}
        <li>{{.}}</li>
    {{end}}
    </ul>
    {{else}}
    <p>No items to display.</p>
    {{end}}
    <p><em>Rendered by SolVM Template Engine at {{.renderTime}}</em></p>
</body>
</html>
]]

local compiled_string_tmpl = template.parse(string_template_definition)

local page_data = {
    pageTitle = "SolVM Dynamic Page",
    mainContent = "Welcome to this page generated by SolVM's template engine!",
    items = {"Apple", "Banana", "Cherry"},
    renderTime = datetime.format(datetime.now(), "15:04:05 MST") -- Using datetime module
}

local rendered_html_from_string = compiled_string_tmpl(page_data)
print("\n--- Rendered HTML from String Template ---")
print(rendered_html_from_string)

-- Example 2: Parse template from file (conceptual)
-- Assuming a file 'examples/templates/page.html' exists with similar content:
-- write_file("examples/templates/page.html", string_template_definition) -- Create the file
-- local file_tmpl = template.parse_file("examples/templates/page.html")
-- local file_data = {
--     pageTitle = "From File",
--     mainContent = "This content comes from a file template",
--     items = {"File Item 1", "File Item 2"},
--     renderTime = datetime.format(datetime.now(), "HH:MM:SSZ")
-- }
-- local rendered_html_from_file = file_tmpl(file_data)
-- print("\n--- Rendered HTML from File Template ---")
-- print(rendered_html_from_file)

-- Example 3: Parse multiple template files (conceptual)
-- Assuming 'examples/templates/layout.html' and 'examples/templates/content.html' exist:
-- write_file("examples/templates/layout.html", "{{define \"layout\"}}Layout Start: {{template \"content\" .}} Layout End{{end}}")
-- write_file("examples/templates/content.html", "{{define \"content\"}}Page Title: {{.pageTitle}}, Content: {{.mainContent}}{{end}}")
-- local multi_tmpl = template.parse_files("examples/templates/layout.html", "examples/templates/content.html")
-- local multi_data = { pageTitle = "Multi-File", mainContent = "Content for multi-file example" }
-- -- To render a specific named template from the set, usually one is the entry point
-- -- or the engine needs to know which top-level template to execute.
-- -- If 'layout' is the main template to render:
-- -- local rendered_html_from_multi = multi_tmpl.execute_template(multi_data, "layout") -- Or similar invocation
-- -- For simplicity, SolVM's template object might directly render the first parsed or a specifically named one.
-- -- print("\n--- Rendered HTML from Multiple Files ---")
-- -- print(rendered_html_from_multi)


-- Example 4: Parse templates using glob pattern (conceptual)
-- Assuming HTML files exist in 'examples/templates/' directory:
-- os.execute("mkdir -p examples/templates_glob")
-- write_file("examples/templates_glob/glob_page1.html", "{{define \"glob1\"}}Glob Page 1: {{.message1}}{{end}}")
-- write_file("examples/templates_glob/glob_page2.html", "{{define \"glob2\"}}Glob Page 2: {{.message2}}{{end}}")
-- local glob_tmpl = template.parse_glob("examples/templates_glob/*.html")
-- local glob_data = { message1 = "Hello from Glob 1", message2 = "Greetings from Glob 2"}
-- -- local rendered_glob1 = glob_tmpl.execute_template(glob_data, "glob1")
-- -- print("\n--- Rendered HTML from Glob (Page 1) ---")
-- -- print(rendered_glob1)

print("\nTemplate Engine demonstration finished.")
-- os.execute("rm -rf examples") -- Cleanup conceptual files
```

### General Purpose Utilities (`utils`)

The `utils` module provides a collection of general-purpose utility functions that can be helpful in various scripting scenarios. These functions cover string manipulation, table operations, and environment interaction. To use them, you would typically import a module named `utils` (e.g., `local utils = import("utils")`).

*   `utils.split(str, separator)`: Splits a string `str` by a `separator` and returns a table of substrings.
*   `utils.join(table_of_strings, separator)`: Joins elements of `table_of_strings` into a single string, with each element separated by `separator`.
*   `utils.escape(path_string)`: Escapes special characters in a `path_string`, which can be useful for file paths or other strings needing system-safe representation.
*   `utils.unescape(escaped_string)`: Reverses the escaping done by `utils.escape`.
*   `utils.unpack(tbl, start_index, end_index)`: Similar to Lua's `table.unpack` (or `unpack` in older Lua versions), this extracts elements from table `tbl` from `start_index` to `end_index` and returns them as multiple values.
*   `utils.getfenv(level_or_func)`: Gets the environment (table of global variables) of a function. `level_or_func` can be a function or a stack level (integer).
*   `utils.setfenv(func, env_table)`: Sets the environment of the given `func` to `env_table`.

**General Utilities (`utils`) Showcase:**
```lua
-- Assuming 'utils' module is available and imported, e.g.:
-- local utils = import("utils") 
-- For demonstration, we'll define mock functions if 'utils' is not a built-in for this example context.
-- If 'utils' is a built-in module, these definitions are not needed.
local utils = {
    split = function(str, sep)
        local result = {}
        for part in string.gmatch(str, "([^"..sep.."]+)") do table.insert(result, part) end
        return result
    end,
    join = function(tbl, sep) return table.concat(tbl, sep) end,
    escape = function(path) return string.gsub(path, "([%-%[%]%^%$%.%*%+%?%(\\%)])", "%%%1") end, -- Basic example
    unescape = function(escaped_path) return string.gsub(escaped_path, "%%(.)", "%1") end, -- Basic example
    unpack = function(tbl, i, j) return table.unpack(tbl, i, j) end,
    getfenv = getfenv or function(f) return _G end, -- Fallback for Lua 5.1 or if not available
    setfenv = setfenv or function(f, env) print("setfenv mock: not fully implemented for this demo") end -- Fallback
}

-- Mock tablex for pretty printing if not available
local tablex = {
    pretty = function(tbl)
        local s = "{\n"
        for k, v in pairs(tbl) do
            s = s .. string.format("  [%s] = %s,\n", tostring(k), tostring(v))
        end
        return s .. "}"
    end
}


print("\n--- Utils Module Demonstration ---")

-- String operations
local str = "hello world from SolVM"
local parts = utils.split(str, " ")
print("Split string: " .. str)
print("Parts:", tablex.pretty(parts))

local joined = utils.join(parts, "-")
print("Joined with '-':", joined)

-- String escaping
local path = "C:\\Program Files\\My App\\file.txt"
local escaped = utils.escape(path)
print("Original path:", path)
print("Escaped path:", escaped)
local unescaped = utils.unescape(escaped)
print("Unescaped path:", unescaped)

-- Table operations
local tbl = {10, 20, 30, 40, 50}
print("Original table for unpack:", tablex.pretty(tbl))
local a, b, c = utils.unpack(tbl, 2, 4) -- Should get 20, 30, 40
print("Unpacked (elements 2 to 4):", a, b, c)

-- Environment operations
local test_env_func_original_env
local function test_env_func()
    print("  Inside test_env_func: x =", x, ", y =", y) -- x and y are expected from the new env
    test_env_func_original_env = utils.getfenv(1) -- Get this function's current environment
end

print("\nEnvironment operations:")
local new_env_for_func = { x = 42, y = 24, print = print } -- Must include print if used inside
utils.setfenv(test_env_func, new_env_for_func)

print("Calling test_env_func after setfenv:")
test_env_func()

print("Environment of test_env_func (retrieved from within):")
if test_env_func_original_env then
    print("  test_env_func_original_env.x = " .. tostring(test_env_func_original_env.x))
    print("  test_env_func_original_env.y = " .. tostring(test_env_func_original_env.y))
else
    print("  Could not retrieve function's environment details in this example setup.")
end
```

### Advanced Type System (`types`)

SolVM provides an advanced `types` module for more granular and robust type checking than Lua's built-in `type()` function. This module is particularly useful for validating data structures, function arguments, and ensuring type safety in complex scripts.

The `types` module offers the following functions:
*   `types.type(value)`: Returns a more specific string representation of the `value`'s type (e.g., distinguishes "integer" from "float" if SolVM supports it, or provides custom type names for SolVM objects).
*   `types.is_number(value)`: Returns `true` if `value` is a number (integer or float), `false` otherwise.
*   `types.is_integer(value)`: Returns `true` if `value` is an integer, `false` otherwise. (Behavior might depend on SolVM's underlying number representation).
*   `types.is_string(value)`: Returns `true` if `value` is a string.
*   `types.is_boolean(value)`: Returns `true` if `value` is a boolean.
*   `types.is_nil(value)`: Returns `true` if `value` is `nil`.
*   `types.is_table(value)`: Returns `true` if `value` is a table.
*   `types.is_function(value)`: Returns `true` if `value` is a Lua function.
*   `types.is_callable(value)`: Returns `true` if `value` is a function or a table with a `__call` metamethod, meaning it can be called like a function.

**Advanced Type System (`types`) Showcase:**
```lua
-- Assuming 'types' module is available and imported, e.g.:
-- local types = import("types")
-- For demonstration, we'll define mock functions if 'types' is not a built-in for this example context.
-- If 'types' is a built-in module, these definitions are not needed.
local types = {
    type = function(v) 
        local t = type(v)
        if t == "number" then
            if math.floor(v) == v then return "integer" else return "float" end
        end
        return t 
    end,
    is_number = function(v) return type(v) == "number" end,
    is_integer = function(v) return type(v) == "number" and math.floor(v) == v end,
    is_string = function(v) return type(v) == "string" end,
    is_boolean = function(v) return type(v) == "boolean" end,
    is_nil = function(v) return v == nil end,
    is_table = function(v) return type(v) == "table" end,
    is_function = function(v) return type(v) == "function" end,
    is_callable = function(v)
        return type(v) == "function" or (type(v) == "table" and getmetatable(v) and getmetatable(v).__call ~= nil)
    end
}

print("\n--- Types Module Demonstration ---")

local values_to_check = {
    42,           -- number (integer)
    42.5,         -- number (float)
    "hello SolVM",-- string
    true,         -- boolean
    nil,          -- nil
    {},           -- table
    function() print("I am a function") end -- function
}

print("Type checking examples for various values:")
for i, v_item in ipairs(values_to_check) do
    print(string.format("\nValue %d (%s):", i, tostring(v_item)))
    print("  Lua's type()     =", type(v_item))
    print("  types.type()     =", types.type(v_item))
    print("  types.is_number  =", types.is_number(v_item))
    print("  types.is_integer =", types.is_integer(v_item))
    print("  types.is_string  =", types.is_string(v_item))
    print("  types.is_boolean =", types.is_boolean(v_item))
    print("  types.is_nil     =", types.is_nil(v_item))
    print("  types.is_table   =", types.is_table(v_item))
    print("  types.is_function=", types.is_function(v_item))
    print("  types.is_callable=", types.is_callable(v_item))
end

-- Custom callable object example
local my_callable_object = {
    message_prefix = "Callable Object says: "
}
-- Set its metatable with a __call method
setmetatable(my_callable_object, {
    __call = function(self, ...)
        local args = {...}
        local message_to_print = self.message_prefix
        for _, arg_val in ipairs(args) do
            message_to_print = message_to_print .. tostring(arg_val) .. " "
        end
        print(message_to_print)
    end
})

print("\nCustom callable object:")
print("  Lua's type()     =", type(my_callable_object))          -- Expected: table
print("  types.type()     =", types.type(my_callable_object))    -- Expected: table (unless types.type has special handling)
print("  types.is_function=", types.is_function(my_callable_object)) -- Expected: false
print("  types.is_callable=", types.is_callable(my_callable_object)) -- Expected: true

print("Attempting to call the custom object:")
my_callable_object("Hello", "from", "SolVM callable!", 123)
```

---

## 3. Creating Your Own Modules: An Example (`math_utils.lua`)

SolVM allows you to structure your code into reusable modules. A module is typically a Lua file that returns a table containing functions and data. This promotes organization and code reuse. As described in [Defining Module Metadata](#defining-module-metadata), modules can also declare information about themselves.

Here's an example of a simple math utilities module, `math_utils.lua`, including metadata:

```lua
-- File: math_utils.lua

-- Define metadata for this module
metadata({
    name = "math_utils",
    version = "1.0.0",
    author = "SolVM Team",
    description = "Basic math operations module for SolVM.",
    repository = "https://github.com/kleeedolinux/SolVM", -- Example repository
    license = "MIT"
})

-- Create a table to hold our module's functions
local math_operations = {}

-- Define an addition function
function math_operations.add(a, b)
    if type(a) ~= "number" or type(b) ~= "number" then
        error("Invalid input: add expects two numbers.")
    end
    return a + b
end

-- Define a subtraction function
function math_operations.subtract(a, b)
    if type(a) ~= "number" or type(b) ~= "number" then
        error("Invalid input: subtract expects two numbers.")
    end
    return a - b
end

-- Define a multiplication function
function math_operations.multiply(a, b)
    if type(a) ~= "number" or type(b) ~= "number" then
        error("Invalid input: multiply expects two numbers.")
    end
    return a * b
end

-- Define a division function with a check for division by zero
function math_operations.divide(a, b)
    if type(a) ~= "number" or type(b) ~= "number" then
        error("Invalid input: divide expects two numbers.")
    end
    if b == 0 then
        error("Critical error: Division by zero is not allowed.")
    end
    return a / b
end

-- Return the table, making its functions available to scripts that import this module
return math_operations
```

You would then use this module in another SolVM script as shown in the "Importing Modules" section. For instance:
```lua
-- In another script:
local math_lib = import("math_utils")
if math_lib then
    print("Imported math_utils. Adding 5 and 3:", math_lib.add(5, 3))
    -- Depending on SolVM's implementation, metadata might be accessible:
    -- if math_lib.metadata then print("Module Version:", math_lib.metadata.version) end
else
    print("Failed to import math_utils.")
end
```
This modular approach is key to building larger, maintainable applications in SolVM.
