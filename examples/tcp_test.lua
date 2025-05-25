-- Example of TCP server and client
print("TCP test script started")

-- Create TCP server
local server = tcp_listen(8080)
print("TCP server listening on port 8080")

-- Handle incoming connections
go(function()
    while true do
        local conn = receive(server)
        if conn then
            print("New connection from:", conn.remote_addr)
            
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
                print("Connection closed")
                conn.close()
            end)
        end
    end
end)

-- Create TCP client
go(function()
    sleep(1) -- Wait for server to start
    
    local conn = tcp_connect("localhost", 8080)
    print("Connected to server")
    
    -- Send some messages
    conn.write("Hello, Server!")
    local response = conn.read()
    print("Server response:", response)
    
    conn.write("How are you?")
    response = conn.read()
    print("Server response:", response)
    
    conn.close()
end)

-- Keep script running
while true do
    sleep(1)
end 