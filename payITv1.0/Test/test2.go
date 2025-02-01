package main

import (
	"fmt"
	"os"
)

func main() {
	files, er := os.ReadDir("../")
	if er != nil {
		panic("Error")
	}
	for _, file := range files {
		fmt.Println(file)
	}
}
