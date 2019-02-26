package main

import (
	"fmt"
	"os"

	"github.com/nlopes/slack"
)

type configuration struct {
	SlackToken    string `env:"SLACK_TOKEN"`
	SlackBotToken string `env:"SLACK_BOT_TOKEN"`
}

var Config = configuration{
	SlackToken:    os.Getenv("SLACK_TOKEN"),
	SlackBotToken: os.Getenv("SLACK_BOT_TOKEN"),
}

var api = slack.New(Config.SlackBotToken)

func main() {
	fmt.Print("Running The server\n")
	RunServer()
}
