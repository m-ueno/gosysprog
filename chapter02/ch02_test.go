package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

// 2.4.7
func TestFprintf(t *testing.T) {
	fmt.Fprintf(os.Stdout, "Write with os.Stdout at %v", time.Now())
}

func TestJSONEncode(t *testing.T) {
	/*
	 * NewEncoder()は与えられたWriterに書き込むencoderを返す
	 **/

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(map[string]string{
		"example": "json",
		"hello":   "world",
	})
}

func TestQ1(t *testing.T) {
	file, err := os.Create("q1.txt")
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(file, "Write with os.Open at %v", time.Now())
}

func TestQ2(t *testing.T) {
	encoder := csv.NewWriter(os.Stdout)
	encoder.Write([]string{"a", "b"})
	encoder.Write([]string{"1", "2", "1,2,3"})
	encoder.Flush()
}
