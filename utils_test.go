package main

import "testing"

func TestGetUserByIDNotSubscribed(t *testing.T) {
	_, err := GetUserByID("testID")
	if err == nil {
		t.Errorf("GetUserByID(testID) do not return an error; want no user found")
	}
}


func TestGetUserByIDSubscribed(t *testing.T) {
	ProtSubscribedUsers["testID"] = SubscribedUser{
		SlackUserID:   "testID",
		SlackUserName: "testUserName",
		SlackChannelId: "testChannelID",
		GithubUser: "testGithubUser",
	}

	user, err := GetUserByID("testID")
	if err != nil {
		t.Errorf("GetUserByID(testID) = ``; want %s", user)
	}

	if user.SlackUserID != "testID" {
		t.Errorf("SlackUserID = %s; want testID", user.SlackUserID)
	}
}
