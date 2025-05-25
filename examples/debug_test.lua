-- Example of debug functionality
print("Debug test script started")

-- Function to demonstrate stack trace
function deep_function()
    function deeper_function()
        print("Current stack trace:")
        trace()
    end
    deeper_function()
end

-- Watch this file for changes
watch_file("debug_test.lua", function()
    print("File changed, reloading...")
    reload_script()
end)

-- Demonstrate stack trace
deep_function()

-- Keep script running
while true do
    sleep(1)
end 