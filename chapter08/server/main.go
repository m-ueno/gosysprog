package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	path := filepath.Join(os.TempDir(), "unixdomainsocket-sample")

	listener, err := net.Listen("unix", path)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Server is running at " + path)
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			fmt.Println("Accept")

			response := http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("Hello client!")),
			}
			response.Write(conn)
			conn.Close()
		}()
	}
}
