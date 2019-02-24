package main

import (
	"os"

	"github.com/nlopes/slack"
)

type configuration struct {
	SlackToken    string `env:"SLACK_TOKEN"`
	SlackBotToken string `env:"SLACK_BOT_TOKEN"`
}

var config = configuration{
	SlackToken:    os.Getenv("SLACK_TOKEN"),
	SlackBotToken: os.Getenv("SLACK_BOT_TOKEN"),
}

var api = slack.New(config.SlackBotToken)

func main() {

	RunServer()
	/* user, err := api.GetUserByEmail("alessio.izzo86@gmail.com")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email) */
}
