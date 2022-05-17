package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	cfg    *ConfigServer
	worker *jobs
}

func (s Server) PluginsAction(c echo.Context) error {
	action := c.QueryParam("action")
	switch action {
	case "list":
		return c.String(http.StatusOK, "{\"plugins\": [ "+strings.Join(plugin_list, ", ")+"]}")
	case "get_status":
		id, err := strconv.Atoi(c.QueryParam("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Incorrect id!")
		}
		var content struct {
			Id       int    `json:"id"`
			Status   string `json:"status"`
			ExitCode int    `json:"exitcode"`
			Stdout   string `json:"stdout"`
			Stderr   string `json:"stderr"`
		}
		content.Id = id
		content.Status = s.worker.list[id].Status
		content.ExitCode = s.worker.list[id].ExitCode
		content.Stdout = s.worker.list[id].Stdout
		content.Stderr = s.worker.list[id].Stderr
		return c.JSON(http.StatusOK, content)
	case "jobs_clear":
		err := s.worker.ClearJobs()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, "ok")
	default:
		return c.String(http.StatusBadRequest, "Incorrect action: \""+action+"\".")
	}
}

func (s Server) PluginDo(c echo.Context) error {
	name := c.Param("name")
	action := c.QueryParam("action")
	switch action {
	case "add":
		param := c.QueryParam("param")
		return c.String(http.StatusOK, strconv.Itoa(s.worker.AddJob(s.cfg.WebRoot, name, param)))
	default:
		return c.String(http.StatusBadRequest, "Incorrect action: \""+action+"\".")
	}
}

func (s Server) InternalAction(c echo.Context) error {
	action := c.QueryParam("action")
	switch action {
	case "set_auth":
		user := c.QueryParam("user")
		pass_hash := c.QueryParam("pass_hash")
		if user == "" || pass_hash == "" {
			return c.String(http.StatusBadRequest, "Param user or pass_hash is empty.")
		}
		s.cfg.set_auth(user, pass_hash)
		return c.String(http.StatusOK, "ok")
	default:
		return c.String(http.StatusBadRequest, "Incorrect action: \""+action+"\".")
	}
}

func (s Server) start() {
	ReloadPluginList(s.cfg.WebRoot + "\\plugins")
	e := echo.New()
	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		password_hash := sha256.Sum256([]byte(password))
		if subtle.ConstantTimeCompare([]byte(username), []byte(s.cfg.User)) == 1 &&
			subtle.ConstantTimeCompare([]byte(hex.EncodeToString(password_hash[:])), []byte(s.cfg.PasswordSHA256)) == 1 {
			return true, nil
		}
		return false, nil
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Static("/", s.cfg.WebRoot)
	e.POST("/api/plugins", func(c echo.Context) error {
		return s.PluginsAction(c)
	})
	e.POST("/api/plugins/:name", s.PluginDo)
	e.POST("/api/internal", s.InternalAction)
	e.Logger.Fatal(e.Start(":" + s.cfg.Port))
}
