package main

func not_main() {

}

/* import (
	"fmt"
	"log"
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

var api = slack.New(config.SlackBotToken,
	slack.OptionDebug(true),
	slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
)

func main() {

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "CGEF2TM7H"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
*/
