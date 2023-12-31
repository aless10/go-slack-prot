package main

import (
	"net/http"
	"strings"
)

func init() {
	http.HandleFunc("/pr-list", commandHandler)
	http.HandleFunc("/repo-list", commandHandler)
	http.HandleFunc("/list-response", listResponseHandler)
	http.HandleFunc("/subscribe", subscribeHandler)
	http.HandleFunc("/github-pr", githubPrHandler)
	http.HandleFunc("/", pingHandler)
}

// RunServer This is the server runner
func RunServer() (err error) {
	err = http.ListenAndServe(strings.Join([]string{ServerConfig.Host, ServerConfig.Port}, ":"), nil)
	if err != nil {
		Error.Println(err)
	}
	return err
}
