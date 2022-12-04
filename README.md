# Mattermost New Channel Notify Plugin

A plugin for Mattermost to notify all users about newly created channels.

![screenshot](https://i.imgur.com/SII7ZEi.png)

## Installation

__Requires Mattermost 6.0 or higher.__

Download the [latest release here](https://gitlab.com/thepill/mattermost-plugin-newchannelnotify/uploads/19c494d648d735746698c6cd73f71c2b/mattermost-plugin-newchannelnotify-0.9.3.tar.gz) (SHA256: `ca281d4751c32d415e7a6ba6c69e47fa5f3efa4619eb72bc0772b76465a55a05`)

In production, deploy and upload your plugin via the [System Console](https://about.mattermost.com/default-plugin-uploads).

Optionally, change `settings` under the plugins settings menu in System Console:
- Bot Name
- Channel to Post
- Include private channels
- Mention
- IncludeChannelPurpose

## Developing 

See https://github.com/mattermost/mattermost-plugin-starter-template#development

Use `make dist` to build distributions of the plugin that you can upload to a Mattermost server.

Use `make deploy` to deploy the plugin to your local server. Before running `make deploy` you need to set a few environment variables:

```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
```
