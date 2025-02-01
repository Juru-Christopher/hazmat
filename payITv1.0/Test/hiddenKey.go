package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func GenPrKey(keySize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, fmt.Errorf("ERR: failed to generate RSA private key: %v", err)
	}
	return privateKey, nil
}

func EncryptData(pubKey *rsa.PublicKey, data string) (string, error) {
	cTxt, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, []byte(data), nil)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %v", err)
	}
	return base64.StdEncoding.EncodeToString(cTxt), nil
}

func DecryptData(prKey *rsa.PrivateKey, cTxt string) (string, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(cTxt)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %v", err)
	}
	data, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, prKey, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %v", err)
	}

	return string(data), nil
}

func main() {
	prKey, err := GenPrKey(2048)
	if err != nil {
		fmt.Println("Error generating RSA keys:", err)
		return
	}
	pubKey := &prKey.PublicKey
	data := "This is a secret message."
	cTxt, err := EncryptData(pubKey, data)
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return
	}
	fmt.Println("Encrypted ciphertext:", cTxt)
	decryptedText, err := DecryptData(prKey, cTxt)
	if err != nil {
		fmt.Println("Error decrypting data:", err)
		return
	}
	fmt.Println("Decrypted plaintext:", decryptedText)
}
