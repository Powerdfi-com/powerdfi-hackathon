package helpers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

const PaddingLength = 32

func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	return append(data, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func EncryptPlainText(text []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// Pad the plaintext using PKCS#7
	padded := PKCS7Padding(text, aes.BlockSize)

	ciphertext := make([]byte, aes.BlockSize+len(padded))
	cbc := cipher.NewCBCEncrypter(block, ciphertext[:aes.BlockSize])
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], padded)

	return ciphertext, nil
}

func DecryptPlainText(text []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(text)-aes.BlockSize)
	cbc := cipher.NewCBCDecrypter(block, text[:aes.BlockSize])
	cbc.CryptBlocks(plaintext, text[aes.BlockSize:])

	// Remove padding using PKCS#7
	return PKCS7Unpadding(plaintext)
}

func PKCS7Unpadding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid data length")
	}
	padding := data[len(data)-1]
	if int(padding) < 1 || int(padding) > aes.BlockSize {
		return nil, errors.New("invalid padding size")
	}
	for i := len(data) - 1; i > len(data)-int(padding); i-- {
		if data[i] != padding {
			return nil, errors.New("invalid padding format")
		}
	}
	return data[:len(data)-int(padding)], nil
}
