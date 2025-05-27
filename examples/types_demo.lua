-- Type checking
local values = {
    42,           -- number
    42.5,         -- number
    "hello",      -- string
    true,         -- boolean
    nil,          -- nil
    {},           -- table
    function() end -- function
}

print("Type checking examples:")
for i, v in ipairs(values) do
    print(string.format("Value %d:", i))
    print("  type() =", type(v))
    print("  types.type() =", types.type(v))
    print("  is_number =", types.is_number(v))
    print("  is_integer =", types.is_integer(v))
    print("  is_string =", types.is_string(v))
    print("  is_boolean =", types.is_boolean(v))
    print("  is_nil =", types.is_nil(v))
    print("  is_table =", types.is_table(v))
    print("  is_function =", types.is_function(v))
    print("  is_callable =", types.is_callable(v))
    print()
end

-- Custom callable object
local obj = {
    __call = function(self, ...)
        print("Called with:", ...)
    end
}
setmetatable(obj, obj)

print("Custom callable object:")
print("is_callable =", types.is_callable(obj))
obj("test", 123) 