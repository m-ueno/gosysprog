package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/lestrrat/go-server-starter/listener"
)

func main() {
	log.Printf("Parent pid=%d\n", os.Getppid())

	//シグナル初期化
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	// Server::Starterからもらったソケットを確認
	listeners, err := listener.ListenAll()
	if err != nil {
		panic(err)
	}

	// 受付・処理中・完了リクエスト数のカウンタ
	var accepted uint64
	var processing uint64
	var processed uint64

	var wg sync.WaitGroup

	connChan := make(chan net.Conn)
	die := make(chan bool, 1)
	go func() {
		listener := listeners[0]
		defer listener.Close()

		for {
			conn, _ := listener.Accept()
			accepted++
			connChan <- conn

			// NOTE:
			// このように書いてもgoroutine外から終わらせることはできない
			select {
			case <-die:
				return
			default:
			}
		}
	}()

L:
	for {
		select {
		case conn := <-connChan:
			log.Printf("* pid=%d remoteaddr=%v\n", os.Getpid(), conn.RemoteAddr())

			wg.Add(1)
			processing++

			go func() {
				time.Sleep(time.Second)

				resp := http.Response{
					StatusCode: 200,
					ProtoMajor: 1,
					ProtoMinor: 0,
					Body:       ioutil.NopCloser(strings.NewReader(fmt.Sprintf("Hello from %d\n", os.Getpid()))),
				}
				resp.Write(conn)
				conn.Close()
				atomic.AddUint64(&processed, 1)
				wg.Done()
			}()
		case <-signals:
			// SIGTERMを受け取ったら新たなAccept()をやめて処理中のリクエストがなくなるのを待つ
			log.Println("drain start...")
			die <-true // 気休め. Accept()を止めることはできない
			wg.Wait()
			log.Println("drain done")
			break L
		}
	}

	log.Printf("bye pid=%v accepted=%d processing=%d processed=%d\n", os.Getpid(), accepted, processing, processed)
}
