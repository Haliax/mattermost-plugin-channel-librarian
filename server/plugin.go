package main

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"sync"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

type NewChannelNotifyPlugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	BotUserId    string
	TeamsToWatch []string
}

func (p *NewChannelNotifyPlugin) OnActivate() error {
	p.API.LogDebug("Plugin loading...")
	config := p.getConfiguration()

	// Ensure default values.
	p.ensureDefaultValues()

	// Ensure the bot.
	botId, err := p.ensureBotExists()
	if err != nil {
		return errors.Wrap(err, "failed to ensure channel librarian bot")
	}
	p.BotUserId = botId

	// Ensure all team
	if config.TeamsToWatch != "" {
		teamsToWatch := strings.Split(config.TeamsToWatch, ";")
		if teamsToWatch != nil {
			p.TeamsToWatch = teamsToWatch
		}
	}

	p.API.LogInfo("Plugin loaded.")
	return nil
}

func (p *NewChannelNotifyPlugin) OnDeactivate() error {
	return nil
}

func (p *NewChannelNotifyPlugin) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
	p.announceNewChannel(c, channel)
}

func (p *NewChannelNotifyPlugin) announceNewChannel(c *plugin.Context, channel *model.Channel) {
	// Ignore DMs.
	if channel.Type == model.ChannelTypeDirect || channel.Type == model.ChannelTypeGroup {
		return
	}

	p.API.LogDebug(fmt.Sprintf("ChannelLibrarian: New channel with id [%s], type [%s] created", channel.Id, channel.Type))

	// Ignore the channel if it is created automatically.
	if channel.CreatorId == "" {
		p.API.LogDebug("ChannelLibrarian: Ignored channel due to not having a valid creator.")
		return
	}

	// Ignore the channel if the team is not watched.
	if !p.isTeamWatched(channel) {
		return
	}

	// Ignore private channels depending on the configuration.
	config := p.getConfiguration()
	isPrivateChannel := false
	if channel.Type == model.ChannelTypePrivate {
		if config.IncludePrivateChannels == false {
			return
		}

		isPrivateChannel = true
	}

	channelToPostTo, err := p.API.GetChannelByName(channel.TeamId, config.ChannelToPost, false)
	if err != nil {
		p.API.LogError(fmt.Sprintf("ChannelLibrarian: Could not find channel to post to: %s", err.Message))
		return
	}

	creator, err := p.API.GetUser(channel.CreatorId)
	if err != nil {
		p.API.LogError(fmt.Sprintf("ChannelLibrarian: Could not find the creator of the channel: %s", err.Message))
		return
	}

	purposeText := ""
	if config.IncludeChannelPurpose && channel.Purpose != "" {
		purposeText = "\n\n**Purpose:** " + channel.Purpose
	}

	privateText := ""
	if isPrivateChannel {
		privateText = " **[private]**"
	}

	message := fmt.Sprintf(
		"%sHello there :wave:. You might want to check out the new%s channel ~%s created by @%s %s",
		config.Mention, privateText, channel.Name, creator.Username, purposeText,
	)
	_ = p.postMessage(channelToPostTo.Id, message)
}
