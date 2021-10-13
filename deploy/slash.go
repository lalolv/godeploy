package deploy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lalolv/godeploy/glob"
	"github.com/rs/zerolog/log"
)

func SlashCommands(c *gin.Context) {
	token := c.PostForm("token")
	// channelID := c.PostForm("channel_id")
	// channelName := c.PostForm("channel_name")
	// userID := c.PostForm("user_id")
	// userName := c.PostForm("user_name")
	command := c.PostForm("command")
	text := c.PostForm("text")
	responseURL := c.PostForm("response_url")
	// apiAppID := c.PostForm("api_app_id")

	if token != glob.SlashToken {
		c.JSON(http.StatusOK, "Invalid token")
		return
	}
	if command != glob.SlashCommand {
		c.JSON(http.StatusOK, "Invalid command")
		return
	}

	go deploySlash(responseURL, text)

	c.JSON(http.StatusOK, gin.H{
		"blocks": []interface{}{
			gin.H{
				"type": "section",
				"text": gin.H{
					"type": "mrkdwn",
					"text": "Send instruction executed, inform you when complete.",
				},
			},
		},
	})
}

func deploySlash(responseURL, shellName string) {
	if !checkPathFile(glob.ShellPath, shellName) {
		replySlashCommand(responseURL, "File not exist")
		return
	}

	// Call system command
	err := system(fmt.Sprintf(`cd %s && sh %s.sh`, glob.ShellPath, shellName))
	if err != nil {
		// response
		replySlashCommand(responseURL, err.Error())
		return
	}

	replySlashCommand(responseURL, shellName+" deployment is done")
}

func replySlashCommand(responseURL, text string) {
	// 回复消息
	buf, _ := json.Marshal(gin.H{
		"text": text,
	})
	resp, err := http.Post(responseURL, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Error().Msg(err.Error())
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	log.Info().Msg(string(respBody))
}
