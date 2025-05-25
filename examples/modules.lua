-- UUID Module Example
local uuid1 = uuid.v4()
local uuid2 = uuid.v4_without_hyphens()
print("UUID v4:", uuid1)
print("UUID v4 without hyphens:", uuid2)
print("Is valid UUID:", uuid.is_valid(uuid1))

-- Random Module Example
local random_num = random.number()
local random_int = random.int(1, 100)
local random_str = random.string(10)
print("Random number:", random_num)
print("Random integer between 1 and 100:", random_int)
print("Random string:", random_str)

-- TOML Module Example
local toml_data = {
    title = "TOML Example",
    owner = {
        name = "John Doe",
        age = 30
    },
    database = {
        enabled = true,
        ports = {8000, 8001, 8002}
    }
}
local toml_str = toml.encode(toml_data)
print("TOML encoded:", toml_str)
local decoded_toml = toml.decode(toml_str)
print("TOML decoded:", decoded_toml.title)

-- YAML Module Example
local yaml_data = {
    name = "YAML Example",
    version = "1.0",
    dependencies = {
        "package1",
        "package2"
    },
    config = {
        debug = true,
        timeout = 30
    }
}
local yaml_str = yaml.encode(yaml_data)
print("YAML encoded:", yaml_str)
local decoded_yaml = yaml.decode(yaml_str)
print("YAML decoded:", decoded_yaml.name)

-- JSONC Module Example
local jsonc_str = [[
{
    "name": "JSONC Example",
    "version": "1.0",
    "config": {
        "debug": true,
        "timeout": 30
    }
}
]]
local decoded_jsonc = jsonc.decode(jsonc_str)
print("JSONC decoded:", decoded_jsonc.name)
local encoded_jsonc = jsonc.encode(decoded_jsonc)
print("JSONC encoded:", encoded_jsonc)

-- Text Module Example
local text_str = "  Hello, World!  "
print("Original:", text_str)
print("Trimmed:", text.trim(text_str))
print("Lowercase:", text.lower(text_str))
print("Uppercase:", text.upper(text_str))
print("Title case:", text.title(text_str))

local words = text.split(text.trim(text_str), " ")
print("Split words:", text.join(words, ", "))

local replaced = text.replace(text_str, "World", "Lua")
print("Replaced:", replaced)

print("Contains 'Hello':", text.contains(text_str, "Hello"))
print("Starts with '  H':", text.starts_with(text_str, "  H"))
print("Ends with '!  ':", text.ends_with(text_str, "!  "))

local padded = text.pad_left("123", 5, "0")
print("Padded left:", padded)

local repeated = text.repeat_str("abc", 3)
print("Repeated:", repeated) 