package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server is running at localhost:8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		/*
			go func() {
				fmt.Printf("Accept %v\n", conn.RemoteAddr())

				request, err := http.ReadRequest(bufio.NewReader(conn))
				if err != nil {
					panic(err)
				}
				dump, err := httputil.DumpRequest(request, true)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(dump))

				response := http.Response{
					StatusCode: 200,
					ProtoMajor: 1,
					ProtoMinor: 0,
					Body:       ioutil.NopCloser(strings.NewReader("Hello, world\n")),
				}
				response.Write(conn)
				conn.Close()
			}()
		*/
		go func() {
			defer conn.Close() // さっきはなかったな...

			fmt.Printf("Accept %v\n", conn.RemoteAddr())

			// for文追加
			// v0と違って1つのgoroutineがおなじconnを回す
			// timeoutかEOFなら終了する
			//   timeoutがないと相手がいなくても固まる
			for {
				conn.SetReadDeadline(time.Now().Add(5 * time.Second))
				request, err := http.ReadRequest(bufio.NewReader(conn))
				if err != nil {
					neterr, ok := err.(net.Error) // ダウンキャスト
					if ok && neterr.Timeout() {
						fmt.Println("Timeout")
						break
					} else if err == io.EOF {
						break
					}
					panic(err)
				}

				dump, err := httputil.DumpRequest(request, true)
				if err != nil {
					panic(err)
				}
				fmt.Println(string(dump))
				content := "Hello, world\n"

				// レスポンスを書き込む
				response := http.Response{
					StatusCode:    200,
					ProtoMajor:    1,
					ProtoMinor:    1,
					ContentLength: int64(len(content)),
					// ReadCloser
					Body: ioutil.NopCloser(strings.NewReader(content)),
				}
				response.Write(conn)
			}
		}()
	}

}
