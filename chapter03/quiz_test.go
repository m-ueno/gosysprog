package main_test

import (
	"archive/zip"
	"bytes"
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

// Q3.5
func myCopyN(w io.Writer, r io.Reader, length int) (written int64, err error) {
	rlimit := io.LimitReader(r, int64(length))
	written, err = io.Copy(w, rlimit)
	return
}

func ExampleCopyN() {
	r := bytes.NewReader([]byte("日本語&フォント"))
	w := bytes.NewBuffer([]byte{})

	size, err := myCopyN(w, r, 10)
	// size, err := io.CopyN(w, r, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(size)
	io.Copy(os.Stdout, w)
	// Output:
	// 10
	// 日本語&
}

// Q3.6
func ExamplePuzzle() {
	computer := strings.NewReader("COMPUTER")
	system := strings.NewReader("SYSTEM")
	programming := strings.NewReader("PROGRAMMING")
	var stream io.Reader
	// ここから
	a := io.NewSectionReader(programming, 5, 1)
	s := io.LimitReader(system, 1)
	c := io.LimitReader(computer, 1)
	i1 := io.NewSectionReader(programming, 8, 1)
	i2 := io.NewSectionReader(programming, 8, 1)

	stream = io.MultiReader(a, s, c, i1, i2)
	// ここまで
	io.Copy(os.Stdout, stream)
	// Output: ASCII
}
