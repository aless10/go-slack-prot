package main

import (
	"github.com/nlopes/slack"
	"log"
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

type ServerConfiguration struct {
	Host string `env:"HOST_ADDRESS"`
	Port string `env:"PORT"`
}

var ServerConfig = ServerConfiguration{
	Host: os.Getenv("HOST_ADDRESS"),
	Port: os.Getenv("PORT"),
}

var Api = slack.New(Config.SlackBotToken)

var ProtSubscribedUsers = make(map[string]SubscribedUser)

type SubscribedUser struct {
	SlackUserID    string
	SlackChannelId string
	GithubUser     string
}

func main() {
	log.Printf("Running The server on %s:%s\n", ServerConfig.Host, ServerConfig.Port)
	err := RunServer()
	if err != nil {
		log.Fatal("Error while running the server %T", err)
	}
}
