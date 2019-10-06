package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/go-github/v25/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func main() {
	payload := issueEventPayload()

	issueID := extractIssueID(payload)
	infoLog("New issue is found!!")
	debugLog("Issue ID: %d\n", issueID)

	client, ctx := getGitHubClient()

	url := os.Getenv("GITHUB_PROJECT_URL")
	if url == "" {
		log.Println("[ERROR] Environment variable GITHUB_PROJECT_URL is not defined in your workflow file")
		os.Exit(1)
	}

	// Project API does not support find Project column ID by URL.
	// So, detecting project type by URL and using different API to get get Project ID are necessary.
	pjType, err := projectType(url)
	errCheck(err)

	parentResource, parentName, err := projectParentName(url)
	errCheck(err)

	var pjID int64
	if pjType == "repository" {
		pjID, err = projectIDByRepo(ctx, client, url, parentResource, parentName)
		errCheck(err)
	} else if pjType == "organization" {
		pjID, err = projectIDByOrg(ctx, client, url, parentName)
		errCheck(err)
	} else if pjType == "user" {
		log.Println("[ERROR] User project is not supported yet")
		os.Exit(1)
	}
	infoLog("Project type:%s\n", pjType)

	pjColumn := os.Getenv("GITHUB_PROJECT_COLUMN_NAME")
	if pjColumn == "" {
		log.Println("[ERROR] Environment variable PROJECT_COLUMN_NAME is not defined in your workflow file")
		os.Exit(1)
	}

	columnID, err := projectColumnID(ctx, client, pjID, pjColumn)
	errCheck(err)

	infoLog("Project card is being added to column %s\n", pjColumn)

	////
	// Add a new opened issue to a designate project column
	////
	err = addIssueToProject(ctx, client, issueID, columnID)
	errCheck(err)

	os.Exit(0)
}

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

// extractIssueID returns GitHub issue ID extracted from event payloads
func extractIssueID(payload github.IssuesEvent) int64 {
	return payload.GetIssue().GetID()
}

func getGitHubClient() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), ctx
}

// GitHub projects belong to repository, organization, and user. (https://developer.github.com/v3/projects/#projects)
// Each project type uses different endpoints to get project information.
// Thus, this function investigates project type based on given URL

func projectType(url string) (string, error) {
	if url == "" {
		return "", errors.New("GITHUB_PROJECT_URL is empty")
	}

	var projectType string
	regUser := regexp.MustCompile(`https://github\.com/users/.+/projects/\d`)
	regOrg := regexp.MustCompile(`https://github\.com/orgs/.+/projects/\d`)
	regRepo := regexp.MustCompile(`https://github\.com/(.+)/.+/projects/\d`) // golang does not support negative lookahead

	if regUser.MatchString(url) {
		projectType = "user"
	} else if regOrg.MatchString(url) {
		projectType = "organization"
	} else if regRepo.MatchString(url) && !regUser.MatchString(url) && !regOrg.MatchString(url) {
		projectType = "repository"
	} else {
		return "", errors.New("GITHUB_PROJECT_URL is an invalid URL")
	}

	return projectType, nil
}

func projectParentName(rawURL string) (string, string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", errors.New("Failed to parse URL")
	}
	// All project URL must be formed with https://github.com/PARENT_RESOURCE/PARENT_NAME/projects/\d style.
	// (e.g)
	//   - (repo) https://github.com/username/reponame/projects/1
	//   - (repo) https://github.com/orgname/reponame/projects/1
	//   - (org) https://github.com/orgs/orgname/projects/1
	//   - (user) https://github.com/users/username/projects/1
	// Thus, organization, user, and repository name can be extracted from given project URL as PARENT_NAME.
	// PARENT_NAME is necessary to get all types of projects but only the repository project needs PARENT_RESOURCE.
	path := strings.Split(u.Path, "/")
	return path[1], path[2], nil
}

func projectIDByRepo(ctx context.Context, client *github.Client, url, owner, repo string) (int64, error) {
	var projectID int64
	projects, res, err := client.Repositories.ListProjects(ctx, owner, repo, nil)
	err = validateGitHubResponse(res, err)
	if err != nil {
		return 0, err
	}

	if projects == nil {
		return 0, errors.New("There are no projects on the repository named " + repo)
	}

	for _, project := range projects {
		infoLog("project url: %v\n", project.GetHTMLURL())
		if project.GetHTMLURL() == url {
			projectID = project.GetID()
			infoLog("Project Name: %s\n", project.GetName())
			debugLog("Project ID: %d\n", projectID)
			break
		}
	}

	if projectID == 0 {
		return 0, errors.New("No such a project url: " + url)
	}
	return projectID, nil
}

func projectIDByOrg(ctx context.Context, client *github.Client, url, org string) (int64, error) {
	var projectID int64
	opt := &github.ProjectListOptions{
		ListOptions: github.ListOptions{
			PerPage: 200,
		},
	}

	projects, res, err := client.Organizations.ListProjects(ctx, org, opt)
	err = validateGitHubResponse(res, err)
	if err != nil {
		return 0, err
	}

	for _, project := range projects {
		if project.GetHTMLURL() == url {
			projectID = project.GetID()
			fmt.Println("Project Name: " + project.GetName())
			infoLog("Project ID: %d\n", projectID)
			break
		}
	}

	if projectID == 0 {
		return 0, errors.New("No such a project url: " + url)
	}
	return projectID, nil
}

func projectColumnID(ctx context.Context, client *github.Client, pjID int64, pjColumn string) (int64, error) {
	var columnID int64
	columns, res, err := client.Projects.ListProjectColumns(ctx, pjID, nil)
	err = validateGitHubResponse(res, err)
	if err != nil {
		return 0, err
	}

	for _, col := range columns {
		if col.GetName() == pjColumn {
			columnID = col.GetID()
			infoLog("Column Name: %s", col.GetName())
			debugLog("Column ID: %d\n", columnID)
			break
		}
	}

	if columnID == 0 {
		return 0, errors.New("No such a column name: " + pjColumn)
	}

	return columnID, nil
}

func addIssueToProject(ctx context.Context, client *github.Client, issueID int64, columnID int64) error {
	opt := &github.ProjectCardOptions{
		ContentID:   issueID,
		ContentType: "Issue",
	}
	card, res, err := client.Projects.CreateProjectCard(ctx, columnID, opt)

	err = validateGitHubResponse(res, err)
	if err != nil {
		return err
	}

	if card.GetID() == 0 {
		return errors.New("Failed to create a card")
	}

	infoLog("Created card %d! issue %d is placed to ColumnID %d", card.GetID(), issueID, columnID)
	return nil
}

func errCheck(err error) {
	if err != nil {
		errorLog(err)
		os.Exit(1)
	}
}

func validateGitHubResponse(res *github.Response, err error) error {
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			return errors.Wrap(err, "Hit GitHub API rate limit")
		}
		return errors.Wrap(err, "Failed to get results from GitHub")
	}

	if !(res.Response.StatusCode == http.StatusOK || res.Response.StatusCode == http.StatusCreated) {
		return errors.Errorf("Invalid status code: %s. Failed to get results from GitHub", res.Status)
	}
	return nil
}
