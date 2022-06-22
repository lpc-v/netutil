package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestSSH(t *testing.T) {
	server := connect("192.168.0.2", "root", "pfsense")
	session, _ := server.NewSession()
	session.Stdout = os.Stdout
	session.Run("ping 8.8.8.8")
}

func TestAuto(t *testing.T) {
	_mainIperf3()
}

func TestMake(t *testing.T) {
	client := Make("39.106.250.109", "root", "root@SFtel")
	if client.tc("eth0", "30", "5", true) {
		log.Println("success")
	}
}

func TestIperf3(t *testing.T) {
	server := Make("39.106.250.109", "root", "root@SFtel")
	server.iperf3Server()
	client := Make("60.205.106.228", "root", "root@SFtel")
	s := client.iperf3Client("172.20.110.221", "10")
	fmt.Println(s)
}

func TestDivide(t *testing.T) {
	a := 9
	b := float64(a / 2.0)
	fmt.Printf("%T, %v", b, b)
}

func TestReadIni(t *testing.T) {
	env := ReadIni("config.ini")
	fmt.Println(env)
}