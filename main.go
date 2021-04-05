package main

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v25/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

func main() {
	// Event name is kept in GITHUB_EVENT_NAME
	// https://help.github.com/en/articles/virtual-environments-for-github-actions#default-environment-variables
	eventName := os.Getenv("GITHUB_EVENT_NAME")
	infoLog("Event name: %s\n", eventName)
	if !(eventName == "issues" || eventName == "pull_request") {
		infoLog("This GitHub event is neither issues nor pull_requests. Stop executing this action.")
		infoLog("Please add 'if github.event_name' to the workflow yaml by following https://github.com/takanabe/github-actions-automate-projects/blob/master/README.md ")
		os.Exit(0)
	}

	var err error
	client, ctx := getGitHubClient()

	url := os.Getenv("GITHUB_PROJECT_URL")
	if url == "" {
		errorLog(errors.New("Environment variable GITHUB_PROJECT_URL is not defined in your workflow file"))
		os.Exit(1)
	}

	// Project API does not support find Project column ID by URL.
	// So, detecting project type by URL and using different API to get get Project ID are necessary.
	pjType, err := projectType(url)
	errCheck(err)

	parentResource, parentName, err := projectParentName(url)
	errCheck(err)

	// eventID stores issue ID or pull-request ID
	var eventID int64
	var projectCards []*github.ProjectCard
	if eventName == "issues" {
		payload := issueEventPayload()
		
		eventID, err = extractIssueID(payload)
		errCheck(err)
		
		repoOwner, repoName, err := repoOwnerAndName(payload.Issue.RepositoryUrl)
	        errCheck(err)
		
		projectCards, err = getProjectCardsFromIssue(ctx, client, payload.Issue, repoOwner, repoName)
		errCheck(err)
	} else if eventName == "pull_request" {
		payload := pullRequestEventPayload()
		eventID, err = extractPullRequestID(payload)
		errCheck(err)
	}

	infoLog("Payload for %s extract correctly", eventName)
	debugLog("Target event ID: %d\n", eventID)

	var pjID int64
	if pjType == "repository" {
		pjID, err = projectIDByRepo(ctx, client, url, parentResource, parentName)
		errCheck(err)
	} else if pjType == "organization" {
		pjID, err = projectIDByOrg(ctx, client, url, parentName)
		errCheck(err)
	} else if pjType == "user" {
		errorLog(errors.New("User project is not supported yet"))
		os.Exit(1)
	}
	infoLog("Project type:%s\n", pjType)

	pjColumn := os.Getenv("GITHUB_PROJECT_COLUMN_NAME")
	if pjColumn == "" {
		errorLog(errors.New("Environment variable PROJECT_COLUMN_NAME is not defined in your workflow file"))
		os.Exit(1)
	}

	columnID, err := projectColumnID(ctx, client, pjID, pjColumn)
	errCheck(err)

	for _, card := range projectCards {
		if *card.ProjectID == pjID {
			// Check card still exists
			// Ignore errors - if the card has been deleted, we want to respect the workflow and create it
			c, _, _ := client.Projects.GetProjectCard(ctx, *card.ID)
			if c == nil {
				continue
			}
			infoLog("Project card is being moved to column %s\n", pjColumn)
			err = moveCardInProject(ctx, client, columnID, card)
			errCheck(err)
			os.Exit(0)
		}
	}

	infoLog("Project card is being added to column %s\n", pjColumn)

	////
	// Add a new opened issue to a designate project column
	////
	err = addToProject(ctx, client, eventID, columnID, eventName)
	errCheck(err)

	os.Exit(0)
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

func repoOwnerAndName(rawApiURL string) (string, string, error) {
	u, err := url.Parse(rawApiURL)
	if err != nil {
		return "", "", errors.New("Failed to parse URL")
	}
	// A "Get Repository API" URL must be formed with https://api.github.com/repos/REPO_OWNER/REPO_NAME style.
	// (e.g)
	//   - (repo) https://github.com/username/reponame/projects/1
	//   - (org) https://github.com/orgname/reponame/1
	// Thus, organization and user repository owner and name can be extracted from given project URL as REPO_OWNER
	// and REPO_NAME, respectively.
	path := strings.Split(u.Path, "/")
	return path[2], path[3], nil
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

func addToProject(ctx context.Context, client *github.Client, eventID, columnID int64, eventName string) error {
	opt := &github.ProjectCardOptions{}

	if eventName == "issues" {
		opt.ContentID = eventID
		opt.ContentType = "Issue"
	} else if eventName == "pull_request" {
		opt.ContentID = eventID
		opt.ContentType = "PullRequest"
	}

	card, res, err := client.Projects.CreateProjectCard(ctx, columnID, opt)

	err = validateGitHubResponse(res, err)
	if err != nil {
		return err
	}

	if card.GetID() == 0 {
		return errors.New("Failed to create a card")
	}

	infoLog("Created card %d! issue %d is placed to ColumnID %d", card.GetID(), eventID, columnID)
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

func getProjectCardsFromIssue(ctx context.Context, client *github.Client, issue *github.Issue, owner string, repo string) ([]*github.ProjectCard, error) {
	cards := []*github.ProjectCard{}
	// TODO: Pagination, but let's assume any project will be in the first set of results.
	events, resp, err := client.Issues.ListIssueEvents(ctx, owner, repo, *issue.Number, nil)
	err = validateGitHubResponse(resp, err)
	if err != nil {
		return cards, err
	}
	for _, event := range events {
		if event.ProjectCard != nil {
			cards = append(cards, event.ProjectCard)
		}
	}
	return cards, nil
}

func moveCardInProject(ctx context.Context, client *github.Client, columnID int64, card *github.ProjectCard) error {
	resp, err := client.Projects.MoveProjectCard(ctx, *card.ID, &github.ProjectCardMoveOptions{ColumnID: columnID, Position: "top"})
	err = validateGitHubResponse(resp, err)
	if err != nil {
		return err
	}
	return nil
}
