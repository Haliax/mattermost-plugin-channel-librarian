package main

import (
	"github.com/mattermost/mattermost-server/v6/model"
)

func (p *NewChannelNotifyPlugin) ensureBotExists() {
	config := p.getConfiguration()

	existingBot, _ := p.API.GetUserByUsername(config.BotUserName)

	// Ensure the bot exists
	if existingBot == nil {
		p.API.LogInfo("Bot user doesnt exist -> creating")
		_, err := p.API.CreateBot(&model.Bot{
			Username:    config.BotUserName,
			DisplayName: "NewChannelBot",
		})

		if err != nil {
			p.API.LogError(err.Message)
		}
	}
}
