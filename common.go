package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
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
		times:      cfg.Section("global").Key("times").RangeInt(10, 1, 15),
		unit:       cfg.Section("global").Key("unit").MustString("M"),
		format:     cfg.Section("global").Key("format").String(),
		pipein:     cfg.Section("server").Key("pipein").String(),
		pipeout:    cfg.Section("server").Key("pipeout").String(),
	}
	fmt.Println(env)
	return env
}

func GetValueFromINI(filename string, section string, key string) string {
	cfg, err := ini.Load(filename)
	if err != nil {
		log.Fatalf("无法读取配置文件%s", filename)
	}

	return cfg.Section(section).Key(key).String()
}
