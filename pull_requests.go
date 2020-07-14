package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/go-github/v25/github"
	"github.com/pkg/errors"
)

// pullRequestEventPayload returns GitHub issue event payload kept in the payload file .
// GITHUB_EVENT_PATH keeps the path to the file that contains the payload of the event that triggered the workflow.
// See: https://developer.github.com/actions/creating-github-actions/accessing-the-runtime-environment/#environment-variables
func pullRequestEventPayload() github.PullRequestEvent {
	var jsonFilePath string
	_, ok := os.LookupEnv("GITHUB_ACTION_LOCAL")
	if ok {
		// Use local test payload
		// https://developer.github.com/v3/activity/events/types/#pullrequestevent
		var err error
		jsonFilePath, err = filepath.Abs("./payload/pull_request_event.json")
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

	payload := github.PullRequestEvent{}
	err = json.Unmarshal(jsonByte, &payload)
	if err != nil {
		errorLog(errors.Wrap(err, "Failed to unmarshal JSON to Go Object"))
	}

	return payload
}

// extractPullRequestID returns GitHub issue ID extracted from event payloads
func extractPullRequestID(payload github.PullRequestEvent) (int64, error) {

	prID := payload.GetPullRequest().GetID()

	if prID == 0 {
		return 0, errors.New("Pull Request ID is 0. Failed to get Pull Request id properly")
	}

	return prID, nil
}
