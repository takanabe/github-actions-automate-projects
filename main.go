package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-github/v25/github"
	"golang.org/x/oauth2"
)

func main() {

	payload := getIssueEvetnPayload()

	// Get issue ID, repository owner and repository name automatically
	issueID := getIssueID(payload)
	owner := getOwner(payload)
	repo := getRepo(payload)

	client, ctx := getGitHubClient()

	// Get project ID and column ID from arguments(project name, project column name), issue ID, repository owner and repository name
	targetProjectID := getProjectID(ctx, client, owner, repo)
	targetColumnID := getProjectColumnID(ctx, client, targetProjectID)

	// Add a new opened issue to a designate column
	addIssueToProject(ctx, client, issueID, targetColumnID)
}

func getIssueEvetnPayload() github.IssuesEvent {
	// GITHUB_EVENT_PATH keeps the path to a file that contains the payload of the event that triggered the workflow
	// See: https://developer.github.com/actions/creating-github-actions/accessing-the-runtime-environment/#environment-variables
	var jsonFilePath string
	_, ok := os.LookupEnv("GITHUB_ACTION_LOCAL")
	if ok {
		// Use local test payload
		// https://developer.github.com/v3/activity/events/types/#issuesevent
		var err error
		jsonFilePath, err = filepath.Abs("./payload/issue_event.json")
		if err != nil {
			log.Fatalf("[ERROR] No such a file: %e", err)
		}
	} else {
		jsonFilePath = os.Getenv("GITHUB_EVENT_PATH")
	}
	jsonFile, err := os.Open(jsonFilePath)
	if err != nil {
		log.Fatalf("[ERROR] Failed to open json: %e", err)
	}
	defer jsonFile.Close()

	// read opened jsonFile as a byte array.
	jsonByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("[ERROR] Failed to read json as a byte array: %e", err)
	}

	payload := github.IssuesEvent{}
	err = json.Unmarshal(jsonByte, &payload)
	if err != nil {
		log.Fatalf("[ERROR] Failed to unmarshal JSON to Go Object: %e", err)
	}
	if payload.GetAction() != "opened" {
		fmt.Println("GitHub action interupts!!")
		fmt.Println("This issue is not new one :D")
		os.Exit(0)
	}
	return payload
}

func getIssueID(payload github.IssuesEvent) int64 {
	var issueID int64
	issueID = payload.GetIssue().GetID()
	fmt.Printf("Issue ID: %d\n", issueID)
	return issueID
}

func getOwner(payload github.IssuesEvent) string {
	owner := payload.GetRepo().GetOwner().GetLogin()
	fmt.Printf("Owner: %s\n", owner)
	return owner
}

func getRepo(payload github.IssuesEvent) string {
	repo := payload.GetRepo().GetName()
	fmt.Printf("Repo: %s\n", repo)
	return repo
}

func getGitHubClient() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), ctx
}

func getProjectID(ctx context.Context, client *github.Client, owner, repo string) int64 {
	pj := os.Getenv("PROJECT_NAME")
	if pj == "" {
		log.Fatal("[Error] Environment variable PROJECT_NAME is not defined in your workflow file")
	}
	fmt.Printf("# Start searching project ID of %s...\n", pj)
	projects, _, err := client.Repositories.ListProjects(ctx, owner, repo, nil)

	if _, ok := err.(*github.RateLimitError); ok {
		log.Fatalf("[Error] Hit rate limit: %e", err)
	}

	var targetProjectID int64

	for _, project := range projects {
		if project.GetName() == pj {
			targetProjectID = project.GetID()
			fmt.Println("Find!!")
			fmt.Println("Project Name: " + project.GetName())
			fmt.Printf("Project ID: %d\n", targetProjectID)
			break
		}
	}

	if targetProjectID == 0 {
		log.Fatalf("[Error] No such a project name %s", pj)
	}

	return targetProjectID
}

func getProjectColumnID(ctx context.Context, client *github.Client, targetProjectID int64) int64 {
	var targetColumnID int64

	pjColumn := os.Getenv("PROJECT_COLUMN_NAME")
	if pjColumn == "" {
		log.Fatal("[Error] Environment variable PROJECT_COLUMN_NAME is not defined in your workflow file")
	}
	fmt.Printf("# Start searching project column ID of %s\n", pjColumn)

	columns, _, err := client.Projects.ListProjectColumns(ctx, targetProjectID, nil)

	if _, ok := err.(*github.RateLimitError); ok {
		log.Fatalf("[Error] Hit rate limit: %e", err)
	}

	for _, col := range columns {
		if col.GetName() == pjColumn {
			targetColumnID = col.GetID()
			fmt.Println("Find!!")
			fmt.Println("Column Name: " + col.GetName())
			fmt.Printf("Column ID: %d\n", targetColumnID)
		}
	}

	if targetColumnID == 0 {
		log.Fatalf("[Error] No such a column name %s", pjColumn)
	}

	return targetColumnID
}

func addIssueToProject(ctx context.Context, client *github.Client, issueID int64, targetColumnID int64) {
	fmt.Println("# Start adding a new issue to  project column")
	opt := &github.ProjectCardOptions{
		ContentID:   issueID,
		ContentType: "Issue",
	}
	card, _, err := client.Projects.CreateProjectCard(ctx, targetColumnID, opt)
	if _, ok := err.(*github.RateLimitError); ok {
		log.Fatalf("[Error] Hit rate limit: %e", err)
	}

	fmt.Printf("Created card %d! issue %d is placed to ColumnID %d", card.GetID(), issueID, targetColumnID)
}
