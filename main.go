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

func commandHandler(w http.ResponseWriter, r *http.Request) {
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := Api.GetUserInfo(command.UserID)
	if err != nil {
		log.Fatalf("User with ID %s not found", command.UserID)
	}

	// 1. For all the repos
	repos, _, err := githubClient.Repositories.List(context.Background(), "", nil)
	for _, repo := range repos {
		owner := repo.Owner.Login
		repoName := repo.Name
		pullsList, _, err := githubClient.PullRequests.List(context.Background(), *owner, *repoName, &github.PullRequestListOptions{State: "open"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, pull := range pullsList {
			// 3. if user.Profile.Email in requested_reviewers email -> append PR info to a list
			for _, assignee := range pull.Assignees {
				log.Println(assignee)
				uResp, _, err := githubClient.Users.Get(context.Background(), *assignee.Login)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				log.Println(uResp)
				if *assignee.Email == user.Profile.Email {
					log.Println("Yes")
				}

			}

		}
	}

	// 4. Return the list to slack

	response := InChannelResponse{user.Profile.FirstName, "in_channel"}
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

type InChannelResponse struct {
	Text         string `json:"text"`
	ResponseType string `json:"response_type"`
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

func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := InChannelResponse{"HELLO CHARLY", "in_channel"}
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

func githubPrHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var t github.PullRequest
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	jsonErr := json.Unmarshal([]byte(body), &t)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	log.Println(t)
	for _, reviewer := range t.Assignees {
		log.Println(reviewer)
	}
	// To delete when reviewers are available
	assignee := t.Assignee
	log.Println(assignee)
	user, err := getUserByEmail(*assignee.Email)
	log.Println(user.Profile.Email)
	if err != nil {
		log.Println(err)
	}
	//w.Header().Set("Content-Type", "application/json")
	//response := InChannelResponse{t.Body, "in_channel"}
	//js, err := json.Marshal(response)
	//w.Write(js)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

}

func getUserByEmail(email string) (*slack.User, error) {
	user, err := Api.GetUserByEmail(email)
	if err != nil {
		log.Println(err)
		return &slack.User{}, err
	}
	log.Println(user.ID, user.Profile.RealName, user.Profile.Email)
	return user, nil
}

func registerEndpoints() {
	http.HandleFunc("/message", commandHandler)
	http.HandleFunc("/github-pr", githubPrHandler)
	http.HandleFunc("/", pingHandler)
	http.HandleFunc("/hello", helloHandler)
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
