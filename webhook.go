package main

import (
	"bytes"
	// "strings"
	"text/template"
	"errors"

	"github.com/mattermost/mattermost-server/model"
)

 type Owner struct {
	 ID int `json:"id"`
	 Permalink string `json:"permalink"`
	 Username string `json:"username"`
	 FullName string `json:"full_name"`
	 Photo string `json:"photo"`
	 GravatarID string `json:"gravatar_id"`
 }

type Project struct {
	ID int `json:"id"`
	Permalink string `json:"permalink"`
	Name string `json:"name"`
	LogoBigURL interface{
	} `json:"logo_big_url"`
}

type AssignedTo struct {
	ID int `json:"id"`
	Permalink string `json:"permalink"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Photo string `json:"photo"`
	GravatarID string `json:"gravatar_id"`
}

type Status struct {
    ID int `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	Color string `json:"color"`
	IsClosed bool `json:"is_closed"`
	IsArchived bool `json:"is_archived"`
}

type Milestone struct { Permalink string `json:"permalink"`

	Project Project
	Owner Owner
	ID int `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	EstimatedStart string `json:"estimated_start"`
	EstimatedFinish string `json:"estimated_finish"`
	CreatedDate string `json:"created_date"`
	ModifiedDate string `json:"modified_date"`
	Closed bool `json:"closed"`
	Disponibility float64 `json:"disponibility"`
}

type UserStory struct {

	CustomAttributesValues struct {
		EiusVeroFacere string `json:"eius vero facere"`
	} `json:"custom_attributes_values"`

	Watchers []int `json:"watchers"`
	Permalink string `json:"permalink"`
	Tags []string `json:"tags"`
	ExternalReference interface{} `json:"external_reference"`

	Project Project
	Owner Owner
	AssignedTo AssignedTo
	Status Status

	Milestone Milestone
	Points []struct {
		Role string `json:"role"`
		Name string `json:"name"`
		Value float64 `json:"value"`
	} `json:"points"`


	ID int `json:"id"`
	IsBlocked bool `json:"is_blocked"`
	BlockedNote string `json:"blocked_note"`
	Ref int `json:"ref"`
	IsClosed bool `json:"is_closed"`
	CreatedDate string `json:"created_date"`
	ModifiedDate string `json:"modified_date"`
	FinishDate interface{} `json:"finish_date"`
	Subject string `json:"subject"`
	Description string `json:"description"`
	ClientRequirement bool `json:"client_requirement"`
	TeamRequirement bool `json:"team_requirement"`
	GeneratedFromIssue interface{} `json:"generated_from_issue"`
	TribeGig interface{} `json:"tribe_gig"`
}

type Data struct {
	CustomAttributesValues struct {
	} `json:"custom_attributes_values"`
	Watchers []interface{} `json:"watchers"`
	Permalink string `json:"permalink"`
	Tags []string `json:"tags"`

	Project Project
	Owner Owner
	AssignedTo AssignedTo
	Status Status
	UserStory UserStory
	Milestone Milestone

	ID int `json:"id"`
	IsBlocked bool `json:"is_blocked"`
	BlockedNote string `json:"blocked_note"`
	Ref int `json:"ref"`
	CreatedDate string `json:"created_date"`
	ModifiedDate string `json:"modified_date"`
	FinishedDate interface{} `json:"finished_date"`
	Subject string `json:"subject"`
	UsOrder int `json:"us_order"`
	TaskboardOrder int `json:"taskboard_order"`
	Description string `json:"description"`
	IsIocaine bool `json:"is_iocaine"`
	ExternalReference interface{} `json:"external_reference"`
}

type Webhook struct {
	Action string `json:"action"`
	Type string `json:"type"`
	Date string `json:"date"`

	By struct {
		ID int `json:"id"`
		Permalink string `json:"permalink"`
		Username string `json:"username"`
		FullName string `json:"full_name"`
		Photo string `json:"photo"`
		GravatarID string `json:"gravatar_id"`
	} `json:"by"`

	Data Data

	Change struct {
		Diff struct {
			AssignedTo struct {
				From string `json:"from"`
				To string `json:"to"`
			} `json:"assigned_to"`
			Status struct {
				From string `json:"from"`
				To string `json:"to"`
			}
			UserStory struct {
				From string `json:"from"`
				To string `json:"to"`
			}
			Description struct {
				From string `json:"from"`
				To string `json:"to"`
			}
			Subject struct {
				From string `json:"from"`
				To string `json:"to"`
			}


		} `json:"diff"`
		Comment string `json:"comment"`
		CommentHTML string `json:"comment_html"`
		DeleteCommentDate interface{
		} `json:"delete_comment_date"`
	} `json:"change"`
}

func (w *Webhook) renderTitle() (string, error) {

	var typeja string

	switch w.Type {
	case "epic": typeja = "エピック"
	case "userstory": typeja = "ユーザーストーリー"
	case "task": typeja = "タスク"
	}

	var actja string

	switch w.Action {
	case "create": actja = "作成"
	case "change": actja = "変更"
	case "delete": actja = "削除"
	}

	tplBody := "" +
		"{{.By.FullName}} が {{.Typeja}}:{{.Data.Subject}} を {{.Actionja}} しました"

	tpl, err := template.New("post").Parse(tplBody)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, struct {
		*Webhook
		Typeja string
		Actionja string
	}{
		Webhook: w,
		Typeja: typeja,
		Actionja: actja,
	}); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (w *Webhook) renderText() (string, error) {

	if w.Action == "create" {
		// 作成は、タイトルをお知らせ、中身がある場合は中身も表示する。
		return w.Data.Description, nil

	} else if w.Action == "change" {
		// コメント
		if (w.Change.Comment != "") {
			return "コメント: " +
				w.Change.Comment, nil
		}
		// 本文
		if (w.Change.Diff.Description.To != "" ) {
			return "本文: " +
				w.Change.Diff.Description.To, nil
		}
		// ステータス
		if (w.Change.Diff.Status.To != "" ) {
			return "ステータス: " +
				w.Change.Diff.Status.From + "から" +
				w.Change.Diff.Status.To + "へ", nil
		}

	} else if w.Action == "delete" {
		return w.Data.Subject + "が" + w.By.FullName + "によって削除されました", nil
	}

	return "", errors.New("unknown action type")
}

func (w *Webhook) SlackAttachment() (*model.SlackAttachment, error) {

	text, err := w.renderText()
	if err != nil {
		return nil, err
	}

	title,  err := w.renderTitle()
	if err != nil {
		return nil, err
	}

	return &model.SlackAttachment{
		Color: "#95b7d0",

		Text: text,
		Title: title,
		TitleLink: w.Data.Permalink,

		AuthorName: "Taiga.io",
		AuthorLink: "http://localhost:8089",
		AuthorIcon: "https://avatars0.githubusercontent.com/u/6905422?s=200&v=4",
	}, nil

}
