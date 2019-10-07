package main_test

import (
	"sync"
	"time"
	"fmt"
)

// 13.2 Go言語の並列処理のための道具
// goroutineの起動よりもループ変数が回るのが早いので
// println()が実行される頃にはtaskは"taskC"になってしまう
func ExampleGoroutinePassingByCaptureNG() {
	tasks := []string{"taskA", "taskB", "taskC"}
	for _, task := range tasks {
		go func() { fmt.Println(task) }()
	}

	time.Sleep(time.Second) // goroutine終了待ち

	// Output:
	// taskC
	// taskC
	// taskC
}

func ExampleGoroutinePassingArgsWait() {
	// 上の改良 goroutineの引数わたしにする
	tasks := []string{"taskA", "taskB", "taskC"}
	for i, task := range tasks {
		go func(n int, _task string) {
			time.Sleep(time.Duration(n * 1000)) // 0, 1000, 2000 milliseconds
			fmt.Println(_task)
		}(i, task)
	}

	time.Sleep(time.Second * 3) // goroutine終了待ち

	// Output:
	// taskA
	// taskB
	// taskC
}

// 13.7 syncパッケージ
// 条件変数
func ExampleCondition() {
	var mutex sync.Mutex
	cond := sync.NewCond(&mutex)

	for _, name := range []string{"A", "B", "C"} {
		go func(name string) {
			// 条件変数condを触るスレッドが複数あるのでロックを取ってから
			// mutex.Lock()
			// defer mutex.Unlock()
			cond.L.Lock()
			defer cond.L.Unlock()

			cond.Wait()
			fmt.Println(name)
		}(name)
	}

	fmt.Println("よーい")
	time.Sleep(time.Second)
	fmt.Println("どん")
	cond.Broadcast() // Broadcastするのは1箇所だけだから、ここではロック取らなくて良い
	time.Sleep(time.Second)

	// A B Cは順不同

	// Output:
	// よーい
	// どん
	// A
	// B
	// C
}