package main_test

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/shirou/gopsutil/process"
)

func ExampleGoPsUtil() {
	p, _ := process.NewProcess(int32(os.Getpid()))
	name, _ := p.Name()
	cmd, _ := p.Cmdline()
	fmt.Printf("name: %s, cmd: %s\n", name, cmd)

	// Output:
	// name: chapter11.test, cmd: /tmp/go-build825168190/b001/chapter11.test
}

func ExampleExec() {
	cmd := exec.Command("sleep", "1")
	err := cmd.Run() // Start and Wait
	if err != nil {
		panic(err)
	}
	state := cmd.ProcessState
	fmt.Printf("%s\n", state.String())
	fmt.Printf("  Pid: %d\n", state.Pid())
	fmt.Printf("  System: %v\n", state.SystemTime())
	fmt.Printf("  User: %v\n", state.UserTime())
	// Output:
	// exit status 0
	//   Pid: 24508
	//   System: 0s
	//   User: 1.064ms
}