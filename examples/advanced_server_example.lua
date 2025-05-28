-- Create HTTP and HTTPS servers
create_server("http", 8080)
create_server("https", 8443, true, "cert.pem", "key.pem")

-- Rate limiting middleware
local requestCounts = {}
use_middleware("http", "/*", function(req)
    local ip = req.headers["X-Forwarded-For"] or "unknown"
    requestCounts[ip] = (requestCounts[ip] or 0) + 1
    
    if requestCounts[ip] > 100 then
        return false
    end
    return true
end)

-- CORS middleware
use_middleware("https", "/api/*", function(req)
    return {
        status = 200,
        headers = {
            ["Access-Control-Allow-Origin"] = "*",
            ["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE",
            ["Access-Control-Allow-Headers"] = "Content-Type, Authorization"
        }
    }
end)

-- REST API endpoints
handle_http("https", "/api/users", function(req)
    if req.method == "GET" then
        return {
            status = 200,
            headers = {
                ["Content-Type"] = "application/json"
            },
            body = '[{"id": 1, "name": "John"}, {"id": 2, "name": "Jane"}]'
        }
    elseif req.method == "POST" then
        return {
            status = 201,
            headers = {
                ["Content-Type"] = "application/json"
            },
            body = '{"message": "User created", "id": 3}'
        }
    end
end)

-- Dynamic route with multiple parameters
handle_http("https", "/products/:category/:id", function(req)
    local category = req.params[1]
    local id = req.params[2]
    return {
        status = 200,
        headers = {
            ["Content-Type"] = "application/json"
        },
        body = string.format('{"category": "%s", "id": "%s", "name": "Product %s"}', category, id, id)
    }
end)

-- WebSocket chat example
local clients = {}
handle_ws("https", "/chat", function(ws)
    local clientId = #clients + 1
    clients[clientId] = ws
    
    ws:send("Welcome to chat! Your ID: " .. clientId)
    
    while true do
        local message = ws:receive()
        if message then
            for id, client in pairs(clients) do
                if id ~= clientId then
                    client:send(string.format("User %d: %s", clientId, message))
                end
            end
        end
    end
end)

-- File upload handler
handle_http("https", "/upload", function(req)
    if req.method == "POST" then
        return {
            status = 200,
            headers = {
                ["Content-Type"] = "application/json"
            },
            body = '{"message": "File uploaded successfully"}'
        }
    end
end)

-- Start both servers
start_server("http")
start_server("https") 