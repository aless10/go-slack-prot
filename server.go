package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type User struct {
	Login string `json:"login"`
	Email string
}

type Label struct {
	LabelName string `json:"name"`
}

type Message struct {
	//ChannelID   string `json:"channel_id"`
	//ChannelName string `json:"channel_name"`
	//Command     string `json:"command"`
	Text string `json:"text"`
	//UserID      string `json:"user_id"`
	UserName string `json:"user_name"`
}

func messageHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("%v:%v\n", "Userid", r.Form.Get("user_id"))
	user_id := r.Form.Get("user_id")
	response := InChannelResponse{user_id, "in_channel"}
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	// user, err := api.GetUseName(message.)
	// if err != nil {
	// 	fmt.Printf("%s\n", err)
	// 	return
	// }
	// fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
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
	w.Write(js)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	response := InChannelResponse{"HELLO CHARLY", "in_channel"}
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

type PullRequest struct {
	HTMLUrl string `json:"html_url"`
	Number  int    `json:"number"`
	State   string `json:"state"`
	ByUser  User   `json:"user"`
}

type GithubResponse struct {
	Action             string      `json:"action"`
	PR                 PullRequest `json:"pull_request"`
	Body               string      `json:"body"`
	Labels             []Label     `json:"labels"`
	CreatedAt          string      `json:"created_at"`
	UpdatedAt          string      `json:"updated_at"`
	Assignee           User        `json:"assignee"`
	Assignees          []User      `json:"assignees"`
	RequestedReviewers []User      `json:"requested_reviewers"`
}

func githubPrHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var t GithubResponse
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	jsonErr := json.Unmarshal([]byte(body), &t)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	log.Println(t)
	for _, reviewer := range t.RequestedReviewers {
		log.Println(reviewer)
	}
	// To delete when reviewers are available
	assignee := t.Assignee
	log.Println(assignee)
	getUserByEmail(assignee.Email)

}

func getUserByEmail(email string) {
	user, err := api.GetUserByEmail(email)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
}

func registerEndpoints() {
	http.HandleFunc("/message", messageHandler)
	http.HandleFunc("/github-pr", githubPrHandler)
	http.HandleFunc("/", pingHandler)
	http.HandleFunc("/hello", helloHandler)
}

// RunServer This is the server runner
func RunServer() {
	registerEndpoints()
	http.ListenAndServe(":8000", nil)
}
