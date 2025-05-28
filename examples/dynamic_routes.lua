-- Create a basic HTTP server
create_server("main", 8080)

-- Authentication middleware
use_middleware("main", "/api/*", function(req)
    local token = req.headers["Authorization"]
    if not token then
        return false
    end
    return true
end)

-- Logging middleware
use_middleware("main", "/*", function(req)
    print("Request:", req.method, req.path)
    return true
end)

-- Basic route
handle_http("main", "/", function(req)
    return {
        status = 200,
        headers = {
            ["Content-Type"] = "text/plain"
        },
        body = "Welcome to SolVM Server!"
    }
end)

-- Dynamic route with parameters
handle_http("main", "/users/:id", function(req)
    local userId = req.params[1]
    return {
        status = 200,
        headers = {
            ["Content-Type"] = "application/json"
        },
        body = string.format('{"id": "%s", "name": "User %s"}', userId, userId)
    }
end)

-- API route with authentication
handle_http("main", "/api/data", function(req)
    return {
        status = 200,
        headers = {
            ["Content-Type"] = "application/json"
        },
        body = '{"message": "Protected data", "timestamp": "' .. os.time() .. '"}'
    }
end)

-- WebSocket example
handle_ws("main", "/ws", function(ws)
    while true do
        local message = ws:receive()
        if message then
            ws:send("Echo: " .. message)
        end
    end
end)

-- Start the server
start_server("main") 