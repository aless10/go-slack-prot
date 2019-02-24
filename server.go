package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	fmt.Println(r.GetBody())
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&message)

	if err != nil {
		panic(err)
	}

	/* 	buf := new(bytes.Buffer)
	   	buf.ReadFrom(r.Body)
	   	body := buf.String()
	*/
	fmt.Println("Request body is \n", message)
	user, err := api.GetUserByEmail(message.text)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
}

func registerEndpoints() {
	//http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/message", messageHandler)
}

// RunServer This is the server runner
func RunServer() {
	registerEndpoints()
	http.ListenAndServe(":8000", nil)
}
