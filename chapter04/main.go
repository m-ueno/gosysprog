package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 4.3.1
	ch := make(chan os.Signal, 1)     // os.Signalが流れるチャネル
	signal.Notify(ch, syscall.SIGINT) // SIGINTをチャネルに送出

	fmt.Println("Waiting SIGINT...")
	<-ch
	fmt.Println("SIGINT arrived")
}
