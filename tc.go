package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)
const G = 1000
const M = 1
const K = 0.001

const N = 10 // 测试次数

type Env struct {
	serverIP string
	clientIP string
	serverPwd string
	clientPwd string
	server_ip string
	serverUser string
	clientUser string

	server SSHClient
	client SSHClient
}

func MakeEnv() *Env {
	en := Env{
		serverIP : "172.23.75.44",
		clientIP : "172.23.75.45",
		serverPwd : "root@SFtel",
		clientPwd : "root@SFtel",
		server_ip : "10.0.99.1",
		serverUser : "root",
		clientUser : "root",
	}
	en.server = Make(en.serverIP, en.serverUser, en.serverPwd)
	en.client = Make(en.clientIP, en.clientUser, en.clientPwd)

	// en.server.iperf3Server()
	return &en
}


func _mainTC() {
	env := MakeEnv()
	
	input := ReadCSV("input/tc.csv")
	fmt.Println("test output")

	for row := 1; row < len(input[0]); row++ {
		for col := 1; col < len(input); col++ {
			loss := input[row][0]
			delay := input[0][col]
			var num float64
			// 测试n次，取最大值
			datas := make([]float64, N)
			max := 0.0
			min := math.MaxFloat64
			maxIdx := 0
			minIdx := 0
			for i := 0; i < N; i++ {
				_, n := env.once(delay, loss)
				// if n > num {
				// 	num = n
				// 	speed = s
				// }
				datas[i] = n
				if n != 0 && n > max {
					max = n
					maxIdx = i
				}
				if n != 0 && n < min {
					min = n
					minIdx = i
				}
			}
			var sum float64
			var div int
			for j := 0; j < N; j++ {
				if datas[j] == 0 || j == maxIdx || j == minIdx{
					continue
				}
				sum += datas[j]
				div++
			}
			res := sum / float64(div)
			log.Printf("loss: %s, delay: %sms.speed: %.3f Mbps", loss, delay, num)
			input[row][col] = fmt.Sprintf("%.3f", res)
		}
	}
	fmt.Println(input)
	filename := generateFilename("tc-iperf3")
	writeCSV(filename, input)
}

func writeCSV(filename string, data [][]string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("无法创建文件")
	}
	defer file.Close()
	// file.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(file)
	for _, o := range data {
		w.Write(o)
		w.Flush()
	}
}

func generateFilename(basename string) string {
	t := time.Now().Format("2006-01-02-15:04:05")
	filename := fmt.Sprintf("out/%s-%s.csv", basename, t)
	return filename
}

func speed2Int(str string) float64 {
	spl := strings.Split(str, " ")
	numStr := spl[0]
	num, _ := strconv.ParseFloat(numStr, 32)
	if strings.Contains(str, "G") {
		num = num * G
	} else if strings.Contains(str, "M") {
		num = num * M
	} else if strings.Contains(str, "K") {
		num = num * K
	}
	return num
}

func (env *Env) once(delay string, loss string) (string, float64) {
	log.Printf("loss: %s%%, delay: %sms", loss, delay)
	if s := env.server.tc("eth0", delay, loss, false); !s {
		log.Println("tc error")
		return "", 0
	}
	var str string
	// 测试时长 20s
	if str = env.client.iperf3Client(env.server_ip, "10"); len(str) == 0 {
		log.Println("iperf3 -c error")
		return "", 0
	}
	speed := searchInfo(str)
	return speed, speed2Int(speed)
}