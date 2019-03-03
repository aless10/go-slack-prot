package main
//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/nlopes/slack"
//	"io/ioutil"
//	"log"
//	"net/http"
//)
//
//type User struct {
//	Login string `json:"login"`
//	Email string
//}
//
//type Label struct {
//	LabelName string `json:"name"`
//}
//
//func commandHandler(w http.ResponseWriter, r *http.Request) {
//	command, err := slack.SlashCommandParse(r)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	user, err := Api.GetUserInfo(command.UserID)
//	if err != nil {
//		log.Fatalf("User with ID %s not found", command.UserID)
//	}
//	fmt.Printf("%v:%v\n", "UserEmail", user.Profile.Email)
//
//	// 1. For all the repos
//	repos, _, err := githubClient.Repositories.List(context.Background(), "", nil)
//	fmt.Println(repos)
//	for i, repo := range repos {
//		fmt.Println(i, repo.GetPullsURL())
//		//pullUrl := repo.GetPullsURL()
//		//githubClient.PullRequests.ListReviewers()
//
//	}
//	//for j, pr := range repo.GetP {
//	//	log.Println(j, pr["html_url"], pr)
//	//}
//	// 2. For all the PR
//	// 3. if user.Profile.Email in requested_reviewers email -> append PR info to a list
//	// 4. Return the list to slack
//	//x, _, err := githubClient.PullRequests.ListReviewers(ctx, owner, repo, number, opt *ListOptions)
//
//	response := InChannelResponse{user.Profile.FirstName, "in_channel"}
//	js, err := json.Marshal(response)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	_, err = w.Write(js)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
//
//type InChannelResponse struct {
//	Text         string `json:"text"`
//	ResponseType string `json:"response_type"`
//}
//
//func pingHandler(w http.ResponseWriter, r *http.Request) {
//	response := InChannelResponse{"pong", "in_channel"}
//	js, err := json.Marshal(response)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	_, err = w.Write(js)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
//
//func helloHandler(w http.ResponseWriter, r *http.Request) {
//	response := InChannelResponse{"HELLO CHARLY", "in_channel"}
//	js, err := json.Marshal(response)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	_, err = w.Write(js)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
//
//type PullRequest struct {
//	HTMLUrl string `json:"html_url"`
//	Number  int    `json:"number"`
//	State   string `json:"state"`
//	ByUser  User   `json:"user"`
//}
//
//type GithubResponse struct {
//	Action             string      `json:"action"`
//	PR                 PullRequest `json:"pull_request"`
//	Body               string      `json:"body"`
//	Labels             []Label     `json:"labels"`
//	CreatedAt          string      `json:"created_at"`
//	UpdatedAt          string      `json:"updated_at"`
//	Assignee           User        `json:"assignee"`
//	Assignees          []User      `json:"assignees"`
//	RequestedReviewers []User      `json:"requested_reviewers"`
//}
//
//func githubPrHandler(w http.ResponseWriter, r *http.Request) {
//	defer r.Body.Close()
//	var t GithubResponse
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	jsonErr := json.Unmarshal([]byte(body), &t)
//
//	if jsonErr != nil {
//		log.Fatal(jsonErr)
//	}
//	log.Println(t)
//	for _, reviewer := range t.RequestedReviewers {
//		log.Println(reviewer)
//	}
//	// To delete when reviewers are available
//	assignee := t.Assignee
//	log.Println(assignee)
//	getUserByEmail(assignee.Email)
//
//}
//
//func getUserByEmail(email string) {
//	user, err := Api.GetUserByEmail(email)
//	if err != nil {
//		fmt.Printf("%s\n", err)
//		return
//	}
//	fmt.Printf("ID: %s, Fullname: %s, Email: %s\n", user.ID, user.Profile.RealName, user.Profile.Email)
//}
//
////func registerEndpoints() {
//	//http.HandleFunc("/message", commandHandler)
//	//http.HandleFunc("/github-pr", githubPrHandler)
//	//http.HandleFunc("/", pingHandler)
//	//http.HandleFunc("/hello", helloHandler)
////}
//
////// RunServer This is the server runner
////func RunServer() {
////	registerEndpoints()
////	err := http.ListenAndServe(":8080", nil)
////	if err != nil {
////		log.Fatal(err)
////	}
////}
