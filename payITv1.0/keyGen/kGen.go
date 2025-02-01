package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func gtK(prKPth, pbKPth string) error {
	prK, er := rsa.GenerateKey(rand.Reader, 2048)
	if er != nil {
		return fmt.Errorf("ERR: generating PRK %v", er)
	}
	prKFle, er := os.Create(prKPth)
	if er != nil {
		return fmt.Errorf("ERR: creating PRKFLE %v", er)
	}
	defer prKFle.Close()
	prKByts := x509.MarshalPKCS1PrivateKey(prK)
	er = pem.Encode(prKFle, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: prKByts,
	})
	if er != nil {
		return fmt.Errorf("ERR: writing PRKBYTS %v", er)
	}
	pbKFle, er := os.Create(pbKPth)
	if er != nil {
		return fmt.Errorf("ERR: creating PBKFLE: %v", er)
	}
	defer pbKFle.Close()
	pbKByts, er := x509.MarshalPKIXPublicKey(&prK.PublicKey)
	if er != nil {
		return fmt.Errorf("ERR: marshalling PBK: %v", er)
	}
	er = pem.Encode(pbKFle, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pbKByts,
	})
	if er != nil {
		return fmt.Errorf("ERR: writing PBKFLE: %v", er)
	}
	return nil
}

func main() {
	prKPth := "./prK.pem"
	pbKPth := "./pbK.pem"
	er := gtK(prKPth, pbKPth)
	if er != nil {
		panic("ERR: generating key pairs.")
	}
	fmt.Println("Keys saved.")
}
