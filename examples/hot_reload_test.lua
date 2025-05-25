-- Example of hot-reload functionality
print("Hot-reload test script started")

-- Configuration that can be modified
local config = {
    message = "Hello, World!",
    count = 0
}

-- Function to print current state
function print_state()
    print("Current state:")
    print("  Message:", config.message)
    print("  Count:", config.count)
    config.count = config.count + 1
end

-- Watch this file for changes
watch_file("hot_reload_test.lua", function()
    print("\nFile changed, reloading...")
    reload_script()
end)

-- Print state every 5 seconds
set_interval(function()
    print_state()
end, 5)

-- Keep script running
while true do
    sleep(1)
end 