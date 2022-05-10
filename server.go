package main

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/ini.v1"
)

func GetPluginList(c echo.Context, dir string) error {
	var plugin_list []string
	fileslist, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	for _, f := range fileslist {
		plugin_list = append(plugin_list, "\""+f.Name()+"\"")
	}
	return c.String(http.StatusOK, "{\"plugins\": ["+strings.Join(plugin_list, ", ")+"]}")
}

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(cfg.Section("main").Key("user").String())) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(cfg.Section("main").Key("password").String())) == 1 {
			return true, nil
		}
		return false, nil
	}))
	e.Static("/", cfg.Section("main").Key("web_root").String())
	e.POST("/api/clean_memory", func(c echo.Context) error {
		runtime.GC()
		return c.String(http.StatusOK, "cleanup: success")
	})
	e.POST("/api/plugins", func(c echo.Context) error {
		return GetPluginList(c, cfg.Section("main").Key("web_root").String()+"\\plugins")
	})
	e.Logger.Fatal(e.Start(":" + cfg.Section("main").Key("server_port").String()))
}
