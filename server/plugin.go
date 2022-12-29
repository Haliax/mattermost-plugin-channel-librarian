package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

const defaultBotName = "newchannelbot"

type NewChannelNotifyPlugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

func (p *NewChannelNotifyPlugin) OnActivate() error {
	p.API.LogInfo("Plugin loaded")
	return nil
}

// https://play.golang.org/p/Qg_uv_inCek
// contains checks if a string is present in a slice
func containsCaseInsensitive(s []string, str string) bool {
	for _, v := range s {
		if strings.ToLower(v) == strings.ToLower(str) {
			return true
		}
	}

	return false
}

func (p *NewChannelNotifyPlugin) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
	log := fmt.Sprintf("ChannelHasBeenCreated for channel with id [%s], type [%s] triggerd", channel.Id, channel.Type)
	p.API.LogDebug(log)

	config := p.getConfiguration()

	// Check if only specific teams are beeing watched and notiefied
	if config.TeamsToWatch != "" {
		team, err := p.API.GetTeam(channel.TeamId)
		if err != nil {
			p.API.LogError(err.Message)
		}
		p.API.LogDebug(fmt.Sprintf("team: %s", team.Name))

		teamsToWatch := strings.Split(config.TeamsToWatch, ";")

		if !containsCaseInsensitive(teamsToWatch, team.Name) {
			p.API.LogDebug(fmt.Sprintf("team %s is not watched - skipping", team.Name))
			return
		}
	}

	if config.BotUserName == "" {
		config.BotUserName = defaultBotName
	}

	if config.ChannelToPost == "" {
		config.ChannelToPost = model.DefaultChannelName
	}

	ChannelPurpose := ""
	if config.IncludeChannelPurpose && channel.Purpose != "" {
		ChannelPurpose = "\n **" + channel.Name + "'s Purpose:** " + channel.Purpose
	}

	newChannelName := channel.Name

	if channel.Type == model.ChannelTypeDirect || channel.Type == model.ChannelTypeGroup {
		return
	}

	if channel.Type == model.ChannelTypePrivate {
		if config.IncludePrivateChannels == false {
			return
		}
		newChannelName += " [Private]"
	}

	p.ensureBotExists()
	bot, _ := p.API.GetUserByUsername(config.BotUserName)

	mainChannel, err := p.API.GetChannelByName(channel.TeamId, config.ChannelToPost, false)
	if err != nil {
		p.API.LogError(err.Message)
	}

	creator, err := p.API.GetUser(channel.CreatorId)
	if err != nil {
		p.API.LogError(err.Message)
	}

	post, err := p.API.CreatePost(&model.Post{
		ChannelId: mainChannel.Id,
		UserId:    bot.Id,
		Message:   fmt.Sprintf("%sHello there :wave:. You might want to check out the new channel ~%s created by @%s %s", config.Mention, newChannelName, creator.Username, ChannelPurpose),
	})

	p.API.LogDebug(fmt.Sprintf("Created post %s", post.Id))

	if err != nil {
		p.API.LogError(err.Message)
	}
}
