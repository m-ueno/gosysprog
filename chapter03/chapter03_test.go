package main_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"testing"
)

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

func TestNewBytesBuffer(t *testing.T) {
	var buf1 bytes.Buffer
	buf2 := bytes.NewBuffer([]byte("abc123"))
	buf3 := bytes.NewBufferString("abc000")

	fmt.Println(buf1, buf2, buf3)
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
