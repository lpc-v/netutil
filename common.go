package main

import (
	"fmt"
	"log"

	"gopkg.in/ini.v1"
)

func ReadIni(filename string) Env {
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatalf("无法读取配置文件%s", filename)
	}

	env := Env{
		serverIP:   cfg.Section("server").Key("ip").String(),
		clientIP:   cfg.Section("client").Key("ip").String(),
		serverPwd:  cfg.Section("server").Key("password").String(),
		clientPwd:  cfg.Section("client").Key("password").String(),
		server_ip:  cfg.Section("client").Key("server_ip").String(),
		serverUser: cfg.Section("server").Key("username").String(),
		clientUser: cfg.Section("client").Key("username").String(),
		serverDev:  cfg.Section("server").Key("dev").String(),
		clientDev:  cfg.Section("client").Key("dev").String(),
		times:      cfg.Section("global").Key("times").RangeInt(10, 1, 15),
		unit:       cfg.Section("global").Key("unit").MustString("M"),
	}
	fmt.Println(env)
	return env
}
