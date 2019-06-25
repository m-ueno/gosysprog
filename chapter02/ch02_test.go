package main_test

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// 2.4.7
func ExampleFprintf() {
	fmt.Fprintf(os.Stdout, "Write with os.Stdout at %v", time.Now())
}

func ExampleJSONEncode() {
	/*
	 * NewEncoder()は与えられたWriterに書き込むencoderを返す
	 **/

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(map[string]string{
		"example": "json",
		"hello":   "world",
	})
	// Output:
	// {
	//     "example": "json",
	//     "hello": "world"
	// }
}

func ExampleQ1() {
	file, err := os.Create("q1.txt")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(file, "Write with os.Open at %v", time.Now())
}

func ExampleQ2() {
	encoder := csv.NewWriter(os.Stdout)
	encoder.Write([]string{"a", "b"})
	encoder.Write([]string{"1", "2", "1,2,3"})
	encoder.Flush()
	// Output:
	// a,b
	// 1,2,"1,2,3"
}
