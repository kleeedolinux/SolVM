package modules

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rand"
	"crypto/rc4"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"

	lua "github.com/yuin/gopher-lua"
)

func RegisterCryptoModule(L *lua.LState) {
	cryptoModule := L.NewTable()
	L.SetGlobal("crypto", cryptoModule)

	
	L.SetField(cryptoModule, "md5", L.NewFunction(func(L *lua.LState) int {
		data := L.CheckString(1)
		hash := md5.Sum([]byte(data))
		L.Push(lua.LString(hex.EncodeToString(hash[:])))
		return 1
	}))

	L.SetField(cryptoModule, "sha1", L.NewFunction(func(L *lua.LState) int {
		data := L.CheckString(1)
		hash := sha1.Sum([]byte(data))
		L.Push(lua.LString(hex.EncodeToString(hash[:])))
		return 1
	}))

	L.SetField(cryptoModule, "sha256", L.NewFunction(func(L *lua.LState) int {
		data := L.CheckString(1)
		hash := sha256.Sum256([]byte(data))
		L.Push(lua.LString(hex.EncodeToString(hash[:])))
		return 1
	}))

	L.SetField(cryptoModule, "sha512", L.NewFunction(func(L *lua.LState) int {
		data := L.CheckString(1)
		hash := sha512.Sum512([]byte(data))
		L.Push(lua.LString(hex.EncodeToString(hash[:])))
		return 1
	}))

	
	L.SetField(cryptoModule, "base64_encode", L.NewFunction(func(L *lua.LState) int {
		data := L.CheckString(1)
		encoded := base64.StdEncoding.EncodeToString([]byte(data))
		L.Push(lua.LString(encoded))
		return 1
	}))

	L.SetField(cryptoModule, "base64_decode", L.NewFunction(func(L *lua.LState) int {
		encoded := L.CheckString(1)
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			L.RaiseError("failed to decode base64: " + err.Error())
			return 0
		}
		L.Push(lua.LString(string(decoded)))
		return 1
	}))

	
	L.SetField(cryptoModule, "aes_encrypt", L.NewFunction(func(L *lua.LState) int {
		plaintext := L.CheckString(1)
		key := L.CheckString(2)
		iv := L.CheckString(3)

		block, err := aes.NewCipher([]byte(key))
		if err != nil {
			L.RaiseError("failed to create cipher: " + err.Error())
			return 0
		}

		mode := cipher.NewCBCEncrypter(block, []byte(iv))
		padded := pkcs7Padding([]byte(plaintext), aes.BlockSize)
		ciphertext := make([]byte, len(padded))
		mode.CryptBlocks(ciphertext, padded)

		L.Push(lua.LString(base64.StdEncoding.EncodeToString(ciphertext)))
		return 1
	}))

	L.SetField(cryptoModule, "aes_decrypt", L.NewFunction(func(L *lua.LState) int {
		ciphertext := L.CheckString(1)
		key := L.CheckString(2)
		iv := L.CheckString(3)

		block, err := aes.NewCipher([]byte(key))
		if err != nil {
			L.RaiseError("failed to create cipher: " + err.Error())
			return 0
		}

		decoded, err := base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			L.RaiseError("failed to decode base64: " + err.Error())
			return 0
		}

		mode := cipher.NewCBCDecrypter(block, []byte(iv))
		plaintext := make([]byte, len(decoded))
		mode.CryptBlocks(plaintext, decoded)

		unpadded, err := pkcs7Unpadding(plaintext)
		if err != nil {
			L.RaiseError("failed to unpad: " + err.Error())
			return 0
		}

		L.Push(lua.LString(string(unpadded)))
		return 1
	}))

	
	L.SetField(cryptoModule, "des_encrypt", L.NewFunction(func(L *lua.LState) int {
		plaintext := L.CheckString(1)
		key := L.CheckString(2)
		iv := L.CheckString(3)

		block, err := des.NewCipher([]byte(key))
		if err != nil {
			L.RaiseError("failed to create cipher: " + err.Error())
			return 0
		}

		mode := cipher.NewCBCEncrypter(block, []byte(iv))
		padded := pkcs7Padding([]byte(plaintext), des.BlockSize)
		ciphertext := make([]byte, len(padded))
		mode.CryptBlocks(ciphertext, padded)

		L.Push(lua.LString(base64.StdEncoding.EncodeToString(ciphertext)))
		return 1
	}))

	L.SetField(cryptoModule, "des_decrypt", L.NewFunction(func(L *lua.LState) int {
		ciphertext := L.CheckString(1)
		key := L.CheckString(2)
		iv := L.CheckString(3)

		block, err := des.NewCipher([]byte(key))
		if err != nil {
			L.RaiseError("failed to create cipher: " + err.Error())
			return 0
		}

		decoded, err := base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			L.RaiseError("failed to decode base64: " + err.Error())
			return 0
		}

		mode := cipher.NewCBCDecrypter(block, []byte(iv))
		plaintext := make([]byte, len(decoded))
		mode.CryptBlocks(plaintext, decoded)

		unpadded, err := pkcs7Unpadding(plaintext)
		if err != nil {
			L.RaiseError("failed to unpad: " + err.Error())
			return 0
		}

		L.Push(lua.LString(string(unpadded)))
		return 1
	}))

	
	L.SetField(cryptoModule, "rc4_encrypt", L.NewFunction(func(L *lua.LState) int {
		plaintext := L.CheckString(1)
		key := L.CheckString(2)

		cipher, err := rc4.NewCipher([]byte(key))
		if err != nil {
			L.RaiseError("failed to create cipher: " + err.Error())
			return 0
		}

		ciphertext := make([]byte, len(plaintext))
		cipher.XORKeyStream(ciphertext, []byte(plaintext))

		L.Push(lua.LString(base64.StdEncoding.EncodeToString(ciphertext)))
		return 1
	}))

	L.SetField(cryptoModule, "rc4_decrypt", L.NewFunction(func(L *lua.LState) int {
		ciphertext := L.CheckString(1)
		key := L.CheckString(2)

		decoded, err := base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			L.RaiseError("failed to decode base64: " + err.Error())
			return 0
		}

		cipher, err := rc4.NewCipher([]byte(key))
		if err != nil {
			L.RaiseError("failed to create cipher: " + err.Error())
			return 0
		}

		plaintext := make([]byte, len(decoded))
		cipher.XORKeyStream(plaintext, decoded)

		L.Push(lua.LString(string(plaintext)))
		return 1
	}))

	
	L.SetField(cryptoModule, "rsa_generate", L.NewFunction(func(L *lua.LState) int {
		bits := L.OptInt(1, 2048)
		privateKey, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			L.RaiseError("failed to generate RSA key: " + err.Error())
			return 0
		}

		privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		})

		publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
		publicKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		})

		result := L.NewTable()
		L.SetField(result, "private", lua.LString(string(privateKeyPEM)))
		L.SetField(result, "public", lua.LString(string(publicKeyPEM)))
		L.Push(result)
		return 1
	}))

	
	L.SetField(cryptoModule, "random_bytes", L.NewFunction(func(L *lua.LState) int {
		length := L.CheckInt(1)
		if length <= 0 {
			L.RaiseError("length must be positive")
			return 0
		}

		bytes := make([]byte, length)
		_, err := io.ReadFull(rand.Reader, bytes)
		if err != nil {
			L.RaiseError("failed to generate random bytes: " + err.Error())
			return 0
		}

		L.Push(lua.LString(hex.EncodeToString(bytes)))
		return 1
	}))
}


func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, len(data)+padding)
	copy(padtext, data)
	for i := len(data); i < len(padtext); i++ {
		padtext[i] = byte(padding)
	}
	return padtext
}


func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, errors.New("invalid padding")
	}
	return data[:length-padding], nil
}
