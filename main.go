package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type configuration struct {
	SlackToken    string `env:"SLACK_TOKEN"`
	SlackBotToken string `env:"SLACK_BOT_ACCESS_TOKEN"`
	GithubToken   string `env:"GITHUB_TOKEN"`
	Organization  string `env:"ORGANIZATION"`
}

var Config = configuration{
	SlackToken:    os.Getenv("SLACK_TOKEN"),
	SlackBotToken: os.Getenv("SLACK_BOT_ACCESS_TOKEN"),
	GithubToken:   os.Getenv("GITHUB_TOKEN"),
	Organization:  os.Getenv("ORGANIZATION"),
}

var Api = slack.New(Config.SlackBotToken)

func createGithubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Config.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

var githubClient = createGithubClient()

type GithubResponse struct {
	PullRequestList []*github.PullRequest
	SlackUser       *SubscribedUser
}

var ProtSubscribedUsers = make(map[string]SubscribedUser)

type SubscribedUser struct {
	SlackUserID    string
	SlackChannelId string
	GithubUser     string
}

func okJSONHandler(rw http.ResponseWriter, r *http.Request) error {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(InChannelResponse{"Request Received!", "in_channel"})
	_, err := rw.Write(response)

	return err
}

func sendResponse(command slack.SlashCommand) {
	user, err := getUserByID(command.UserID)
	if err != nil {
		log.Fatalf("User with ID %s not found", command.UserID)
	}
	response := GithubResponse{SlackUser: &user}

	// 1. For all the repos
	repos, _, err := githubClient.Repositories.List(context.Background(), "", nil)
	//repos, _, err := githubClient.Repositories.ListByOrg(context.Background(), Config.Organization, nil)
	for _, repo := range repos {
		owner := repo.Owner.Login
		repoName := repo.Name
		pullsList, _, err := githubClient.PullRequests.List(context.Background(), *owner, *repoName, &github.PullRequestListOptions{State: "open"})
		if err != nil {
			log.Println(err)
			return
		}
		for _, pull := range pullsList {
			// TODO remove this assignment, range over pull.RequestedReviewers
			for _, assignee := range pull.Assignees {
				log.Println(assignee)
				if err != nil {
					log.Println(err)
					return
				}
				if *assignee.Login == user.GithubUser {
					log.Println("Yes")
					response.PullRequestList = append(response.PullRequestList, pull)
				}
			}
		}

	}
	msgOptionText := fmt.Sprintf("Here are your PR")
	msgOptionTitle := slack.MsgOptionText(msgOptionText, true)
	sendMessage(user.SlackChannelId, msgOptionTitle, *makeMessage(response.PullRequestList ...))

}

func sendResponseHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HEre we are baby")
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Println(err)
	}
	log.Println(command)
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	// 0. Return the response to slack
	defer r.Body.Close()
	err := okJSONHandler(w, r)
	command, err := slack.SlashCommandParse(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	go sendResponse(command)

}

type InChannelResponse struct {
	Text         string `json:"text"`
	ResponseType string `json:"response_type"`
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, exists := ProtSubscribedUsers[command.UserID]
	if !exists {

		newUser := SubscribedUser{
			SlackUserID:    command.UserID,
			SlackChannelId: command.ChannelID,
			GithubUser:     command.Text,}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ProtSubscribedUsers[newUser.SlackUserID] = newUser
		response := InChannelResponse{"User added", "in_channel"}
		js, err := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(js)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		msgText := fmt.Sprintf("%s Already subscribed", user.GithubUser)
		response := InChannelResponse{msgText, "in_channel"}
		js, err := json.Marshal(response)
		_, err = w.Write(js)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := InChannelResponse{"pong", "in_channel"}
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
			// Todo fix this assignment
			pull.Labels = append(pull.Labels, &github.Label{Name: github.String("No Label Added"), Color: github.String("green"),})
		}
		labelInfo.Name = pull.Labels[0].Name
		labelInfo.Color = pull.Labels[0].Color
		msgOptionText := fmt.Sprintf("%s requests your review on this PR", *pull.User.Login)
		msgAttTitle := fmt.Sprintf("%s -> %s", *pull.Title, *pull.Base.Repo.Name)
		msgAttText := fmt.Sprintf("%s\n Label: %s", *pull.HTMLURL, *labelInfo.Name)
		log.Println(pull)
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
		log.Println(err)

	}
	log.Println(respChannel)

}

func githubPrHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var t GithubPResponse
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	jsonErr := json.Unmarshal([]byte(body), &t)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	log.Println(t.Action)
	log.Println(t.PR.State != github.String("closed"))
	if t.PR.State != github.String("closed") {
		option := makeMessage(&t.PR)
		log.Println(option)
		for _, reviewer := range t.PR.Assignees {
			user, err := getUserGithubName(*reviewer.Login)
			if err != nil {
				log.Println(err)
			}
			sendMessage(user.SlackChannelId, *option)
		}
	}

}

func getUserByID(id string) (SubscribedUser, error) {
	return ProtSubscribedUsers[id], nil

}

func getUserGithubName(login string) (SubscribedUser, error) {

	for _, subscribedUser := range ProtSubscribedUsers {
		user, err := getUserByID(subscribedUser.SlackUserID)
		if err != nil {
			log.Println(err)
		}
		if user.GithubUser == login {
			return user, nil
		}
	}
	nullUser := SubscribedUser{}
	return nullUser, nil
}

func registerEndpoints() {
	http.HandleFunc("/message", commandHandler)
	http.HandleFunc("/send-response", sendResponseHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/github-pr", githubPrHandler)
	http.HandleFunc("/", pingHandler)
}

// RunServer This is the server runner
func RunServer() {
	registerEndpoints()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Print("Running The server\n")
	RunServer()
}
