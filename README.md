# Mattermost New Channel Notify Plugin

A plugin for Mattermost to notify all users about newly created channels.

![screenshot](https://i.imgur.com/SII7ZEi.png)

## Installation

__Requires Mattermost 6.0 or higher.__

Download the [latest release here](https://gitlab.com/thepill/mattermost-plugin-newchannelnotify/uploads/82bbcab1589d1997d5aca02ee8fcba2c/mattermost-plugin-newchannelnotify-0.11.0.tar.gz) (SHA256: `8e781f1375453314802f478947ce3f3a32c43e121d79a3082304a35de70732f2`)

In production, deploy and upload your plugin via the [System Console](https://about.mattermost.com/default-plugin-uploads).

Optionally, change `settings` under the plugins settings menu in System Console:
- Bot Name
- Channel to Post
- Include private channels
- Mention (see Note below)
- IncludeChannelPurpose
- TeamsToWatch

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
