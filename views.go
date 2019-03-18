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
			return
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
		Info.Println(option)
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
	err := okJSONHandler(w, r)
	command, err := slack.SlashCommandParse(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch command.Command {
	case "/prot":
		go sendResponse(command)
	case "/protlist":
		go sendResponsePrList(command)
	}

}

func listResponseHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	err := okJSONHandler(w, r)
	response, err := listResponseParse(r)
	if err != nil {
		Error.Printf(err.Error(), "while parsing command. Returning", http.StatusInternalServerError)
		return
	}

	Info.Println(response)
	go sendSingleRepoResponse(response)
}
