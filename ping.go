package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

type PingResult struct {
	Loss string
	Min  string
	Avg  string
	Max  string
	Mdev string
}

func ReadCSV(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("无法读取文件%s", path)
	}

	defer f.Close()

	r := csv.NewReader(f)
	data, _ := r.ReadAll()
	return data
}


func ping(count string, ip string) PingResult {
	sysType := runtime.GOOS
	var flag string
	if sysType == "windows" {
		flag = "-n"
	} else {
		flag = "-c"
	}
	cmd := exec.Command("ping", flag, count, ip)
	stdout := &bytes.Buffer{}
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	res := PingResult{
		Loss: "-",
		Min:  "-",
		Avg:  "-",
		Max:  "-",
		Mdev: "-",
	}
	cmd.Run()
	// loss
	re, _ := regexp.Compile("\\d+\\.?\\d*%")
	s := re.FindString(string(stdout.Bytes()))
	res.Loss = s

	// min/max/avg
	var rttre *regexp.Regexp
	if sysType == "windows" {
		rttre, _ = regexp.Compile("最短.*")
		s = rttre.FindString(string(stdout.Bytes()))
		pattern, _ := regexp.Compile("\\d+\\.?\\d*ms")
		arr := pattern.FindAllString(s, -1)
		res.Min = arr[0]
		res.Max = arr[1]
		res.Avg = arr[2]
	} else { // linux freebsd
		rttre, _ = regexp.Compile("min/avg/max/.* = (.*)/(.*)/(.*)/(.*)")
		s = rttre.FindString(string(stdout.Bytes()))
		if len(s) != 0 {
			idx := strings.Index(s, "=")
			s = s[idx+2 : len(s)-3]
			arr := strings.Split(s, "/")
			res.Min = arr[0]
			res.Avg = arr[1]
			res.Max = arr[2]
			res.Mdev = arr[3]
		}
	}
	return res
}


func _mainPing() {
	data := ReadCSV("input/ping.csv")
	ips := data[0][1:]
	// rounds := data[1][1:]
	ttlTimes := data[2][1:]
	// intervals := data[4][1:]

	var mu sync.Mutex
	wg := sync.WaitGroup{}
	writeData := []PingResult{}
	for idx, ip := range ips {
		// args := fmt.Sprintf("-c %s %s", ttlTimes[idx], ip)
		fmt.Println("ping ", ip)
		wg.Add(1)
		go func(i int, ip string) {
			res := ping(ttlTimes[i], ip)
			mu.Lock()
			defer mu.Unlock()
			writeData = append(writeData, res)
			wg.Done()
		}(idx, ip)
	}

	wg.Wait()
	t := time.Now().Format("2006-01-02-15:04:05")
	filename := fmt.Sprintf("out/ping-%s.csv", t)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("无法创建文件")
	}
	defer file.Close()
	// file.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(file)
	w.Write([]string{"ip", "loss", "min", "avg", "max", "mdev"}) // header
	w.Flush()
	for i, ps := range writeData {
		row := []string{ips[i], ps.Loss, ps.Min, ps.Avg, ps.Max, ps.Mdev}
		fmt.Println(row)
		w.Write(row)
	}
	w.Flush()
}
