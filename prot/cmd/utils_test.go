package main

import (
	"log"
	"os"
	"testing"
)

func init() {
	testLogFile, err := os.OpenFile("./log/test_log.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
	}
	initLogs(testLogFile)
}

func TestGetUserByIDNotSubscribed(t *testing.T) {
	_, err := GetUserByID("testID")
	if err == nil {
		t.Errorf("GetUserByID(testID) do not return an error; want no user found")
	}
}

func TestGetUserByIDSubscribed(t *testing.T) {
	ProtSubscribedUsers["testID"] = SubscribedUser{
		SlackUserID:    "testID",
		SlackUserName:  "testUserName",
		SlackChannelId: "testChannelID",
		GithubUser:     "testGithubUser",
	}

	user, err := GetUserByID("testID")
	if err != nil {
		t.Errorf("GetUserByID(testID) = ``; want %s", user)
	}

	if user.SlackUserID != "testID" {
		t.Errorf("SlackUserID = %s; want testID", user.SlackUserID)
	}
}

func TestGetUserGithubNameNoUser(t *testing.T) {

	emptyUser := SubscribedUser{}
	user, err := GetUserGithubName("NoUseForAName")
	if user != emptyUser || err != nil {
		t.Errorf("user or err are not nil")
	}
}

func TestGetUserGithubName(t *testing.T) {

	ProtSubscribedUsers["testID"] = SubscribedUser{
		SlackUserID:    "testID",
		SlackUserName:  "testUserName",
		SlackChannelId: "testChannelID",
		GithubUser:     "testGithubUser",
	}

	user, err := GetUserGithubName("testGithubUser")
	if err != nil {
		t.Errorf("Expected User with SlackUserID = testID, got err %s", err)
	}

	if user.SlackUserID != "testID" {
		t.Errorf("Expected User with SlackUserID = testID, got %s", user.SlackUserID)
	}
}
