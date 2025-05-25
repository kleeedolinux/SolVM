-- Create channels
chan("numbers", 5)
chan("results", 5)

-- Spawn a worker goroutine
go(function()
    while true do
        local num = receive("numbers")
        if num == nil then
            break
        end
        print("Worker received:", num)
        sleep(0.5)  -- Simulate work
        send("results", num * 2)
    end
end)

-- Send some numbers
for i = 1, 5 do
    print("Sending:", i)
    send("numbers", i)
end

-- Receive results
for i = 1, 5 do
    local result = receive("results")
    print("Received result:", result)
end

-- Demonstrate select
chan("ch1")
chan("ch2")

go(function()
    sleep(0.2)
    send("ch1", "message from ch1")
end)

go(function()
    sleep(0.3)
    send("ch2", "message from ch2")
end)

local value, channel = select("ch1", "ch2")
print("Select received:", value, "from", channel)

-- Wait for all goroutines to complete
wait() 