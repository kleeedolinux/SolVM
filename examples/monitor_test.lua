-- Register error handler
on_error(function(err)
    print("Error occurred:", err)
end)

-- Create some channels and goroutines
chan("test", 5)

-- Spawn a goroutine that might cause a memory leak
go(function()
    local table = {}
    for i = 1, 1000 do
        table[i] = "leak test " .. i
    end
    print("Created large table in goroutine")
    sleep(1)
end)

-- Check memory usage
local initial_mem = check_memory()
print("Initial memory stats:")
print("  Allocation diff:", initial_mem.alloc_diff)
print("  Total allocation diff:", initial_mem.total_alloc_diff)
print("  System memory diff:", initial_mem.sys_diff)
print("  Number of goroutines:", initial_mem.goroutines)

-- Create another goroutine
go(function()
    print("Second goroutine running")
    sleep(0.5)
end)

-- Check memory again after some time
sleep(1)
local final_mem = check_memory()
print("\nFinal memory stats:")
print("  Allocation diff:", final_mem.alloc_diff)
print("  Total allocation diff:", final_mem.total_alloc_diff)
print("  System memory diff:", final_mem.sys_diff)
print("  Number of goroutines:", final_mem.goroutines)

-- Get list of goroutines
local goroutines = get_goroutines()
print("\nActive goroutines:")
for id, name in pairs(goroutines) do
    print("  Goroutine", id, ":", name)
end

-- Wait for all goroutines to complete
wait() 