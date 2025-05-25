-- Example of UDP communication
print("UDP test script started")

-- Create UDP receiver
local udp = udp_recvfrom(8080)
print("UDP receiver listening on port 8080")

-- Handle incoming messages
go(function()
    while true do
        local msg = udp.receive()
        if msg then
            print("Received from", msg.addr, ":", msg.data)
            -- Send response
            udp_sendto(msg.addr, msg.port, "Echo: " .. msg.data)
        end
    end
end)

-- Send some test messages
go(function()
    sleep(1) -- Wait for receiver to start
    
    -- Send to localhost
    udp_sendto("localhost", 8080, "Hello, UDP!")
    udp_sendto("localhost", 8080, "How are you?")
    
    -- Send to broadcast
    udp_sendto("255.255.255.255", 8080, "Broadcast message!")
end)

-- Keep script running
while true do
    sleep(1)
end 