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

var LogFile = os.Getenv("APP_LOG_FILE")
var Api = slack.New(Config.SlackBotToken)

var ProtSubscribedUsers = make(map[string]SubscribedUser)

type SubscribedUser struct {
	SlackUserID    string
	SlackUserName  string
	SlackChannelId string
	GithubUser     string
}

func main() {
	log.Println(LogFile)
	logFile, err := os.OpenFile(LogFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	}
	initLogs(logFile)
	defer logFile.Close()
	log.Printf("Running The server on %s:%s\n", ServerConfig.Host, ServerConfig.Port)
	serverErr := RunServer()
	if serverErr != nil {
		log.Fatalf("Error while running the server %s", err)
	}
}
