package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"log"
)

// GenerateSymmetricKey generates a random AES key.
func GenerateSymmetricKey(keySize int) ([]byte, error) {
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// EncryptWithAES encrypts the data using AES-GCM.
func EncryptWithAES(data, key []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, data, nil)
	return ciphertext, nonce, nil
}

// DecryptWithAES decrypts the data using AES-GCM.
func DecryptWithAES(ciphertext, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

// EncryptWithRSA encrypts the symmetric key using RSA.
func EncryptWithRSA(key []byte, publicKey *rsa.PublicKey) ([]byte, error) {
	return rsa.EncryptOAEP(nil, rand.Reader, publicKey, key, nil)
}

// DecryptWithRSA decrypts the symmetric key using RSA.
func DecryptWithRSA(encryptedKey []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	return rsa.DecryptOAEP(nil, rand.Reader, privateKey, encryptedKey, nil)
}

func main() {
	// Generate RSA keys
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate RSA key pair: %v", err)
	}
	publicKey := &privateKey.PublicKey

	// Original data
	data := []byte("This is some data that is too large for RSA encryption directly.")

	// Step 1: Generate a symmetric key
	symmetricKey, err := GenerateSymmetricKey(32) // 32 bytes = 256-bit AES key
	if err != nil {
		log.Fatalf("Failed to generate symmetric key: %v", err)
	}

	// Step 2: Encrypt data using the symmetric key
	ciphertext, nonce, err := EncryptWithAES(data, symmetricKey)
	if err != nil {
		log.Fatalf("Failed to encrypt data with AES: %v", err)
	}

	// Step 3: Encrypt the symmetric key using RSA
	encryptedSymmetricKey, err := EncryptWithRSA(symmetricKey, publicKey)
	if err != nil {
		log.Fatalf("Failed to encrypt symmetric key with RSA: %v", err)
	}

	// Encode outputs for storage/transmission
	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	encodedEncryptedKey := base64.StdEncoding.EncodeToString(encryptedSymmetricKey)
	encodedNonce := base64.StdEncoding.EncodeToString(nonce)

	fmt.Println("Encrypted Symmetric Key (RSA):", encodedEncryptedKey)
	fmt.Println("Ciphertext (AES):", encodedCiphertext)
	fmt.Println("Nonce:", encodedNonce)

	// Decryption process:
	// Step 4: Decode the base64 strings
	decodedEncryptedKey, _ := base64.StdEncoding.DecodeString(encodedEncryptedKey)
	decodedCiphertext, _ := base64.StdEncoding.DecodeString(encodedCiphertext)
	decodedNonce, _ := base64.StdEncoding.DecodeString(encodedNonce)

	// Step 5: Decrypt the symmetric key using RSA
	decryptedSymmetricKey, err := DecryptWithRSA(decodedEncryptedKey, privateKey)
	if err != nil {
		log.Fatalf("Failed to decrypt symmetric key with RSA: %v", err)
	}

	// Step 6: Decrypt the data using AES
	decryptedData, err := DecryptWithAES(decodedCiphertext, decryptedSymmetricKey, decodedNonce)
	if err != nil {
		log.Fatalf("Failed to decrypt data with AES: %v", err)
	}

	fmt.Println("Decrypted Data:", string(decryptedData))
}
