package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func Encrypt(data []byte) (string, error) {
	pass := os.Getenv("ENCRYPTION_KEY")
	if pass == "" {
		return "", fmt.Errorf("encryption key not set")
	}

	key := pbkdf2.Key([]byte(pass), []byte(pass), 100000, 32, sha256.New)
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}
	ct := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ct), nil
}

func main() {
	os.Setenv("ENCRYPTION_KEY", "your-secure-key") // Set this securely in production
	encrypted, err := Encrypt([]byte("Sensitive Data"))
	if err != nil {
		fmt.Println("Encryption failed:", err)
		return
	}
	fmt.Println("Encrypted data:", encrypted)
}
