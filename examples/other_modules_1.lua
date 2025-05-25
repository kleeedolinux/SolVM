-- Dotenv example
dotenv.load(".env")
local api_key = dotenv.get("API_KEY", "default_key")
print("API Key:", api_key)

-- Datetime example
local now = datetime.now()
print("Current timestamp:", now)
print("Formatted time:", datetime.format(now, "2006-01-02 15:04:05"))

local future = datetime.add(now, "24h")
print("Tomorrow:", datetime.format(future, "2006-01-02 15:04:05"))

local diff = datetime.diff(now, future)
print("Time difference:", diff)

-- CSV example
local data = {
    {"Name", "Age", "City"},
    {"John", "30", "New York"},
    {"Alice", "25", "London"}
}

csv.write("data.csv", data)
local loaded = csv.read("data.csv")
print("CSV Data:", loaded)

-- File Transfer example
ft.download("https://example.com/file.txt", "downloaded.txt")
ft.copy("downloaded.txt", "backup.txt")
ft.move("backup.txt", "archive/backup.txt")

-- INI example
local config = {
    database = {
        host = "localhost",
        port = "5432",
        user = "admin",
        password = "secret"
    },
    server = {
        port = "8080",
        timeout = "30s"
    }
}

ini.write("config.ini", config)
local loaded_config = ini.read("config.ini")
print("Database host:", loaded_config.database.host)
print("Server port:", loaded_config.server.port) 