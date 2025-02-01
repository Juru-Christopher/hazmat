package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

// GenerateRSAKeys generates an RSA key pair (private and public keys).
func GenerateRSAKeys() (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %v", err)
	}
	return privateKey, nil
}

// EncryptWithRSA encrypts a symmetric AES key using RSA.
func EncryptWithRSA(publicKey *rsa.PublicKey, symmetricKey []byte) (string, error) {
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, symmetricKey, nil)
	if err != nil {
		return "", fmt.Errorf("RSA encryption failed: %v", err)
	}
	return base64.StdEncoding.EncodeToString(encryptedKey), nil
}

// DecryptWithRSA decrypts an encrypted symmetric AES key using RSA.
func DecryptWithRSA(privateKey *rsa.PrivateKey, encryptedKey string) ([]byte, error) {
	encryptedKeyBytes, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode RSA encrypted key: %v", err)
	}
	symmetricKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, encryptedKeyBytes, nil)
	if err != nil {
		return nil, fmt.Errorf("RSA decryption failed: %v", err)
	}
	return symmetricKey, nil
}

// EncryptWithAES encrypts a plaintext message using a symmetric key (AES).
func EncryptWithAES(key []byte, plaintext string) (string, string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Generate a random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", "", fmt.Errorf("failed to generate IV: %v", err)
	}

	// Pad plaintext to be a multiple of the block size
	paddedPlaintext := pad([]byte(plaintext), aes.BlockSize)

	// Encrypt the plaintext
	ciphertext := make([]byte, len(paddedPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	// Encode ciphertext and IV as base64
	return base64.StdEncoding.EncodeToString(ciphertext), base64.StdEncoding.EncodeToString(iv), nil
}

// DecryptWithAES decrypts a ciphertext message using a symmetric key (AES).
func DecryptWithAES(key []byte, ciphertext, iv string) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %v", err)
	}
	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return "", fmt.Errorf("failed to decode IV: %v", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Decrypt the ciphertext
	decrypted := make([]byte, len(ciphertextBytes))
	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(decrypted, ciphertextBytes)

	// Remove padding
	plaintext, err := unpad(decrypted)
	if err != nil {
		return "", fmt.Errorf("failed to unpad plaintext: %v", err)
	}

	return string(plaintext), nil
}

// pad adds PKCS#7 padding to plaintext.
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

// unpad removes PKCS#7 padding from plaintext.
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
	// Generate RSA keys
	privateKey, err := GenerateRSAKeys()
	if err != nil {
		fmt.Println("Error generating RSA keys:", err)
		return
	}
	publicKey := &privateKey.PublicKey

	// Plaintext message
	plaintext := "This is a very long message that exceeds the maximum size for RSA encryption. It will be encrypted using a hybrid RSA-AES scheme."

	// Generate a random AES key
	symmetricKey := make([]byte, 32) // AES-256
	if _, err := io.ReadFull(rand.Reader, symmetricKey); err != nil {
		fmt.Println("Error generating AES key:", err)
		return
	}

	// Encrypt the AES key with RSA
	encryptedSymmetricKey, err := EncryptWithRSA(publicKey, symmetricKey)
	if err != nil {
		fmt.Println("Error encrypting AES key with RSA:", err)
		return
	}

	// Encrypt the message with AES
	ciphertext, iv, err := EncryptWithAES(symmetricKey, plaintext)
	if err != nil {
		fmt.Println("Error encrypting message with AES:", err)
		return
	}

	fmt.Println("Encrypted AES key:", encryptedSymmetricKey)
	fmt.Println("Encrypted message:", ciphertext)
	fmt.Println("IV:", iv)

	// Decrypt the AES key with RSA
	decryptedSymmetricKey, err := DecryptWithRSA(privateKey, encryptedSymmetricKey)
	if err != nil {
		fmt.Println("Error decrypting AES key with RSA:", err)
		return
	}

	// Decrypt the message with AES
	decryptedMessage, err := DecryptWithAES(decryptedSymmetricKey, ciphertext, iv)
	if err != nil {
		fmt.Println("Error decrypting message with AES:", err)
		return
	}

	fmt.Println("Decrypted message:", decryptedMessage)
}
