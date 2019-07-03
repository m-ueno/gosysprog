package main_test

import (
	"context"
	"fmt"
	"time"
)

// 通知のパターン：
// 受信側がチャネルをブロッキングしておいて、適当な値を送信する
// 受信側がチャネルをブロッキングしておいて、チャネルを閉じる

func ExampleChanNotify() {
	fmt.Println("start sub()")

	done := make(chan bool)
	go func() {
		fmt.Println("sub() is finished")
		done <- true
	}()
	<-done
	fmt.Println("finish")

	// Output:
	// start sub()
	// sub() is finished
	// finish
}

func ExampleCtxNotify() {
	fmt.Println("start sub()")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Println("sub() is finished")
		cancel()
	}()

	<-ctx.Done()
	fmt.Println("finish")

	// Output:
	// start sub()
	// sub() is finished
	// finish
}

// Q4.1
func ExampleTimer() {
	ch := time.After(3 * time.Second)
	<-ch
	fmt.Println("Bang!")
	// Output: Bang!
}
