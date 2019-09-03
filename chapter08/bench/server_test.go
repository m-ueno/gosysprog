package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func BenchmarkTCPServer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:18888")
		if err != nil {
			panic(err)
		}
		request, err := http.NewRequest("get", "http://localhost:18888", nil)
		if err != nil {
			panic(err)
		}
		request.Write(conn)
		response, err := http.ReadResponse(bufio.NewReader(conn), request)
		if err != nil {
			panic(err)
		}
		_, err = httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkUDSStreamServer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("unix", filepath.Join(os.TempDir(), "bench-unixdomainsocket-stream"))
		if err != nil {
			panic(err)
		}
		request, err := http.NewRequest("get", "http://localhost:18888", nil)
		if err != nil {
			panic(err)
		}
		request.Write(conn)
		response, err := http.ReadResponse(bufio.NewReader(conn), request)
		if err != nil {
			panic(err)
		}
		_, err = httputil.DumpResponse(response, true)
		if err != nil {
			panic(err)
		}
	}
}

func TestMain(m *testing.M) {
	go UnixDomainSocketStreamServer()
	go TCPServer()
	fmt.Println("prepare...")
	time.Sleep(time.Second)
	fmt.Println("prepare done")

	code := m.Run()
	os.Exit(code)
}
