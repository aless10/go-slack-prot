package main

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"io/ioutil"
	"net/http"
)

func okJSONHandler(rw http.ResponseWriter, r *http.Request) error {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(InChannelResponse{"Request Received!", "in_channel", []slack.Attachment{}})
	_, err := rw.Write(response)

	return err
}

func koJSONHandler(rw http.ResponseWriter, r *http.Request) error {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(InChannelResponse{fmt.Sprintf("Something went wrong: %d!", http.StatusInternalServerError), "in_channel", []slack.Attachment{}})
	_, err := rw.Write(response)

	return err
}

func UserNotSubscribedHandler(channelID string) {
	msg := fmt.Sprintf("You are not subscribed. Please type /subscribe [your-github-username]")
	msgOption := slack.MsgOptionText(msg, true)
	sendMessage(channelID, msgOption)

}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := InChannelResponse{"pong", "in_channel", []slack.Attachment{}}
	js, err := json.Marshal(response)
	if err != nil {
		Error.Printf(err.Error(), "while marshalling response. Returning", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(js)
	if err != nil {
		Error.Printf(err.Error(), "while writing response. Returning", http.StatusInternalServerError)
		return
	}
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		Error.Printf(err.Error(), "while parsing command. Returning", http.StatusInternalServerError)
		return
	}
	user, exists := ProtSubscribedUsers[command.UserID]
	if !exists {

		newUser := SubscribedUser{
			SlackUserID:    command.UserID,
			SlackChannelId: command.ChannelID,
			GithubUser:     command.Text,}

		ProtSubscribedUsers[newUser.SlackUserID] = newUser
		response := InChannelResponse{"User " + newUser.GithubUser + " added", "in_channel", []slack.Attachment{}}
		js, err := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(js)
		if err != nil {
			Error.Printf(err.Error(), "while writing response. Returning", http.StatusInternalServerError)
			koJSONHandler(w, r)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		msgText := fmt.Sprintf("%s Already subscribed", user.GithubUser)
		response := InChannelResponse{msgText, "in_channel", []slack.Attachment{}}
		js, err := json.Marshal(response)
		_, err = w.Write(js)
		if err != nil {
			Error.Printf(err.Error(), "while writing response. Returning", http.StatusInternalServerError)
			return
		}
	}
}

func githubPrHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var t GithubPResponse
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		Error.Println(err)
	}
	jsonErr := json.Unmarshal([]byte(body), &t)

	if jsonErr != nil {
		Error.Println(jsonErr)
	}
	if *t.PR.State != "closed" {
		option := makeMessage(&t.PR)
		for _, reviewer := range t.PR.Assignees {
			user, err := GetUserGithubName(*reviewer.Login)
			if err != nil {
				Error.Println(err)
			}
			sendMessage(user.SlackChannelId, *option)
		}
	}

}

func commandHandler(w http.ResponseWriter, r *http.Request) {
	// 0. Return the response to slack
	defer r.Body.Close()
	command, err := slack.SlashCommandParse(r)
	if err != nil {
		Error.Printf("Error while parsing the command: %s. Returning status error %d", err.Error(), http.StatusInternalServerError)
		koJSONHandler(w, r)
	} else {
		err := okJSONHandler(w, r)
		user, err := GetUserByID(command.UserID)
		if err != nil {
			Error.Printf("User with ID %s not found", command.UserID)
			UserNotSubscribedHandler(command.ChannelID)

		} else {
			switch command.Command {
			case "/prot":
				go sendResponse(user)
			case "/protlist":
				go sendResponsePrList(user)

			}
		}
	}
}

func listResponseHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	response, err := listResponseParse(r)
	if err != nil {
		Error.Printf(err.Error(), "while parsing command. Returning", http.StatusInternalServerError)
		koJSONHandler(w, r)
	} else {
		okJSONHandler(w, r)
		Info.Println(response)
		go sendSingleRepoResponse(response)

	}

}
