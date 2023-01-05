package main

import (
	"fmt"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
	"strings"
)

const defaultBotName = "channel-librarian"

func (p *NewChannelNotifyPlugin) EnsureDefaultValues() {
	config := p.getConfiguration()

	if config.BotUserName == "" {
		config.BotUserName = defaultBotName
	}

	if config.ChannelToPost == "" {
		config.ChannelToPost = model.DefaultChannelName
	}

	if config.MessageTemplate == "" {
		config.MessageTemplate = "Hello there :wave:. You might want to checkout the new channel channel.link created by channel.creator"
	}

	// Ensure all team
	if config.TeamsToWatch != "" {
		teamsToWatch := strings.Split(config.TeamsToWatch, ";")
		if teamsToWatch != nil {
			p.TeamsToWatch = teamsToWatch
		}
	}

	if config.IgnoredPatterns != "" {
		ignoredPatterns := strings.Split(config.IgnoredPatterns, ";")
		if ignoredPatterns != nil {
			p.IgnoredPatterns = ignoredPatterns
		}
	}

	if config.BlacklistedPurposePatterns != "" {
		blacklistedPurposePatterns := strings.Split(config.BlacklistedPurposePatterns, ";")
		if blacklistedPurposePatterns != nil {
			p.BlacklistedPurposePatterns = blacklistedPurposePatterns
		}
	}
}

func (p *NewChannelNotifyPlugin) EnsureBotExists() (string, error) {
	config := p.getConfiguration()

	// Check whether the bot exists.
	existingBot, _ := p.API.GetUserByUsername(config.BotUserName)

	// Otherwise create it.
	if existingBot == nil {
		p.API.LogInfo("ChannelLibrarian: Specified bot user does not exist. Creating...")

		bot, err := p.API.CreateBot(&model.Bot{
			Username:    config.BotUserName,
			DisplayName: "Channel Librarian",
			Description: "Created by the Channel Librarian plugin.",
		})
		if err != nil {
			p.API.LogError(err.Message)
			return "", errors.Wrap(err, "ChannelLibrarian: Failed to ensure the bot.")
		}

		return bot.UserId, nil
	}

	return existingBot.Id, nil
}

func (p *NewChannelNotifyPlugin) IsTeamWatched(channel *model.Channel) bool {
	// Watch all teams by default.
	if p.TeamsToWatch == nil || len(p.TeamsToWatch) <= 0 {
		return true
	}

	team, err := p.API.GetTeam(channel.TeamId)
	if err != nil {
		p.API.LogError(fmt.Sprintf("ChannelLibrarian: Cannot fetch associated team of message: %s", err.Message))
		return false
	}

	if !ContainsValueCaseInsensitive(p.TeamsToWatch, team.Name) {
		return false
	}

	return true
}

func (p *NewChannelNotifyPlugin) IsChannelIgnored(channel *model.Channel) bool {
	if p.IgnoredPatterns == nil || len(p.IgnoredPatterns) <= 0 {
		return false
	}

	if ContainsStringCaseInsensitive(p.IgnoredPatterns, channel.Name) {
		p.API.LogDebug(fmt.Sprintf("ChannelLibrarian: Refusing to announce the channel as its name contains ignored patterns: %s", channel.Name))
		return true
	}

	return false
}

func (p *NewChannelNotifyPlugin) HasBlacklistedPurposePatterns(channel *model.Channel) bool {
	if p.BlacklistedPurposePatterns == nil || len(p.BlacklistedPurposePatterns) <= 0 {
		return false
	}

	if ContainsStringCaseInsensitive(p.BlacklistedPurposePatterns, channel.Purpose) {
		p.API.LogDebug("ChannelLibrarian: Refusing to announce the channel as its purpose contains ignored patterns.")
		return true
	}

	return false
}

func (p *NewChannelNotifyPlugin) postMessage(channelId string, message string) error {
	bot, err := p.API.GetUser(p.BotUserId)
	if err != nil {
		p.API.LogError(err.Message)
		return errors.Wrap(err, "ChannelLibrarian: Bot user could not be found.")
	}

	post, err := p.API.CreatePost(&model.Post{
		ChannelId: channelId,
		UserId:    bot.Id,
		Message:   message,
	})

	if err != nil {
		p.API.LogError(err.Message)
		return errors.Wrap(err, "ChannelLibrarian: Message could not be posted.")
	}

	p.API.LogDebug(fmt.Sprintf("ChannelLibrarian: Created post %s", post.Id))
	return nil
}
