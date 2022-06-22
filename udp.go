package main

import (
	"log"
	"net"
)

func _main() {
	l, _ := net.ListenUDP("udp", &net.UDPAddr{
		IP: net.ParseIP("127.0.0.1"),
		Port: 7777,
	})
	defer l.Close()
	for {
		buf := [1024]byte{}
		n, addr, _ := l.ReadFromUDP(buf[:])
		log.Printf("src: %v, n: %v, data: %v", addr, n, string(buf[:n]))
		l.WriteToUDP(buf[:n], addr)
	}
}
