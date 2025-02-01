package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

func rdK(repoURL string) (string, error) {
	resp, er := http.Get(repoURL)
	if er != nil {
		return "", er
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ERR: bad status: %s", resp.Status)
	}
	dat, er := io.ReadAll(resp.Body)
	if er != nil {
		return "", er
	}
	return string(dat), nil
}

func main() {
	homeDir := ""
	switch runtime.GOOS {
	case "windows":
		homeDir = os.Getenv("USERPROFILE")
	default:
		homeDir = os.Getenv("HOME") + "/Desktop/malware"
	}
	if homeDir == "" {
		homeDir = "/home"
	}
	fmt.Println(homeDir)
}
