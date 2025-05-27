-- String operations
local str = "hello world"
local parts = utils.split(str, " ")
print("Split:", tablex.pretty(parts))

local joined = utils.join(parts, "-")
print("Joined:", joined)

-- String escaping
local path = "C:\\Program Files\\App\\file.txt"
local escaped = utils.escape(path)
print("Escaped:", escaped)
local unescaped = utils.unescape(escaped)
print("Unescaped:", unescaped)

-- Table operations
local tbl = {1, 2, 3, 4, 5}
local a, b, c = utils.unpack(tbl, 2, 4)
print("Unpacked:", a, b, c)

-- Environment operations
local function test()
    local env = utils.getfenv(1)
    print("Function environment:", tablex.pretty(env))
end

local newEnv = {x = 42, y = 24}
utils.setfenv(test, newEnv)
test() 