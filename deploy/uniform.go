package deploy

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/lalolv/godeploy/glob"
	"github.com/lalolv/godeploy/utils"
	"github.com/rs/zerolog/log"
)

// 部署
func Uniform(c *gin.Context) {
	// Get shell param
	shell := c.Param("shell")

	if !checkPathFile(glob.ShellPath, shell) {
		c.String(http.StatusBadRequest, "File not exist")
		return
	}

	// Call system command
	err := system(fmt.Sprintf(`cd %s && sh %s.sh`, glob.ShellPath, shell))
	if err != nil {
		// slack
		go utils.SendSlack(err.Error())
		// response
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// slack
	go utils.SendSlack(shell + " deployment is done")

	// response
	c.String(http.StatusOK, "OK")
}

// Check shell path and file exist
func checkPathFile(path, shellName string) bool {
	// Check file exist
	_, err := os.Stat(fmt.Sprintf(`%s/%s.sh`, path, shellName))
	if os.IsNotExist(err) {
		log.Error().Msg("File not exist")
		return false
	}

	return true
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
