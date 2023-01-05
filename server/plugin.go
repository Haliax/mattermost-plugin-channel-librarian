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

	BotUserId                  string
	TeamsToWatch               []string
	IgnoredPatterns            []string
	BlacklistedPurposePatterns []string
}

func (p *NewChannelNotifyPlugin) OnActivate() error {
	p.API.LogDebug("Plugin loading...")

	// Ensure default values.
	p.EnsureDefaultValues()

	// Ensure the bot.
	botId, err := p.EnsureBotExists()
	if err != nil {
		return errors.Wrap(err, "failed to ensure channel librarian bot")
	}
	p.BotUserId = botId

	p.API.LogInfo("Plugin loaded.")
	return nil
}

func (p *NewChannelNotifyPlugin) OnDeactivate() error {
	return nil
}

func (p *NewChannelNotifyPlugin) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
	p.AnnounceNewChannel(c, channel)
}

func (p *NewChannelNotifyPlugin) AnnounceNewChannel(c *plugin.Context, channel *model.Channel) {
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
	if !p.IsTeamWatched(channel) || p.IsChannelIgnored(channel) {
		return
	}

	// Ignore private channels depending on the configuration.
	config := p.getConfiguration()
	isPrivateChannel := false
	if channel.Type == model.ChannelTypePrivate {
		if config.IncludePrivateChannels == false {
			return
		}

		if p.HasBlacklistedPurposePatterns(channel) {
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
		purposeText = channel.Purpose
	}

	privateText := "public"
	if isPrivateChannel {
		privateText = "private"
	}

	message := FormatTemplate(config.MessageTemplate, creator.Username, channel.DisplayName, channel.Name, purposeText, privateText)
	_ = p.postMessage(channelToPostTo.Id, message)
}

func FormatTemplate(template string, creatorName string, channelDisplayName string, channelName string, channelPurpose string, channelType string) string {
	message := template

	message = strings.Replace(message, "channel.creator", "@"+creatorName, -1)
	message = strings.Replace(message, "channel.name", channelDisplayName, -1)
	message = strings.Replace(message, "channel.link", "~"+channelName, -1)
	message = strings.Replace(message, "channel.purpose", channelPurpose, -1)
	message = strings.Replace(message, "channel.type", channelType, -1)

	return message
}
