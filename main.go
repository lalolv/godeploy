package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"

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

	// slack
	slackToken = c.Get("slack", "token")
	slackChannel = c.Get("slack", "channel")
	if slackToken != "" {
		slackUtil = goapis.SlackUtil{Token: slackToken}
		text := fmt.Sprintf("[%s] Deploy service is done.", "Startup")
		sendSlack(text)
	}

	// Creates a gin router with default middlewares:
	// logger and recovery (crash-free) middlewares
	router := gin.New()
	router.POST("/deploy/:shell", deploy)

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
			fmt.Println(err.Error())
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
		return err
	}
	// 输出执行结果
	fmt.Printf("%s", out.String())

	return nil
}
