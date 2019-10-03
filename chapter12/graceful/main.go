package main

import (
	// "context"
	"fmt"
	"github.com/lestrrat/go-server-starter/listener"
	"io/ioutil"
	"net/http"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	//シグナル初期化
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	// Server::Starterからもらったソケットを確認
	listeners, err := listener.ListenAll()
	if err != nil {
		panic(err)
	}

	connChan := make(chan net.Conn)
	done := make(chan bool, 1)
	wg := &sync.WaitGroup{}

	go func() {
		for {
			conn, _ := listeners[0].Accept()
			connChan <- conn
		}
	}()

	go func() {
		<-signals
		done <- true
	}()

L:
	for {
		select {
		case conn := <-connChan:
			fmt.Printf("accept new connection %v\n", conn)

			wg.Add(1)

			go func() {
				fmt.Printf("Remote: %v\n", conn.RemoteAddr())

				time.Sleep(time.Second)

				resp := http.Response{
					StatusCode: 200,
					ProtoMajor: 1,
					ProtoMinor: 0,
					Body:       ioutil.NopCloser(strings.NewReader("Hello!\n")),
				}
				resp.Write(conn)
				conn.Close()
				wg.Done()
			}()
		case <-done:
		// SIGTERMを受け取ったら新たなAccept()をやめて処理中のリクエストがなくなるのを待つ
			fmt.Println("drain start...")
			wg.Wait()
			fmt.Println("drain done")
			break L
		}
	}

	fmt.Println("* bye")
}
