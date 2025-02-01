package main

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// Encrypt encrypts plaintext using AES with a password, without IV or salt.
func Encrypt(password, plaintext string) (string, error) {
	// Derive a 32-byte key from the password using SHA-256
	key := sha256.Sum256([]byte(password))

	// Create AES cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Pad the plaintext to a multiple of the block size
	paddedPlaintext := pad([]byte(plaintext), aes.BlockSize)

	// Encrypt the plaintext in ECB mode (insecure for most use cases)
	ciphertext := make([]byte, len(paddedPlaintext))
	for i := 0; i < len(paddedPlaintext); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], paddedPlaintext[i:i+aes.BlockSize])
	}

	// Return ciphertext encoded as base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES with a password, without IV or salt.
func Decrypt(password, ciphertextB64 string) (string, error) {
	// Derive a 32-byte key from the password using SHA-256
	key := sha256.Sum256([]byte(password))

	// Decode the base64 ciphertext
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %v", err)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Decrypt the ciphertext in ECB mode (insecure for most use cases)
	decrypted := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += aes.BlockSize {
		block.Decrypt(decrypted[i:i+aes.BlockSize], ciphertext[i:i+aes.BlockSize])
	}

	// Remove padding
	plaintext, err := unpad(decrypted)
	if err != nil {
		return "", fmt.Errorf("failed to unpad plaintext: %v", err)
	}

	return string(plaintext), nil
}

// pad adds padding to the plaintext to make it a multiple of the block size.
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

// unpad removes padding from the decrypted plaintext.
func unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:length-padding], nil
}

func main() {
	password := "simple_password"
	plaintext := "This is a secret message."

	// Encrypt the plaintext
	ciphertext, err := Encrypt(password, plaintext)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}
	fmt.Println("Ciphertext:", ciphertext)

	// Decrypt the ciphertext
	decrypted, err := Decrypt(password, ciphertext)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}
	fmt.Println("Decrypted:", decrypted)
}
