package main

import (
    //  "fmt"
    //  "bytes"
    // "crypto/subtle"
    "encoding/json"
	"net/http"
    "sync/atomic"

    "github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type Configuration struct {
 	Enabled  bool
 	Secret   string
 	UserName string
}

type Plugin struct {
    plugin.MattermostPlugin

    api           plugin.API
    configuration atomic.Value
}

func (p *Plugin) OnActivate(api plugin.API) error {
	p.api = api
	return p.OnConfigurationChange()
}

func (p *Plugin) config() *Configuration {
	return p.configuration.Load().(*Configuration)
}

func (p *Plugin) OnConfigurationChange() error {
	var configuration Configuration
	err := p.API.LoadPluginConfiguration(&configuration)
	p.configuration.Store(&configuration)
	return err
}


func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	config := p.config()

	if !config.Enabled || config.Secret == "" || config.UserName == "" {
		http.Error(w, "This plugin is not configured.", http.StatusForbidden)
		return
	} else if r.URL.Path != "/webhook" {
		http.NotFound(w, r)
		return
	} else if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var webhook Webhook

	// mm
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
    } else if r.URL.Query().Get("channel") == "" || r.URL.Query().Get("team") == "" {
        p.API.LogInfo("log1")
        http.Error(w, "you must provide a team and a chanel name.", http.StatusBadRequest)
    } else if team, err := p.API.GetTeamByName(r.URL.Query().Get("team")); err != nil {
        p.API.LogInfo(err.Message)
        http.Error(w, err.Message, err.StatusCode)
    } else if channel, err := p.API.GetChannelByName(team.Id, r.URL.Query().Get("channel"), false); err != nil {
        p.API.LogInfo(err.Message)
        http.Error(w, err.Message, err.StatusCode)
    } else if user, err := p.API.GetUserByUsername(config.UserName); err != nil {
        p.API.LogInfo(err.Message)
        http.Error(w, err.Message, err.StatusCode)
    } else if attachment, err := webhook.SlackAttachment(); err != nil {
        p.API.LogInfo("log5")
        http.Error(w, "", http.StatusBadRequest)
	} else if _, err := p.API.CreatePost(&model.Post{

		ChannelId: channel.Id,
		Type:      model.POST_SLACK_ATTACHMENT,
		UserId:    user.Id,
		Props: map[string]interface{}{
			"from_webhook":  "true",
			"attachments":   []*model.SlackAttachment{attachment},
			"override_username": "Taiga.io",
			"override_icon_url": "https://avatars0.githubusercontent.com/u/6905422?s=200&v=4",
		},

	}); err != nil {
		http.Error(w, err.Message, err.StatusCode)
	}
}


