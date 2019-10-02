package main

import (
	"fmt"
	"github.com/edsrzf/mmap-go"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

func ExampleMMap() {
	var testData = []byte("0123456789ABCDEF")
	var testPath = filepath.Join(os.TempDir(), "testdata")
	err := ioutil.WriteFile(testPath, testData, 0644)
	if err != nil {
		panic(err)
	}

	// OpenFile(): 前章参照
	f, err := os.OpenFile(testPath, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	m, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		panic(err)
	}
	defer m.Unmap()

	m[9] = 'X'
	m.Flush()

	fileData, err := ioutil.ReadAll(f)

	fmt.Printf("original: %s\n", testData)
	fmt.Printf("mmap:     %s\n", m)
	fmt.Printf("file:     %s\n", fileData)

	// Output:
	// original: 0123456789ABCDEF
	// mmap:     012345678XABCDEF
	// file:     012345678XABCDEF
}

func ExampleMMapCopy() {
	var testData = []byte("0123456789ABCDEF")
	var testPath = filepath.Join(os.TempDir(), "testdata")
	err := ioutil.WriteFile(testPath, testData, 0644)
	if err != nil {
		panic(err)
	}

	// OpenFile(): 前章参照
	f, err := os.OpenFile(testPath, os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	m, err := mmap.Map(f, mmap.COPY, 0)
	if err != nil {
		panic(err)
	}
	defer m.Unmap()

	m[9] = 'X'
	m.Flush()

	fileData, err := ioutil.ReadAll(f)

	fmt.Printf("original: %s\n", testData)
	fmt.Printf("mmap:     %s\n", m)
	fmt.Printf("file:     %s\n", fileData)

	// Output:
	// original: 0123456789ABCDEF
	// mmap:     012345678XABCDEF
	// file:     0123456789ABCDEF
}

func ExampleMMapExec() {
	binfile := "./hello"
	f, err := os.Open(binfile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	m, err := mmap.Map(f, mmap.EXEC, 0)
	if err != nil {
		panic(err)
	}
	defer m.Unmap()

	m[10000]= 'X'

	fdPath := fmt.Sprintf("/proc/self/fd/%d", f.Fd())
	syscall.Exec(fdPath, []string{}, nil)

	// Output: hello!
}