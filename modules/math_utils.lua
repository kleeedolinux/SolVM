local math_utils = {}

metadata({
    name = "math_utils",
    version = "1.0.0",
    author = "SolVM Team",
    description = "Basic math operations module",
    repository = "https://github.com/yourusername/solvm",
    license = "MIT"
})

function math_utils.add(a, b)
    return a + b
end

function math_utils.subtract(a, b)
    return a - b
end

function math_utils.multiply(a, b)
    return a * b
end

function math_utils.divide(a, b)
    if b == 0 then
        error("Division by zero")
    end
    return a / b
end

return math_utils 