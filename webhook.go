package main

import (
	// "bytes"
	// "strings"
	// "text/template"

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
		} `json:"diff"`
		Comment string `json:"comment"`
		CommentHTML string `json:"comment_html"`
		DeleteCommentDate interface{
		} `json:"delete_comment_date"`
	} `json:"change"`
}


func (w *Webhook) SlackAttachment() (*model.SlackAttachment, error) {

	return &model.SlackAttachment{
		Fallback: "Fallback",
		Color:    "#95b7d0",
		Pretext:  "This is Pre-text",
		Text:     w.Data.UserStory.Subject,
	}, nil
}
