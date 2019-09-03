package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	conn, err := net.Dial("unix",
		filepath.Join(os.TempDir(), "unixdomainsocket-sample"))
	if err != nil {
		panic(err)
	}
	request, err := http.NewRequest("get", "http://localhost:8888", nil)
	if err != nil {
		panic(err)
	}
	request.Write(conn)

	fmt.Println("done")
}
