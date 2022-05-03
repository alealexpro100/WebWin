package main

import (
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"

	"github.com/spf13/viper"
)

type Login struct {
	User     string `form:"user"`
	Password string `form:"password"`
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic("Coudn't read config file!")
	}
	router := gin.Default()
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		viper.GetString("main.user"): viper.GetString("main.pass"),
	}))
	authorized.StaticFS("/", http.Dir(viper.GetString("main.web_root")))
	authorized.POST("/api/:plugin/*action", func(c *gin.Context) {
		cCopy := c.Copy()
		{
			var result string
			var output string
			plugin := cCopy.Param("plugin")
			action := cCopy.Param("action")
			options := cCopy.DefaultPostForm("options", "")
			switch os := action; os {
			case "/status":
				result = "unknown"
			case "/run":
				cmd := exec.Command("powershell", "-NoLogo", "-NoProfile", viper.GetString("main.web_root")+"\\plugins\\"+plugin+"\\exec.ps1")
				stdout, err := cmd.CombinedOutput()
				print(string(stdout))
				if err != nil {
					print(err)
					result = err.Error()
				} else {
					result = "success"
				}
				output = string(stdout)
			default:
				result = "wrong action"
			}
			c.JSON(200, gin.H{
				"result":  result,
				"plugin":  plugin,
				"options": options,
				"action":  action,
				"output":  output,
			})
		}
	})

	router.Run(":" + viper.GetString("main.port"))
}
