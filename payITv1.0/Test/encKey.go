package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"log"
)

// EncryptPrivateKey encrypts the private key with a password.
func EncryptPrivateKey(privateKey *rsa.PrivateKey, password string) (string, error) {
	// Marshal the private key into DER format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	// Derive AES key from password (use a proper KDF in production)
	key := deriveKey(password)

	// Encrypt the private key bytes
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, privateKeyBytes, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPrivateKey decrypts the private key using the password.
func DecryptPrivateKey(encryptedPrivateKey string, password string) (*rsa.PrivateKey, error) {
	// Decode the Base64 encrypted private key
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedPrivateKey)
	if err != nil {
		return nil, err
	}

	// Derive AES key from password (use a proper KDF in production)
	key := deriveKey(password)

	// Decrypt the private key bytes
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedBytes) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := encryptedBytes[:nonceSize], encryptedBytes[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	// Parse the private key
	privateKey, err := x509.ParsePKCS1PrivateKey(plaintext)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// Helper: Derive AES key from password (for simplicity, using a fixed size key)
func deriveKey(password string) []byte {
	key := make([]byte, 32)     // AES-256 requires 32 bytes key
	copy(key, []byte(password)) // For simplicity, just copying password bytes (not secure for production)
	return key
}

func main() {
	// Generate RSA keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Error generating RSA keys: %v", err)
	}

	publicKey := &privateKey.PublicKey
	fmt.Println("Public Key:", base64.StdEncoding.EncodeToString(x509.MarshalPKCS1PublicKey(publicKey)))

	// Password for encryption
	password := "strongpassword123"

	// Encrypt the private key
	encryptedPrivateKey, err := EncryptPrivateKey(privateKey, password)
	if err != nil {
		log.Fatalf("Error encrypting private key: %v", err)
	}
	fmt.Println("Encrypted Private Key:", encryptedPrivateKey)

	// Decrypt the private key
	decryptedPrivateKey, err := DecryptPrivateKey(encryptedPrivateKey, password)
	if err != nil {
		log.Fatalf("Error decrypting private key: %v", err)
	}

	// Verify that the decrypted key matches the original
	if decryptedPrivateKey.Equal(privateKey) {
		fmt.Println("Decryption successful, keys match!")
	} else {
		fmt.Println("Decryption failed, keys do not match.")
	}
}
