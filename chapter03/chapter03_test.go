package main_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

// 3.1
func ExampleUseReaderDirectly() {
	r := bytes.NewBuffer([]byte("asdfasdf")) // p.45 bytes.Bufferだけ覚えとけ
	buffer := make([]byte, 1024)
	size, err := r.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println(size)
	// Output:
	// 8
}

// 3.2.1
func ExampleHelper_ReadAll() {
	reader := bytes.NewBufferString("abcd") // p.45 bytes.Bufferだけ覚えとけ
	buffer, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}

	fmt.Println(buffer)
	fmt.Println(string(buffer))

	// Output:
	// [97 98 99 100]
	// abcd
}

// 3.2.2
func ExampleHelper_Copy() {
	reader := bytes.NewBufferString("abcde")
	// reader := strings.NewReader("abcde") // p.45 bytes.Bufferだけ覚えとけ
	writer := bytes.NewBufferString("")

	size, err := io.Copy(writer, reader)
	if err != nil {
		panic(err)
	}
	fmt.Println(size)
	fmt.Println(reader)
	fmt.Println(writer)

	// Output:
	// 5
	//
	// abcde
}

// 3.2.2
func ExampleHelper_CopyBuffer() {
	// 8KB buffer
	buffer := make([]byte, 8*1024)
	src := bytes.NewBuffer([]byte{})
	dst := bytes.NewBufferString("")

	for i := 0; i < 1024; i++ {
		src.Write([]byte("1234567\n"))
	}

	written, err := io.CopyBuffer(dst, src, buffer)
	if err != nil {
		panic(err)
	}
	fmt.Println(written)
	// Output:
	// 8192
}

// 3.3.2
func ExampleInterface_Cast() {
	// io.ReadCloserインタフェースが要求されているとき、ダミーのClose()を生やす
	reader := strings.NewReader("テストデータ")
	var _ io.ReadCloser = ioutil.NopCloser(reader)
	// Output:
}

// 3.3.2
func Example_NewReadWriter() {
	var reader *bufio.Reader = bufio.NewReader(strings.NewReader("テストデータ"))
	var writer *bufio.Writer = bufio.NewWriter(os.Stdout)
	readWriter := bufio.NewReadWriter(reader, writer)
	var _ io.ReadWriter = readWriter

	readWriter.Write([]byte("追加データ"))
	readWriter.Flush()

	_, err := io.Copy(os.Stdout, readWriter)
	if err != nil {
		panic(err)
	}

	// Output:
	// 追加データテストデータ
}

// 3.4.2
func Example_FileRead() {
	file, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.Copy(os.Stdout, file) // os.File implements io.Reader

	// Output:
	// Hello!
}

// 3.4.3
func TestNetRead(t *testing.T) {
	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		panic(err)
	}

	conn.Write([]byte("GET / HTTP/1.0\r\nHost: ascii.jp\r\n\r\n"))
	res, err := http.ReadResponse(bufio.NewReader(conn), nil)
	// bufio.NewReader: io.Reader -> bufio.Reader

	fmt.Println(res.Header)
	defer res.Body.Close() // なくても動くけど
	io.Copy(os.Stdout, res.Body)
}

func ExampleNewBytesBuffer() {
	var buf1 bytes.Buffer
	buf2 := bytes.NewBuffer([]byte("abc123")) // ポインタ
	buf3 := bytes.NewBufferString("abc000")   // ポインタ

	fmt.Println(reflect.TypeOf(buf1))
	fmt.Println(reflect.TypeOf(buf2))
	fmt.Println(reflect.TypeOf(buf3))

	// Output:
	// bytes.Buffer
	// *bytes.Buffer
	// *bytes.Buffer
}

func TestNewReader(t *testing.T) {
	_ = bytes.NewReader([]byte{0x10, 0x20, 0x30})
}

// 3.5
func ExampleSectionReader() {
	reader := strings.NewReader("example of io.SectionReader\n")
	sectionReader := io.NewSectionReader(reader, 11, 7)
	io.Copy(os.Stdout, sectionReader)
	// Output: io.Sect
}
