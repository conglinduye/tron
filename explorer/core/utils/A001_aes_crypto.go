package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// encryptPassword 对aes加密密码进行hash
func encryptPassword(password, salt string) []byte {
	encryptedKey := pbkdf2.Key([]byte(password), []byte(salt), 1, 256/8, sha512.New)
	return encryptedKey
}

// AesCtrEncrypt 使用密码，佐料对数据进行加密
func AesCtrEncrypt(passwd, salt string, data []byte) []byte {
	key := encryptPassword(passwd, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(data))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	// fmt.Printf("%v\n%v\n", iv, hexEncode(iv))

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], data)

	return cipherText
}

// AesCtrDecrypt 使用密码，佐料对数据进行解密
func AesCtrDecrypt(passwd, salt string, data []byte) []byte {
	key := encryptPassword(passwd, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	iv := data[:aes.BlockSize]
	cipherText := data[aes.BlockSize:]
	// fmt.Printf("%v\n%v\n", iv, hexEncode(iv))

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return cipherText
}
