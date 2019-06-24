package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func fn_2_4() {
	file, err := os.Create("test.txt")
	defer file.Close()
	if err != nil {
		panic(err)
	}

	file.Write([]byte("os.File example\n"))
}

func useBuffer() {
	var buffer bytes.Buffer
	buffer.Write([]byte("buffer Example\n"))
	fmt.Println(buffer.String())
}

func fn_2_4_5() {
	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		panic(err)
	}
	// conn.Write([]byte("GET / HTTP/1.0\r\nHost: ascii.jp\r\n\r\n"))
	io.WriteString(conn, "GET / HTTP/1.0\r\nHost: ascii.jp\r\n\r\n")
	io.Copy(os.Stdout, conn)
}

func fn_2_4_5_http() {
	req, err := http.NewRequest("GET", "https://ascii.jp", nil)
	// (*http.Request, error)
	if err != nil {
		panic(err)
	}

	req.Write(os.Stdout)
}

func fn_2_4_6() {
	file, err := os.Create("multi.txt")
	if err != nil {
		panic(err)
	}
	writer := io.MultiWriter(file, os.Stdout)
	io.WriteString(writer, "io.MultiWriter example\n")
}

func fn_2_4_6_gzip() {}

func fn_2_4_6_bufio_writer() {
	buffer := bufio.NewWriter(os.Stdout)
	buffer.WriteString("bufio.Writer ")
	buffer.Flush()
	buffer.WriteString("example")
	// buffer.Flush()  // Flush()を呼ばないと出力されずに消滅
}

func fn_2_4_7() {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	encoder.Encode(map[string]string{
		"example": "encoding/json",
		"hello":   "world",
	})
}

func main() {
	// useBuffer()
	// fn_2_4_5()
	// fn_2_4_5_http()
	// fn_2_4_6()
	fn_2_4_6_bufio_writer()
	// fn_2_4_7()
}
