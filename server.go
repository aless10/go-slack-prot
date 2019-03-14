package main

import (
	"log"
	"net/http"
	"strings"
)

func registerEndpoints() {
	http.HandleFunc("/pr-list", commandHandler)
	http.HandleFunc("/repo-list", commandHandler)
	http.HandleFunc("/list-response", listResponseHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/github-pr", githubPrHandler)
	http.HandleFunc("/", pingHandler)
}

// RunServer This is the server runner
func RunServer() (err error) {
	registerEndpoints()
	err = http.ListenAndServe(strings.Join([]string{ServerConfig.Host, ServerConfig.Port}, ":"), nil)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
