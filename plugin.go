package main

import (
	"fmt"
	// "bytes"
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
	err := p.api.LoadPluginConfiguration(&configuration)
	p.configuration.Store(&configuration)
	return err
}


func (p *Plugin) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	config := p.config()

	if !config.Enabled || config.Secret == "" || config.UserName == "" {
		http.Error(w, "This plugin is not configured.", http.StatusForbidden)
		return
	}

	var webhook Webhook
	json.NewDecoder(r.Body).Decode(&webhook);

	// mm
	team, _ := p.api.GetTeamByName("test");
	channel, _ := p.api.GetChannelByName("test", team.Id);
	user, _ := p.api.GetUserByUsername(config.UserName);

	attachment, err := webhook.SlackAttachment();
	if err != nil {
		return
	}

	if _, err := p.api.CreatePost(&model.Post{
		ChannelId: channel.Id,
		Type:      model.POST_SLACK_ATTACHMENT,
		UserId:    user.Id,
		Props: map[string]interface{}{
			"from_webhook":  "true",
			"use_user_icon": "true",
			"attachments":   []*model.SlackAttachment{attachment},
		},
	}); err != nil {
		http.Error(w, err.Message, err.StatusCode)
	}

	fmt.Fprint(w, "ok.")

}


