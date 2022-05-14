package main

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var WebRoot string

type Server struct {
	Port        string
	User        string
	PasswordMD5 string
}

func PluginsAction(c echo.Context) error {
	action := c.QueryParam("action")
	switch action {
	case "list":
		return c.String(http.StatusOK, "{\"plugins\": [ "+strings.Join(plugin_list, ", ")+"]}")
	case "jobs_clear":
		err := ClearJobs()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		} else {
			return c.String(http.StatusOK, "ok")
		}
	default:
		return c.String(http.StatusBadRequest, "Incorrect action: \""+action+"\".")
	}
}

func PluginDo(c echo.Context) error {
	name := c.Param("name")
	action := c.QueryParam("action")
	param := c.QueryParam("param")
	switch action {
	case "add":
		return c.String(http.StatusOK, strconv.Itoa(AddJob(WebRoot, name, param)))
	case "get_status":
		id, err := strconv.Atoi(param)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Incorrect id!")
		} else {
			var content struct {
				Status   string `json:"status"`
				ExitCode int    `json:"exitcode"`
				Stdout   string `json:"stdout"`
				Stderr   string `json:"stderr"`
			}
			content.Status = jobs[id].Status
			content.ExitCode = jobs[id].ExitCode
			content.Stdout = jobs[id].Stdout
			content.Stderr = jobs[id].Stderr
			return c.JSON(http.StatusOK, content)
		}
	default:
		return c.String(http.StatusBadRequest, "Incorrect action: \""+action+"\".")
	}
}

func (s Server) start() {
	e := echo.New()
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		password_hash := md5.Sum([]byte(password))
		if subtle.ConstantTimeCompare([]byte(username), []byte(s.User)) == 1 &&
			subtle.ConstantTimeCompare([]byte(hex.EncodeToString(password_hash[:])), []byte(s.PasswordMD5)) == 1 {
			return true, nil
		}
		return false, nil
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Static("/", WebRoot)
	e.POST("/api/plugins", func(c echo.Context) error {
		return PluginsAction(c)
	})
	e.POST("/api/plugins/:name", PluginDo)
	e.Logger.Fatal(e.Start(":" + s.Port))
}
