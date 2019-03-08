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
}

var Config = configuration{
	SlackToken:    os.Getenv("SLACK_TOKEN"),
	SlackBotToken: os.Getenv("SLACK_BOT_ACCESS_TOKEN"),
	GithubToken:   os.Getenv("GITHUB_TOKEN"),
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

type ProtSubscribedUsersMap map[string]SubscribedUser


func subscribeUser() *ProtSubscribedUsersMap {


}


type SubscribedUser struct {
	SlackUserID    string
	SlackChannelId string
	GithubUser     string
}

func okJSONHandler(rw http.ResponseWriter, r *http.Request) error {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(slack.SlackResponse{
		Ok: true,
	})
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
	for _, repo := range repos {
		owner := repo.Owner.Login
		repoName := repo.Name
		pullsList, _, err := githubClient.PullRequests.List(context.Background(), *owner, *repoName, &github.PullRequestListOptions{State: "open"})
		if err != nil {
			log.Println(err)
			return
		}
		for _, pull := range pullsList {
			// TODO remove this assignement, range over pull.RequestedReviewers
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
	log.Println(response.PullRequestList)
	//msgOptionText := fmt.Sprintf("These are PR that you have to review")
	msgOptions, msgAttachments := makeMessage(response.PullRequestList[0])
	//msgOptions := slack.MsgOptionText(msgOptionText, true)
	//msgAttText := fmt.Sprintf("Repo")
	//msgAttachments := slack.MsgOptionAttachments(slack.Attachment{Title: "Title",
	//	Text:  msgAttText,
	//	Color: "green",
	//})
	sendMessage(user.SlackChannelId, *msgOptions, *msgAttachments)
}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	// 0. Return the response to slack
	err := okJSONHandler(w, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	command, err := slack.SlashCommandParse(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(command)

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
	Action string      `json:"action"`
	PR     github.PullRequest `json:"pull_request"`
}

func makeMessage(pr *github.PullRequest) (*slack.MsgOption, *slack.MsgOption) {
	labelInfo := github.Label{}
	if len(pr.Labels) == 0 {
		// Todo fix this assignment
		*labelInfo.Name = fmt.Sprintf("No Label Added")
		*labelInfo.Color = fmt.Sprintf(`green`)
	} else {
		labelInfo.Name = pr.Labels[0].Name
		labelInfo.Color = pr.Labels[0].Color
	}

	msgOptionText := fmt.Sprintf("%s requests your review on this PR", *pr.User.Login)
	msgOptions := slack.MsgOptionText(msgOptionText, true)
	msgAttText := fmt.Sprintf("%s\n Repo: %s\n Label: %s", *pr.Base.Repo.Name, *pr.HTMLURL, *labelInfo.Name)
	msgAttachments := slack.MsgOptionAttachments(slack.Attachment{Title: *pr.Title,
		Text:  msgAttText,
		Color: *labelInfo.Color,
	})
	return &msgOptions, &msgAttachments
}

func sendMessage(channelId string, options ...slack.MsgOption) {
	respChannel, _, err := Api.PostMessage(channelId, options...)
	if err != nil {
		log.Println(err)
		return
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
	option, attachments := makeMessage(&t.PR)
	log.Println(option)

	for _, reviewer := range t.PR.Assignees {
		user, err := getUserGithubName(*reviewer.Login)
		if err != nil {
			log.Println(err)
		}
		go sendMessage(user.SlackChannelId, *option, *attachments)
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
