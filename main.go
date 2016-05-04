package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/ccding/go-config-reader/config"
	"github.com/gin-gonic/gin"
)

var (
	shPath string
	shFile string
)

func main() {
	// 读取配置文件
	c := config.NewConfig("app.conf")
	err := c.Read()
	// error handle
	if err != nil {
		fmt.Println(err.Error())
	}

	shFile = c.Get("shell", "file")
	shPath = c.Get("shell", "path")

	// Creates a gin router with default middlewares:
	// logger and recovery (crash-free) middlewares
	router := gin.New()
	router.POST("/deploy", deploy)

	// port
	router.Run(":" + c.Get("server", "port"))
}

// 部署
func deploy(c *gin.Context) {
	system(fmt.Sprintf(`cd %s && sh %s`, shPath, shFile))
	c.String(http.StatusOK, "OK")
}

// 调用系统指令的方法，参数s 就是调用的shell命令
func system(s string) {
	cmd := exec.Command("/bin/sh", "-c", s) //调用Command函数
	var out bytes.Buffer                    //缓冲字节

	cmd.Stdout = &out //标准输出
	err := cmd.Run()  //运行指令 ，做判断
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", out.String()) //输出执行结果
}
