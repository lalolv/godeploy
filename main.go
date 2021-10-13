package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ccding/go-config-reader/config"
	"github.com/gin-gonic/gin"
	"github.com/lalolv/goapis"
	"github.com/lalolv/godeploy/deploy"
	"github.com/lalolv/godeploy/glob"
	"github.com/lalolv/godeploy/utils"
)

func main() {
	// 读取配置文件
	c := config.NewConfig("app.conf")
	err := c.Read()
	// error handle
	if err != nil {
		fmt.Println(err.Error())
	}

	glob.ShellPath = c.Get("shell", "path")

	// 设置日志
	logPath, _ := utils.PathExists("log")
	if logPath {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
		// 写入文件
		logOut, err := os.OpenFile("log/api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm|os.ModeAppend)
		if err != nil {
			fmt.Println("err", err)
		}
		log.Logger = log.Output(logOut)
	}

	// slack
	logStartup := fmt.Sprintf("[%s] Deploy service is done.", "Startup")
	glob.SlackToken = c.Get("slack", "token")
	glob.SlackChannel = c.Get("slack", "channel")
	if glob.SlackToken != "" {
		glob.SlackUtil = goapis.SlackUtil{Token: glob.SlackToken}
		utils.SendSlack(logStartup)
	}

	// slash
	glob.SlashToken = c.Get("slash", "token")
	glob.SlashCommand = fmt.Sprintf("/%s", c.Get("slash", "command"))

	// Creates a gin router with default middlewares:
	// logger and recovery (crash-free) middlewares
	router := gin.New()
	// routers
	router.POST("/deploy/:shell", deploy.Uniform)
	router.POST("/deploy/slash", deploy.SlashCommands)
	// log
	log.Log().Msg(logStartup)

	// port
	router.Run(":" + c.Get("server", "port"))
}
