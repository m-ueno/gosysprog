package main_test

import (
	"errors"
	"fmt"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"testing"
	"time"
)

func fib(k int) int {
	if k <= 1 {
		return 1
	}

	return fib(k-1) + fib(k-2)
}

func fibWorker(id int, tasks chan int, wg *sync.WaitGroup) {
	// closeされるまでループ
	for k := range tasks {
		fmt.Println(id, k, fib(k))
		wg.Done()
	}
}

// 14.2.7 ワーカープール
// fib(1), ..., fib(N)をCPUコア数だけ並列計算
func ExampleWorkerPool() {
	n := 43
	tasks := make(chan int, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		tasks <- i
	}

	for i := 0; i < runtime.NumCPU()/2; i++ {
		// go func(ch chan int, wg *sync.WaitGroup) {fibWorker(ch, wg) }(tasks, &wg)
		go fibWorker(i, tasks, &wg)
	}

	// 全てのワーカーが終了する. 無くても動くがブロック中のgoroutineが終わらない
	close(tasks)
	wg.Wait() // 終了待ち
	fmt.Println("bye")

	// Output:
	// 0 1
}

// Goroutineを確実に終了する (timeoutつきfib. 同期)
func fibWithTimeout(n int) (int, error) {
	// 計算スレッドとメインスレッドの2個を並列実行
	result := make(chan int)

	// 計算スレッド
	// このスレッドはいつ終わる?
	go func() {
		result <- fib(n)
	}()

	select {
	case k := <-result:
		return k, nil
	case <-time.After(5 * time.Second):
		close(result) // いる？
		return 0, errors.New("timed out")
	}
}

func TestFibTimeout(t *testing.T) {
	if k, _ := fibWithTimeout(1); k != 1 {
		t.Fatalf("fib(1) must be 1")
	}
	if k, _ := fibWithTimeout(2); k != 2 {
		t.Fatalf("fib(2) must be 2, but %d", k)
	}
	if _, err := fibWithTimeout(50); err != nil {
		t.Logf("fib(50) should be timeout: %s", err)
	}
}

func ExampleTryToWriteClosedChannel() {
	func() {
		ch := make(chan bool)
		defer close(ch)

		go func() {
			fmt.Println("gofunc")
			time.Sleep(time.Second)
			ch <- true
		}()

		select {
		case <-time.After(time.Millisecond * 1500):
			fmt.Println("timeout")
			break
		case <-ch:
			fmt.Println("done")
		}
	}()
	// Output:
	// gofunc
	// done
}

func fibChannel(n int) chan int {
	ch := make(chan int)
	go func() {
		ch <- fib(n)
	}()

	return ch
}

func leak() {
	// 読み込みでブロックされる例
	ch := make(chan int)
	<-ch
}

func leak2() chan bool {
	// 書き込みでブロックされる例
	ch := make(chan bool)
	go func() {
		ch <- true
		fmt.Println("done")
	}()
	return ch
}

// 終わらないgoroutine (リーク) を多数起動するとOSに殺される
func ExampleLeak() {
	for i := 0; i < 1e7; i++ {
		_ = leak2()
	}
	// Output:
}

// https://youtu.be/SmoM1InWXr0?t=935
func ExampleNilChannel() {
	c := make(chan int)
	d := make(chan bool)

	go func(src chan int) {
		for {
			select {
			case src <- 42:
			case <-d:
				src = nil
			default:
				fmt.Println("default")
				return
			}
		}
	}(c)

	fmt.Printf("%d\n", <-c)
	fmt.Printf("%d\n", <-c)
	// d <-true
	fmt.Printf("%d\n", <-c)
	// Output:
	// 42
	// 42
	// 42

}
