local data = {
    name = "test",
    values = {1, 2, 3}
}

local json_str = json_encode(data)
print("Encoded:", json_str)

local decoded = json_decode(json_str)
print("Decoded name:", decoded.name)

print("Sleeping for 1 second...")
sleep(1)
print("Awake!") 