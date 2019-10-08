package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
	//シグナル初期化
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)

	// Server::Starterからもらったソケットを確認
	listeners, err := listener.ListenAll()
	if err != nil {
		panic(err)
	}
	listener := listeners[0]

	// サーバスレッドを起動する
	// * 接続要求を無限ループで待ち、クライアントが来たら非同期にレスポンス送信
	// * listner.Close()が呼ばれたとき、shutdownチャネル閉じていたら、graceful shutdown
	shutdown := make(chan bool)
	done := make(chan bool)
	go func() {

		// 受付・処理中・完了リクエスト数のカウンタ
		var accepted uint64
		var processing uint64
		var processed uint64

		var wg sync.WaitGroup
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("svr: Accept() returned error: %s\n", err)

				select {
				case <-shutdown:
					log.Println("svr: channel shutdown is closed, so waiting all goroutines finish")
					wg.Wait()
					log.Printf("svr: bye pid=%v accepted=%d processing=%d processed=%d\n", os.Getpid(), accepted, processing, processed)
					close(done)
					return
				default:
				}

				panic(err)
			}
			accepted++
			wg.Add(1)

			// log.Printf("*svr: pid=%d remoteaddr=%v\n", os.Getpid(), conn.RemoteAddr())

			go func() {
				atomic.AddUint64(&processing, 1)
				defer atomic.AddUint64(&processed, 1)
				defer wg.Done()

				time.Sleep(time.Second * 3)

				resp := http.Response{
					StatusCode: 200,
					ProtoMajor: 1,
					ProtoMinor: 0,
					Body:       ioutil.NopCloser(strings.NewReader(fmt.Sprintf("Hello from %d\n", os.Getpid()))),
				}
				resp.Write(conn)
				conn.Close()
			}()
		}
	}()

	// メインスレッドはシグナルを待つ
	select {
	case <-signals:
		// SIGTERMを受け取ったら新たなAccept()をやめて処理中のリクエストがなくなるのを待つ
		log.Println("main: shutdown start...")

		close(shutdown)
		err := listener.Close()
		if err != nil {
			log.Println("main:", err)
		}
		log.Println("main: waiting server thread shutdown")
		<-done
		log.Println("main: shutdown done")
	}

}
