package main_test

import (
	"archive/zip"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// Q3.2
func TestDummyFile(t *testing.T) {
	buf := make([]byte, 1024)
	rand.Reader.Read(buf)
	fmt.Println(buf)

	// あとはこれをファイルに書き出す
	file, err := os.Create("rand.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(buf)

	// 確かに1024バイトのファイルができる
	// $ wc --bytes rand.txt
	// 1024 rand.txt
}

// Q3.3
func ExampleCreateZipFileWithInterface() {
	source := strings.NewReader("This is sample text to be compressed.")

	file, err := os.Create("3.3.zip")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	// zipファイルの読み書きに構造体そのものではなく、インタフェースを使う
	var writer io.Writer
	writer, err = zipWriter.Create("sample.txt")
	if err != nil {
		panic(err)
	}
	io.Copy(writer, source)

	// Output:
}
