package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: netutil ping or netutil iperf3")
		return
	}
	cmd := os.Args[1]

	if cmd == "ping" {
		_mainPing()
	} else if cmd == "iperf3" {
		_mainIperf3()
	} else if cmd == "tc" {
		_mainTC()
	} else {
		fmt.Println("command error")
	}
}
