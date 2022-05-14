package main

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}
	ReloadPluginList(cfg.Section("main").Key("web_root").String() + "\\plugins")
	server := new(Server)
	server.Port = cfg.Section("main").Key("web_port").String()
	WebRoot = cfg.Section("main").Key("web_root").String()
	server.User = cfg.Section("main").Key("user").String()
	server.PasswordMD5 = cfg.Section("main").Key("password_md5").String()
	server.start()
}
