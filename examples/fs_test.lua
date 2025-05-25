-- Read a file
local config = read_file("config.json")
local data = json_decode(config)

-- Write to a file
write_file("output.txt", "Hello, World!")

-- List directory contents
local files = list_dir(".")
for _, file in ipairs(files) do
    if not file.is_dir then
        print("File:", file.name, "Size:", file.size)
    end
end