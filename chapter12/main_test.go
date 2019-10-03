package main_test

import (
	"fmt"
	"os"
	"os/exec"
)


func ExampleSendSignal() {
	// spawn child process
	cmd := exec.Command("tail", "-f", "/dev/null")
	cmd.Start()

	// send signal
	// Cmd構造体のProcessメンバ
	cmd.Process.Signal(os.Interrupt)

	fmt.Println("done")
	// Output: done
}