package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	c *ssh.Client
}

func (client SSHClient) execute(cmd string) error{
	s, _ := client.c.NewSession()
	return s.Run(cmd)
}

// tc qdisc add dev eth0 root netem delay 100ms loss 5%
// tc qdisc change dev eth0 root netem delay 100ms loss 5%
func (client SSHClient) tc(intf string, delay string, loss string, init bool) bool {
	var cmd string
	if init {
		cmd = fmt.Sprintf("tc qdisc add dev %s root netem delay %sms loss %s%%", intf, delay, loss)
	} else {
		cmd = fmt.Sprintf("tc qdisc change dev %s root netem delay %sms loss %s%%", intf, delay, loss)
	}
	if err := client.execute(cmd); err != nil {
		return false
	}
	return true
}

// iperf3 -s
func (client SSHClient) iperf3Server() error {
	if err := client.execute("iperf3 -s"); err != nil {
		return err
	}
	// session, _ := client.c.NewSession()
	// b := bytes.Buffer{}
	// session.Stdout = &b
	// session.Run("iperf3 -s")
	// fmt.Println("----------")
	// fmt.Println(len(b.String()))
	return nil
}


// iperf3 -c x.x.x.x -t 20
func (client SSHClient) iperf3Client(ip string, seconds string) string {
	cmd := fmt.Sprintf("iperf3 -c %s -t %s -R", ip, seconds)
	fmt.Println(cmd)
	session, _ := client.c.NewSession()
	b := bytes.Buffer{}
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		return ""
	}
	return b.String()
}


func Make(host string, user string, pwd string) SSHClient {
	c := connect(host, user, pwd)
	client := SSHClient{}
	client.c = c
	return client
}