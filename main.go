package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ccding/go-config-reader/config"
	"github.com/gin-gonic/gin"
	"github.com/lalolv/goapis"
)

var (
	shPath       string
	slackUtil    goapis.SlackUtil
	slackToken   string
	slackChannel string
)

func main() {
	// 读取配置文件
	c := config.NewConfig("app.conf")
	err := c.Read()
	// error handle
	if err != nil {
		fmt.Println(err.Error())
	}

	shPath = c.Get("shell", "path")

	// 设置日志
	logPath, _ := pathExists("log")
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
	slackToken = c.Get("slack", "token")
	slackChannel = c.Get("slack", "channel")
	if slackToken != "" {
		slackUtil = goapis.SlackUtil{Token: slackToken}
		sendSlack(logStartup)
	}

	// Creates a gin router with default middlewares:
	// logger and recovery (crash-free) middlewares
	router := gin.New()
	// routers
	router.POST("/deploy/:shell", deploy)
	// log
	log.Log().Msg(logStartup)

	// port
	router.Run(":" + c.Get("server", "port"))
}

// 部署
func deploy(c *gin.Context) {
	// Get shell param
	shell := c.Param("shell")
	// Check file exist
	_, err := os.Stat(fmt.Sprintf(`%s/%s.sh`, shPath, shell))
	if os.IsNotExist(err) {
		log.Error().Msg("File not exist")
		c.String(http.StatusBadRequest, "File not exist")
		return
	}

	// Call system command
	err = system(fmt.Sprintf(`cd %s && sh %s.sh`, shPath, shell))
	if err != nil {
		// slack
		go sendSlack(err.Error())
		// response
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// slack
	go sendSlack("[Deploy] Success to deploy " + shell)

	// response
	c.String(http.StatusOK, "OK")
}

// Send to slack message
func sendSlack(text string) {
	if slackUtil.Token != "" && slackChannel != "" {
		err := slackUtil.PostMessage(slackChannel, text)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}
}

// 调用系统指令的方法，参数s 就是调用的shell命令
func system(s string) error {
	// 调用Command函数
	cmd := exec.Command("/bin/sh", "-c", s)
	// 缓冲字节
	var out bytes.Buffer

	// 标准输出
	cmd.Stdout = &out
	// 运行指令 ，做判断
	err := cmd.Run()
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}

	// 输出执行结果
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		log.Log().Msg(scanner.Text())
	}

	return nil
}

// pathExists 判断文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			return true, nil
		}
	}
	return false, err
}
