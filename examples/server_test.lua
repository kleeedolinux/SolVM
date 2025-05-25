-- Create HTTP server
create_server("http_server", 8080, false)

-- Handle HTTP requests
handle_http("http_server", "/", function(req)
    print("Received HTTP request:")
    print("  Method:", req.method)
    print("  Path:", req.path)
    print("  Query:", req.query)
    
    -- Return response
    return {
        status = 200,
        headers = {
            ["Content-Type"] = "application/json"
        },
        body = json_encode({
            message = "Hello from SolVM HTTP server!",
            request = {
                method = req.method,
                path = req.path,
                query = req.query
            }
        })
    }
end)

-- Handle WebSocket connections
handle_ws("http_server", "/ws", function(ws)
    print("New WebSocket connection")
    
    -- Send welcome message
    ws.send(json_encode({
        type = "welcome",
        message = "Connected to SolVM WebSocket server!"
    }))
    
    -- Echo messages
    while true do
        local message = ws.receive()
        if message == nil then
            break
        end
        
        print("Received WebSocket message:", message)
        ws.send(json_encode({
            type = "echo",
            message = message
        }))
    end
    
    print("WebSocket connection closed")
end)

-- Start the server
print("Starting server on port 8080...")
start_server("http_server")

-- Keep the script running
while true do
    sleep(1)
end 