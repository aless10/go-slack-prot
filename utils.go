package main

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
)

func getUserByID(id string) (SubscribedUser, error) {
	return ProtSubscribedUsers[id], nil

}

func getUserGithubName(login string) (SubscribedUser, error) {

	for _, subscribedUser := range ProtSubscribedUsers {
		user, err := getUserByID(subscribedUser.SlackUserID)
		if err != nil {
			log.Println(err)
		}
		if user.GithubUser == login {
			return user, nil
		}
	}
	nullUser := SubscribedUser{}
	return nullUser, nil
}

func createGithubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Config.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
