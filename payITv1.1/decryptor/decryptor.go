package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func mapFleExts(rt string, exts []string) ([]string, error) {
	var fleExts []string
	extMap := make(map[string]struct{})
	for _, ext := range exts {
		extMap[ext] = struct{}{}
	}
	er := filepath.Walk(rt, func(pth string, inf os.FileInfo, er error) error {
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

func Decrypt(b64Ctxt, pass string) ([]byte, error) {
	ct, er := base64.StdEncoding.DecodeString(b64Ctxt)
	if er != nil {
		return nil, er
	}
	key := pbkdf2.Key([]byte(pass), []byte(pass), 100000, 32, sha256.New)
	aesBlock, er := aes.NewCipher(key)
	if er != nil {
		return nil, er
	}
	gcm, er := cipher.NewGCM(aesBlock)
	if er != nil {
		return nil, er
	}
	nonceSize := 12
	nonce, ctxt := ct[:nonceSize], ct[nonceSize:]
	pTxt, er := gcm.Open([]byte(nonce), []byte(nonce), []byte(ctxt), nil)
	if er != nil {
		return nil, er
	}
	return pTxt[12:], nil
}

func main() {
	var pass string
	fmt.Print("Enter Password: ")
	fmt.Scanln(&pass)
	exts := []string{".juru"}
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
	flePths, er := mapFleExts(rtDir, exts)
	if er != nil {
		panic("ERR: mapping files")
	}
	for _, flePth := range flePths {
		b64Ct, er := os.ReadFile(flePth)
		if er != nil {
			panic("ERR: reading file data")
		}
		fmt.Printf("Decrypting  ==>  %v\n", filepath.Base(flePth))
		fleDat, er := Decrypt(string(b64Ct), pass)
		if er != nil {
			panic("ERR: decrypting files")
		}
		time.Sleep(200 * time.Millisecond)
		os.WriteFile(flePth, []byte(fleDat), 0644)
	}
	fmt.Println("Finalizing")
	time.Sleep(1 * time.Second)
	fmt.Println("All Done.")
}
