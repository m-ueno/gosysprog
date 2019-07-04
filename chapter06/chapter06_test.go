package chapter06

import (
	"testing"
	"time"
)

// 6.5
func TestHelloServer(t *testing.T) {
	go func() {
		HelloServer()
	}()
	time.Sleep(10 * time.Millisecond)
	TCPClient()
	time.Sleep(time.Second)
}
