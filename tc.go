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

type Env struct {
	serverIP   string
	clientIP   string
	serverPwd  string
	clientPwd  string
	server_ip  string
	serverUser string
	clientUser string

	server SSHClient
	client SSHClient

	times  int
	unit   string
	format string
	pipein string
	pipeout string
}

func MakeEnv() *Env {
	en := ReadIni("config.ini")
	en.server = Make(en.serverIP, en.serverUser, en.serverPwd)
	en.client = Make(en.clientIP, en.clientUser, en.clientPwd)

	// en.server.iperf3Server()
	return &en
}

func _mainTC() {
	env := MakeEnv()
	log.Println("读取配置文件")
	input := ReadCSV("input/tc.csv")
	log.Println("读取输入文件")
	fmt.Println(input)
	// 创建文件保存结果
	filename := generateFilename("tc-iperf3")
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("无法创建文件 ", err)
	}
	defer file.Close()
	w := csv.NewWriter(file)
	w.Write(input[0])
	w.Flush()

	for row := 1; row < len(input); row++ {
		for col := 1; col < len(input[0]); col++ {
			loss := input[row][0]
			delay := input[0][col]
			// 测试n次
			datas := make([]float64, env.times)
			max := 0.0
			min := math.MaxFloat64
			maxIdx := 0
			minIdx := 0
			// 设置delay loss
			if err := env.server.ipfw(env.pipein, loss, delay); !err {
				log.Println("ipfw error")
				continue
			}
			if err := env.server.ipfw(env.pipeout, loss, delay); !err {
				log.Println("ipfw error")
				continue
			}
			for i := 0; i < env.times; i++ {
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
			for j := 0; j < env.times; j++ {
				if datas[j] == 0 || j == maxIdx || j == minIdx {
					continue
				}
				sum += datas[j]
				div++
			}
			res := sum / float64(div)
			fmt.Println(datas)
			log.Printf("loss: %s, delay: %sms. ==> speed: %.3f/%.3f/%.3f Mbps", loss, delay, datas[minIdx], res, datas[maxIdx])
			if env.format == "avg" {
				input[row][col] = fmt.Sprintf("%.3f", res)
			} else {
				input[row][col] = fmt.Sprintf("%.3f/%.3f/%.3f", datas[minIdx], res, datas[maxIdx])
			}
			fmt.Println(input)
		}
		w.Write(input[row])
		w.Flush()
	}
	fmt.Println(input)
}

func WriteCSV(filename string, data [][]string) {
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

// unit G/M/K
func speed2Int(str string, unit string) float64 {
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
	switch unit {
	case "G":
		num = num / 1000
	case "K":
		num = num * 1000
	default:
	}
	return num
}

func (env *Env) once(delay string, loss string) (string, float64) {
	// lossint, _ := strconv.ParseFloat(loss, 64)
	// delayint, _ := strconv.ParseFloat(delay, 64)
	// halfLoss := fmt.Sprintf("%.3f", lossint/2)
	// halfdelay := fmt.Sprintf("%.3f", delayint/2)
	// log.Printf("loss: %s%%, delay: %sms", loss, delay)
	// if success := env.server.ipfw(loss, delay); !success {
	// 	log.Println("ipfw error")
	// 	return "", 0
	// }
	// if s := env.client.tc(env.clientDev, halfdelay, halfLoss, false); !s {
	// 	log.Println("tc client error")
	// 	return "", 0
	// }
	var str string
	// 测试时长 10s
	if str = env.client.iperf3Client(env.server_ip, "10"); len(str) == 0 {
		log.Println("iperf3 -c error")
		return "", 0
	}
	speed := searchInfo(str)
	return speed, speed2Int(speed, env.unit)
}
