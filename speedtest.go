package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func connect(host string, user string, pwd string) *ssh.Client {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	host = fmt.Sprintf("%s:22", host)
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil
	}
	return client

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	// session, err := client.NewSession()
	// if err != nil {
	// 	log.Fatal("Failed to create session: ", err)
	// }
	// return session
	// defer session.Close()

	// // Once a Session is created, you can execute a single command on
	// // the remote side using the Run method.
	// var b bytes.Buffer
	// session.Stdout = &b
	// if err := session.Run("iperf3 -s"); err != nil {
	// 	log.Fatal("Failed to run: " + err.Error())
	// }
	// fmt.Println(b.String())
}

// func ReadCSV(path string) [][]string {
// 	f, err := os.Open(path)
// 	if err != nil {
// 		log.Fatalf("无法读取文件%s", path)
// 	}

// 	defer f.Close()

// 	r := csv.NewReader(f)
// 	data, _ := r.ReadAll()
// 	return data
// }

func searchInfo(data string) string {
	re, _ := regexp.Compile(".*/sec")
	arr := re.FindAllString(data, -1)
	for _, v := range arr {
		fmt.Println(v)
	}
	str := strings.Split(arr[len(arr)-1], " ")
	res := strings.Join(str[len(str)-2:], " ")
	return res
}

func _mainIperf3() {
	input := ReadCSV("input/iperf3.csv")
	input = input[1:]
	output := make([][]string, len(input))
	
	for idx, row := range input {
		server := Make(row[0], row[1], row[2])
		client := Make(row[3], row[4], row[5])
		if server.c == nil || client.c == nil {
			output[idx] = row
			log.Println("ssh 连接失败")
			continue
		}
		ch := make(chan string)
		server.iperf3Server()
		go func() {
			round, _ := strconv.Atoi(row[8])
			for i := 0; i < round; i++ {
				var s string
				if s = client.iperf3Client(row[6], row[7]); len(s) == 0 {
					ch <- "-"
					continue
				}
				res := searchInfo(s)
				ch <- res
			}
			close(ch)
		}()
		output[idx] = row
		for speed := range ch {
			output[idx] = append(output[idx], speed)
		}
	}
	fmt.Println(output)
	t := time.Now().Format("2006-01-02-15:04:05")
	filename := fmt.Sprintf("out/iperf3-%s.csv", t)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("无法创建文件")
	}
	defer file.Close()
	// file.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(file)
	header := []string{"server", "username", "pwd", "client", "username", "pwd", "server_ip", "time", "round", "speed"}
	w.Write(header) // header
	w.Flush()
	for _, o := range output {
		w.Write(o)
		w.Flush()
	}
}
