package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GetUserByID(id string) (SubscribedUser, error) {
	user := ProtSubscribedUsers[id]
	if user.SlackUserID == "" {
		Warning.Println("No user found. Please subscribe")
		return SubscribedUser{}, fmt.Errorf("no user found")
	}
	return user, nil

}

func GetUserGithubName(login string) (SubscribedUser, error) {

	for _, subscribedUser := range ProtSubscribedUsers {
		user, err := GetUserByID(subscribedUser.SlackUserID)
		if err != nil {
			Warning.Printf("User %s is not subscribed", subscribedUser.SlackUserName)
		}
		if user.GithubUser == login {
			return user, nil
		}
	}

	return SubscribedUser{}, nil
}

func createGithubClient() *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: Config.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}
