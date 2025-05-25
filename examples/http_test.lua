-- Test GET request
local response = http_get("https://httpbin.org/get")
print("GET Response:")
print("  Status:", response.status)
print("  Body:", response.body)

-- Test POST request with JSON
local data = {
    name = "test",
    value = 42
}
local json_data = json_encode(data)
local post_response = http_post("https://httpbin.org/post", json_data)
print("\nPOST Response:")
print("  Status:", post_response.status)
print("  Body:", post_response.body)

-- Test custom request with headers
local headers = {
    ["X-Custom-Header"] = "test-value",
    ["Authorization"] = "Bearer test-token"
}
local custom_response = http_request("GET", "https://httpbin.org/headers", "", headers)
print("\nCustom Request Response:")
print("  Status:", custom_response.status)
print("  Body:", custom_response.body)

-- Test error handling
on_error(function(err)
    print("Error occurred:", err)
end)

-- This should trigger an error
local error_response = http_get("https://invalid-url-that-does-not-exist.com")
if error_response == nil then
    print("Expected error for invalid URL")
end 