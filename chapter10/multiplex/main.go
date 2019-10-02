package main

import (
	"fmt"
	"syscall"
)

func main(){
	efd, err := syscall.EpollCreate(100)
	if err != nil {
		panic(err)
	}
	fd, err := syscall.Open("./test", syscall.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	ev1 := syscall.EpollEvent{
		Events: 1,
		Fd: int32(fd),
		Pad: 1,
	}

	for {
		events := make([]syscall.EpollEvent, 10)

	}
}