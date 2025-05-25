import("math_utils")

local math = math_utils
local result = math.add(5, 3)
print("5 + 3 =", result)

result = math.multiply(4, 6)
print("4 * 6 =", result)

local data = {
    operation = "subtraction",
    a = 10,
    b = 4
}

local json_str = json_encode(data)
print("Encoded:", json_str)

local decoded = json_decode(json_str)
result = math.subtract(decoded.a, decoded.b)
print("10 - 4 =", result) 