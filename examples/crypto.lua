-- Hash functions
local data = "Hello, World!"
print("MD5:", crypto.md5(data))
print("SHA1:", crypto.sha1(data))
print("SHA256:", crypto.sha256(data))
print("SHA512:", crypto.sha512(data))

-- Base64 encoding/decoding
local encoded = crypto.base64_encode(data)
print("Base64 encoded:", encoded)
local decoded = crypto.base64_decode(encoded)
print("Base64 decoded:", decoded)

-- AES encryption/decryption
local key = "1234567890123456"  -- 16 bytes for AES-128
local iv = "1234567890123456"   -- 16 bytes for AES-128
local encrypted = crypto.aes_encrypt(data, key, iv)
print("AES encrypted:", encrypted)
local decrypted = crypto.aes_decrypt(encrypted, key, iv)
print("AES decrypted:", decrypted)

-- DES encryption/decryption
local des_key = "12345678"  -- 8 bytes for DES
local des_iv = "12345678"   -- 8 bytes for DES
local des_encrypted = crypto.des_encrypt(data, des_key, des_iv)
print("DES encrypted:", des_encrypted)
local des_decrypted = crypto.des_decrypt(des_encrypted, des_key, des_iv)
print("DES decrypted:", des_decrypted)

-- RC4 encryption/decryption
local rc4_key = "mysecretkey"
local rc4_encrypted = crypto.rc4_encrypt(data, rc4_key)
print("RC4 encrypted:", rc4_encrypted)
local rc4_decrypted = crypto.rc4_decrypt(rc4_encrypted, rc4_key)
print("RC4 decrypted:", rc4_decrypted)

-- RSA key generation
local rsa_keys = crypto.rsa_generate(2048)
print("RSA Private Key:", rsa_keys.private)
print("RSA Public Key:", rsa_keys.public)

-- Random bytes generation
local random = crypto.random_bytes(32)
print("Random bytes (32):", random) 