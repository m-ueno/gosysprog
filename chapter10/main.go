package main

import (
	"log"
	"sync"
	"syscall"
	"time"
)

type FileLock struct {
	l  sync.Mutex
	fd int
}

func NewFileLock(filename string) *FileLock {
	if filename == "" {
		panic("filename needed")
	}
	fd, err := syscall.Open(filename, syscall.O_CREAT|syscall.O_RDONLY, 0750)
	if err != nil {
		panic(err)
	}
	return &FileLock{fd: fd}
}

func (m *FileLock) Lock() {
	m.l.Lock()
	if err := syscall.Flock(m.fd, syscall.LOCK_EX); err != nil {
		panic(err)
	}
}

func (m *FileLock) Unlock() {
	if err := syscall.Flock(m.fd, syscall.LOCK_UN); err != nil {
		panic(err)
	}
	m.l.Unlock()
}

/*
syscall.Flock()のモードフラグ / リソースのロック
LOCK_SH 共有ロック
LOCK_EX 排他ロック
LOCK_UN ロック解除
LOCK_NB ノンブロッキングモード Go言語の場合は並行処理が比較的簡単書けるので...
*/
// 10.2.3

func sub() {
	l := NewFileLock("main.go")
	log.Println("try locking")
	l.Lock()
	log.Println("locked")
	time.Sleep(10 * time.Second)
	l.Unlock()
	log.Println("unlocked")
}

func main() {
	done := make(chan bool, 2)
	go func() {
		sub()
		done <- true
	}()
	go func() {
		sub()
		done <- true
	}()
	log.Println("* waiting...")
	<-done
	<-done
	log.Println("* done")
}
