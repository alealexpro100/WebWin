package main

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/ini.v1"
)

func GetPluginList(c echo.Context, dir string) error {
	var plugin_list []string
	fileslist, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Failed to read directory: %v\n", err)
		return c.String(http.StatusBadRequest, "Plugins directory not found.")
	}
	for _, f := range fileslist {
		content, err := os.ReadFile(dir + "\\" + f.Name() + "\\manifest.json")
		if err != nil {
			fmt.Printf("Failed to read file: %v\n", err)
		} else {
			plugin_list = append(plugin_list, "{ \"id\": \""+f.Name()+"\", "+string(content[1:len(content)-1])+"}")
		}
	}
	return c.String(http.StatusOK, "{\"plugins\": [ "+strings.Join(plugin_list, ", ")+"]}")
}

func PluginsAction(c echo.Context, dir string) error {
	action := c.QueryParam("action")
	switch action {
	case "list":
		return GetPluginList(c, dir)
	default:
		return c.String(http.StatusBadRequest, "Incorrect action: \""+action+"\".")
	}
}

func PluginAction(c echo.Context) error {
	name := c.Param("name")
	action := c.QueryParam("action")
	return c.String(http.StatusOK, name+" "+action)
}

func main() {
	e := echo.New()
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		os.Exit(1)
	}
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		password_hash := md5.Sum([]byte(password))
		if subtle.ConstantTimeCompare([]byte(username), []byte(cfg.Section("main").Key("user").String())) == 1 &&
			subtle.ConstantTimeCompare([]byte(hex.EncodeToString(password_hash[:])), []byte(cfg.Section("main").Key("password_md5").String())) == 1 {
			return true, nil
		}
		return false, nil
	}))
	e.Static("/", cfg.Section("main").Key("web_root").String())
	e.POST("/api/plugins", func(c echo.Context) error {
		return PluginsAction(c, cfg.Section("main").Key("web_root").String()+"\\plugins")
	})
	e.POST("/users/plugins/:name", PluginAction)
	e.Logger.Fatal(e.Start(":" + cfg.Section("main").Key("server_port").String()))
}
