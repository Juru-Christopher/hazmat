package main

import (
	"encoding/base64"
	"fmt"
)

// ScrambleBytes scrambles a byte slice using a simple XOR operation with a key.
func ScrambleBytes(data []byte, key byte) []byte {
	scrambled := make([]byte, len(data))
	for i, b := range data {
		scrambled[i] = b ^ key
	}
	return scrambled
}

// UnscrambleBytes reverses the scrambling process using the same key.
func UnscrambleBytes(data []byte, key byte) []byte {
	// XOR again to reverse
	return ScrambleBytes(data, key)
}

func main() {
	// Base64 string to scramble
	originalString := "SGVsbG8sIFdvcmxkIQ==" // "Hello, World!" in Base64
	fmt.Println("Original Base64 string:", originalString)

	// Decode the Base64 string to bytes
	decodedBytes, err := base64.StdEncoding.DecodeString(originalString)
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		return
	}

	// Scramble the bytes
	key := byte(42) // Choose a key for scrambling (must be the same for unscrambling)
	scrambledBytes := ScrambleBytes(decodedBytes, key)
	scrambledString := base64.StdEncoding.EncodeToString(scrambledBytes)
	fmt.Println("Scrambled Base64 string:", scrambledString)

	// Reverse the scrambling
	scrambledBytesDecoded, err := base64.StdEncoding.DecodeString(scrambledString)
	if err != nil {
		fmt.Println("Error decoding scrambled Base64:", err)
		return
	}
	unscrambledBytes := UnscrambleBytes(scrambledBytesDecoded, key)
	unscrambledString := base64.StdEncoding.EncodeToString(unscrambledBytes)
	fmt.Println("Unscrambled Base64 string:", unscrambledString)
}
