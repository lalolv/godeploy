package utils

import (
	"github.com/lalolv/godeploy/glob"
	"github.com/rs/zerolog/log"
)

// Send to slack message
func SendSlack(text string) {
	if glob.SlackUtil.Token != "" && glob.SlackChannel != "" {
		err := glob.SlackUtil.PostMessage(glob.SlackChannel, text)
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}
}
