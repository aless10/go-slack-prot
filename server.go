package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nlopes/slack/slackevents"
)

type Message struct {
	channel_id   string
	channel_name string
	command      string
	response_url string
	team_domain  string
	team_id      string
	text         string
	token        string
	trigger_id   string
	user_id      string
	user_name    string
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	message := Message{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)

	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: Config.SlackToken}))
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	fmt.Println("Events Api: ", eventsAPIEvent)
	fmt.Println("Request body is \n", body)
	fmt.Println("Request body message is \n", message)
	body_json, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Println("Json marshalled: ", body_json)
	user, err := api.GetUserByEmail("message.text")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
}

type InChannelResponse struct {
	text          string
	response_type string
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	response := InChannelResponse{"pong", "in_channel"}
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type Profile struct {
	Text          string
	response_type string
	Pippo         string
}

func foo(w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Alex", "in_channel", "Pluto"}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func registerEndpoints() {
	http.HandleFunc("/message", messageHandler)
	http.HandleFunc("/", foo)
	http.HandleFunc("/foo", foo)
}

// RunServer This is the server runner
func RunServer() {
	registerEndpoints()
	http.ListenAndServe(":8000", nil)
}
