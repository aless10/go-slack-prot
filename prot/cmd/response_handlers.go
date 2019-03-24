package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"log"
	"net/http"
)

type GithubResponse struct {
	PullRequestList []*github.PullRequest
	SlackUser       *SubscribedUser
}

type lightInteractionCallback struct {
	SelectedRepo string
	User         slack.User
}

var githubClient = createGithubClient()

func sendResponsePrList(user SubscribedUser) {
	// 1. For all the repos
	repos, _, _ := githubClient.Repositories.List(context.Background(), "", nil)
	//repos, _, err := githubClient.Repositories.ListByOrg(context.Background(), Config.Organization, nil)

	msgText := fmt.Sprintf("Here are your Repos")
	msgPreText := fmt.Sprintf("Select the one you want to review")
	attAction := slack.AttachmentAction{Name: fmt.Sprintf("GoGet"),
		Text: fmt.Sprintf("Select a repo"),
		Type: fmt.Sprintf("select"),
	}
	for _, repo := range repos {
		attAction.Options = append(attAction.Options,
			slack.AttachmentActionOption{Text: *repo.FullName,
				Value: *repo.Name})
	}
	msgAtt := slack.Attachment{
		CallbackID: "GoGet",
		Text:       msgText,
		Pretext:    msgPreText,
		Color:      "#3AA3E3",
		Actions:    []slack.AttachmentAction{attAction},
	}
	sendMessage(user.SlackChannelId, slack.MsgOptionAttachments(msgAtt))

}

func sendResponse(user SubscribedUser) {
	response := GithubResponse{SlackUser: &user}

	// 1. For all the repos
	repos, _, _ := githubClient.Repositories.List(context.Background(), "", nil)
	//repos, _, err := githubClient.Repositories.ListByOrg(context.Background(), Config.Organization, nil)
	for _, repo := range repos {
		pullsList, _, err := githubClient.PullRequests.List(context.Background(), *repo.Owner.Login, *repo.Name, &github.PullRequestListOptions{State: "open"})
		if err != nil {
			Error.Println(err)
			return
		}
		for _, pull := range pullsList {
			// TODO remove this assignment, range over pull.RequestedReviewers
			for _, assignee := range pull.Assignees {
				Info.Println(assignee)
				if err != nil {
					Error.Println(err)
					return
				}
				if *assignee.Login == user.GithubUser {
					response.PullRequestList = append(response.PullRequestList, pull)
				}
			}
		}

	}
	var msgOptionText string
	if len(response.PullRequestList) > 0 {
		msgOptionText = fmt.Sprintf("Here are your PR!")
	} else {
		msgOptionText = fmt.Sprintf("You have no PR to review!")
	}
	msgOptionTitle := slack.MsgOptionText(msgOptionText, true)
	sendMessage(user.SlackChannelId, msgOptionTitle, *makeMessage(response.PullRequestList ...))

}

type InChannelResponse struct {
	Text         string             `json:"text"`
	ResponseType string             `json:"response_type"`
	Attachments  []slack.Attachment `json:"attachments"`
}

type GithubPResponse struct {
	Action string             `json:"action"`
	PR     github.PullRequest `json:"pull_request"`
}

func makeMessage(pr ...*github.PullRequest) *slack.MsgOption {
	var atts []slack.Attachment
	for _, pull := range pr {
		labelInfo := github.Label{}
		if len(pull.Labels) == 0 {
			pull.Labels = append(pull.Labels, &github.Label{Name: github.String("No Label Added"), Color: github.String("green"),})
		}
		labelInfo.Name = pull.Labels[0].Name
		labelInfo.Color = pull.Labels[0].Color
		msgOptionText := fmt.Sprintf("%s requests your review on this PR", *pull.User.Login)
		msgAttTitle := fmt.Sprintf("%s -> %s", *pull.Title, *pull.Base.Repo.Name)
		msgAttText := fmt.Sprintf("%s\n Label: %s", *pull.HTMLURL, *labelInfo.Name)
		Info.Println("PR: ", pull)
		atts = append(atts, slack.Attachment{Title: msgAttTitle,
			Text:       msgAttText,
			Color:      *labelInfo.Color,
			AuthorName: msgOptionText,
		})

	}

	msgAttachments := slack.MsgOptionAttachments(atts...)
	return &msgAttachments
}

func sendMessage(channelId string, options ...slack.MsgOption) {
	respChannel, _, err := Api.PostMessage(channelId, options...)
	if err != nil {
		Error.Println(err)

	}
	Info.Println("Message sent to ", respChannel)

}

// SlashCommandParse will parse the request of the slash command
func listResponseParse(r *http.Request) (lightInteractionCallback, error) {
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var rBody slack.InteractionCallback
	var s lightInteractionCallback
	if err := r.ParseForm(); err != nil {
		return s, err
	}
	err := json.Unmarshal([]byte(r.PostForm["payload"][0]), &rBody)
	if err != nil {
		return s, err
	}

	s.User = rBody.User
	s.SelectedRepo = rBody.ActionCallback.Actions[0].SelectedOptions[0].Value

	return s, nil
}

func sendSingleRepoResponse(r lightInteractionCallback) {
	user, err := GetUserByID(r.User.ID)
	if err != nil {
		Error.Printf("User with ID %s not found", r.User.ID)
	}
	response := GithubResponse{SlackUser: &user}
	pullsList, _, err := githubClient.PullRequests.List(context.Background(), "aless10", r.SelectedRepo, &github.PullRequestListOptions{State: "open"})
	if err != nil {
		log.Println(err)
		return
	}

	msgOptionText := fmt.Sprintf("Here are your PR of Repo %s", r.SelectedRepo)

	if len(pullsList) == 0 {
		msgOptionText = fmt.Sprintf("No PR to review for Repo %s", r.SelectedRepo)
	}

	for _, pull := range pullsList {
		// TODO remove this assignment, range over pull.RequestedReviewers
		for _, assignee := range pull.Assignees {
			Info.Println(assignee)
			if err != nil {
				log.Println(err)
				return
			}
			if *assignee.Login == user.GithubUser {
				response.PullRequestList = append(response.PullRequestList, pull)
			}
		}
	}

	msgOptionTitle := slack.MsgOptionText(msgOptionText, true)
	sendMessage(user.SlackChannelId, msgOptionTitle, *makeMessage(response.PullRequestList ...))

}
