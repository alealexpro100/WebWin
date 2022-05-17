package main

import (
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigServer struct {
	Port           string
	WebRoot        string
	User           string
	PasswordSHA256 string
	ConfigFile     *ini.File
}

func (cfg *ConfigServer) init() {
	var err error
	cfg.ConfigFile, err = ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}
	cfg.get()
}

func (cfg *ConfigServer) get() {
	cfg.Port = cfg.ConfigFile.Section("main").Key("web_port").String()
	cfg.WebRoot = cfg.ConfigFile.Section("main").Key("web_root").String()
	cfg.User = cfg.ConfigFile.Section("main").Key("user").String()
	cfg.PasswordSHA256 = cfg.ConfigFile.Section("main").Key("password_sha256").String()
}

func (cfg *ConfigServer) set_auth(user, pass_hash string) {
	cfg.ConfigFile.Section("main").Key("user").SetValue(user)
	cfg.ConfigFile.Section("main").Key("password_sha256").SetValue(pass_hash)
	cfg.User = user
	cfg.PasswordSHA256 = pass_hash
	cfg.ConfigFile.SaveTo("config.ini")
}
