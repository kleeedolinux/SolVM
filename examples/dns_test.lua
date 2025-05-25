-- Example of DNS resolution
print("DNS test script started")

-- List of domains to resolve
local domains = {
    "google.com",
    "github.com",
    "example.com",
    "localhost"
}

-- Resolve each domain
for _, domain in ipairs(domains) do
    local ip = resolve_dns(domain)
    print(string.format("Domain: %-15s IP: %s", domain, ip))
end

-- Keep script running
while true do
    sleep(1)
end 