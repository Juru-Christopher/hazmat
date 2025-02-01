package main

import (
	"fmt"
	"os"

	"github.com/MarkKremer/microphone"
)

func main() {
	mic, _, err := microphone.OpenDefaultStream(44100, 1)
	if err != nil {
		fmt.Println("Error opening microphone:", err)
		os.Exit(1)
	}
	defer mic.Close()

	samples := make([][2]float64, 1024)
	for {
		n, err := mic.Stream(samples)
		if !err {
			fmt.Println("Error reading from microphone:", err)
			break
		}
		// Process samples here
		fmt.Println("Read", n, "samples")
	}
}
