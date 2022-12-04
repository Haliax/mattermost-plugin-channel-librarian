# Mattermost New Channel Notify Plugin

A plugin for Mattermost to notify all users about newly created channels.

![screenshot](https://i.imgur.com/SII7ZEi.png)

## Installation

__Requires Mattermost 6.0 or higher.__

Download the [latest release here](https://gitlab.com/thepill/mattermost-plugin-newchannelnotify/uploads/260cc9f92af33bb923cd405ec758ad0a/mattermost-plugin-newchannelnotify-0.10.0.tar.gz) (SHA256: `87555a032ada2ee258e40515007937dd01500c3f4556de34cbeb367ef8db30db`)

In production, deploy and upload your plugin via the [System Console](https://about.mattermost.com/default-plugin-uploads).

Optionally, change `settings` under the plugins settings menu in System Console:
- Bot Name
- Channel to Post
- Include private channels
- Mention
- IncludeChannelPurpose
- TeamsToWatch

## Developing 

See https://github.com/mattermost/mattermost-plugin-starter-template#development

Use `make dist` to build distributions of the plugin that you can upload to a Mattermost server.

Use `make deploy` to deploy the plugin to your local server. Before running `make deploy` you need to set a few environment variables:

```
export MM_SERVICESETTINGS_SITEURL=http://localhost:8065
export MM_ADMIN_USERNAME=admin
export MM_ADMIN_PASSWORD=password
```
