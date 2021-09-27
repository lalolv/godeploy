package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/ccding/go-config-reader/config"
	"github.com/gin-gonic/gin"
)

var (
	shPath string
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
	system(fmt.Sprintf(`cd %s && sh %s.sh`, shPath, shell))
	c.String(http.StatusOK, "OK")
}

// 调用系统指令的方法，参数s 就是调用的shell命令
func system(s string) {
	// 调用Command函数
	cmd := exec.Command("/bin/sh", "-c", s)
	// 缓冲字节
	var out bytes.Buffer

	// 标准输出
	cmd.Stdout = &out
	// 运行指令 ，做判断
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	// 输出执行结果
	fmt.Printf("%s", out.String())
}
