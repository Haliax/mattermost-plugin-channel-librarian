# Mattermost New Channel Notify Plugin

This plugin is a fork of [gitlab.com/thepill/mattermost-plugin-newchannelnotify](https://gitlab.com/thepill/mattermost-plugin-newchannelnotify).

A plugin for Mattermost to notify all users about newly created channels.

![screenshot](https://i.imgur.com/SII7ZEi.png)


## Notes and Acknowledgements

This plugin is a fork of [gitlab.com/thepill/mattermost-plugin-newchannelnotify](https://gitlab.com/thepill/mattermost-plugin-newchannelnotify).


## Installation

**This plugin requires Mattermost 6.0 or higher.**

<!-- Download the [latest release here](https://gitlab.com/thepill/mattermost-plugin-newchannelnotify/uploads/cb855f926098701e017c97de403ee3d3/mattermost-plugin-newchannelnotify-0.12.0.tar.gz) (SHA256: `36bbc87c1712fa899c7b89ceed0fed48f7a7682a4af8916a2eca8332fb6f475e`) -->
<!-- In production, deploy and upload your plugin via the [System Console](https://about.mattermost.com/default-plugin-uploads). -->

Optionally, change `settings` under the plugins settings menu in System Console:
- Bot Name
- Channel to Post
- Include private channels
- IncludeChannelPurpose
- TeamsToWatch
- MessageTemplate

**A note about Mentions**:
Mentions will only work properly if the bot account has been assigned to the team and, if needed, the channel where it will post.


## Developing 

See https://github.com/mattermost/mattermost-plugin-starter-template#development

Use `make dist` to build distributions of the plugin that you can upload to a Mattermost server.

Use `make deploy` to deploy the plugin to your local server. Before running `make deploy` you need to set a few environment variables:

```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
```
