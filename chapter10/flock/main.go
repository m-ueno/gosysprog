package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var sharedFile = "./shared.txt"

func do(workerID int, done chan bool) {
	for i := 0; i < 1000; i++ {
		f, err := os.Open(sharedFile)

		// data, err := ioutil.ReadFile(sharedFile)
		if err != nil {
			panic(err)
		}
		data := make([]byte, 16)
		length, err := f.Read(data)
		if err != nil {
			fmt.Println("Error: failed to read")
			continue
		}
		f.Close()

		fmt.Printf("data: %v\n", data)                        // 48, 49, ...
		fmt.Printf("string(data[:len]): %v\n", data[:length]) // 48, 49, ...

		n, err := strconv.Atoi(string(data[:length]))
		if err != nil {
			panic(err)
		}
		fmt.Printf("n: %v\n", n)                     // 1
		fmt.Printf("Itoa(n): %v\n", strconv.Itoa(n)) // 1
		err = ioutil.WriteFile(sharedFile, []byte(strconv.Itoa(n+1)), 0644)
		if err != nil {
			panic(err)
		}
	}
	done <- true
}

func main() {
	os.Remove(sharedFile)
	ioutil.WriteFile(sharedFile, []byte("0"), 0644)
	done := make(chan bool, 2)
	go do(1, done)
	do(2, done)

	<-done
	<-done
}
