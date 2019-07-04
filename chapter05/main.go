package main

import (
	"io"
	"os"
)

func main() {
	file, err := os.Open("main.go") // Set debug point
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(os.Stdout, file)
	// Output:
}
