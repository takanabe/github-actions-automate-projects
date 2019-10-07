package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-github/v25/github"
	"github.com/pkg/errors"
)

// issueEventPayload returns GitHub issue event payload kept in the payload file .
// GITHUB_EVENT_PATH keeps the path to the file that contains the payload of the event that triggered the workflow.
// See: https://developer.github.com/actions/creating-github-actions/accessing-the-runtime-environment/#environment-variables
func issueEventPayload() github.IssuesEvent {
	var jsonFilePath string
	_, ok := os.LookupEnv("GITHUB_ACTION_LOCAL")
	if ok {
		// Use local test payload
		// https://developer.github.com/v3/activity/events/types/#issuesevent
		var err error
		jsonFilePath, err = filepath.Abs("./payload/issue_event.json")
		if err != nil {
			errorLog(err)
		}
	} else {
		jsonFilePath = os.Getenv("GITHUB_EVENT_PATH")
	}
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		errorLog(errors.Wrap(err, "Failed to open json"))
	}
	defer jsonFile.Close()

	// read opened jsonFile as a byte array.
	jsonByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		errorLog(errors.Wrap(err, "Failed to read json as a byte array"))
	}

	payload := github.IssuesEvent{}
	err = json.Unmarshal(jsonByte, &payload)
	if err != nil {
		errorLog(errors.Wrap(err, "Failed to unmarshal JSON to Go Object"))
	}
	if payload.GetAction() != "opened" {
		infoLog("GitHub action interupts!!")
		infoLog("This issue is not new one :D")
		os.Exit(0)
	}

	return payload
}

// extractIssueID returns GitHub issue ID extracted from event payloads
func extractIssueID(payload github.IssuesEvent) (int64, error) {

	issueID := payload.GetIssue().GetID()

	if issueID == 0 {
		return 0, errors.New("Issue ID is 0. Failed to get issue id properly")
	}

	return issueID, nil
}
