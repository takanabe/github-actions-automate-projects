# add-issues-to-project-column

[dockerhub]: https://hub.docker.com/r/takanabe/add-new-issues-to-project-column
[license]: https://github.com/takanabe/add-new-issues-to-project-column/blob/master/LICENSE

[![Docker Automated build](https://img.shields.io/docker/automated/takanabe/add-new-issues-to-project-column.svg?logo=docker)][dockerhub]
[![License](https://img.shields.io/github/license/takanabe/add-new-issues-to-project-column.svg)][license]

GitHub Actions adding GitHub Issues & Pull requests to a specified GitHub Project column automatically :recycle:. This GitHub Action is inspired by https://github.com/masutaka/github-actions-all-in-one-project

## Usage

GitHub Projects belong to organizations, repositories, and users. This GitHub action currently does not support user-based GitHub Project. Create `.github/workflows/issues.yml` file on your repository and edit like below.

For any type of GitHub Projects, you need to change `GITHUB_PROJECT_URL` and `GITHUB_PROJECT_COLUMN_NAME` depending on your GitHub Project URL and column name to which you want to add new cards.

### Repository-based project

```yml
name: Add a new GitHub Project card linked to a GitHub issue to a specified project column
on: [issues, pull_request]
jobs:
  add-new-issues-to-project-column:
    runs-on: ubuntu-latest
    steps:
    - name: add-new-issues-to-repository-based-project-column
      uses: docker://takanabe/add-new-issues-to-project-column:v0.0.1
      if: github.event_name == 'issues' && github.event.action == 'opened'
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        GITHUB_PROJECT_URL: https://github.com/takanabe/add-new-issues-to-project-column/projects/1
        GITHUB_PROJECT_COLUMN_NAME: To do
    - name: add-new-prs-to-repository-based-project-column
      uses: docker://takanabe/add-new-issues-to-project-column:v0.0.1
      if: github.event_name == 'pull_request' && github.event.action == 'opened'
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        GITHUB_PROJECT_URL: https://github.com/takanabe/add-new-issues-to-project-column/projects/1
        GITHUB_PROJECT_COLUMN_NAME: To do
```

1. Replace the URL set on `GITHUB_PROJECT_URL` to the URL of your repository project to place issues/pull-requests
1. Replace the URL set on `GITHUB_PROJECT_COLUMN_NAME` to the string which your repository project has and want to place issues/pull-requests

### Organization-based project

```yml
name: Add a new GitHub issue to a designate project column
on: [issues, pull_request]
jobs:
  add-new-issues-to-project-column:
    runs-on: ubuntu-latest
    steps:
    - name: add-new-issues-to-organization-based-project-column
      uses: docker://takanabe/add-new-issues-to-project-column:v0.0.1
      if: github.event_name == 'issues' && github.event.action == 'opened'
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_PERSONAL_TOKEN_TO_ADD_PROJECT }}
        GITHUB_PROJECT_URL: https://github.com/orgs/organization_name/projects/1
        GITHUB_PROJECT_COLUMN_NAME: To Do
```

1. Replace the URL set on `GITHUB_PROJECT_URL` to the URL of your repository project to place issues/pull-requests
1. Replace the URL set on `GITHUB_PROJECT_COLUMN_NAME` to the string which your repository project has and want to place issues/pull-requests
1. Replace the secret set on ${{ secrets.GITHUB_PERSONAL_TOKEN_TO_ADD_PROJECT }} to your personal GitHub token
   1. Create a new personal access token from https://github.com/settings/tokens
   1. Create a new personal access token from https://github.com/organization_name/repository_name/settings/secrets with the value of personal access token you created above
   1. Replace the pesonal token name from ${{ secrets.GITHUB_PERSONAL_TOKEN_TO_ADD_PROJECT }} to ${{ secrets.YOUR_NEW_PERSONAL_TOKEN }}

### User-based project

User-based project is not supported yet

## Configurations

### Environment variables

| Environment variable       | Value                                                                                                                                       | Description                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| :------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------ | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| GITHUB_TOKEN               | ${{ secrets.GITHUB_TOKEN }}, ${{ secrets.GITHUB_PERSONAL_ACCESS_TOKEN }}                                                                    | An Access token to access your repository and projects. if you use repository-based projects, ${{ secrets.GITHUB_TOKEN }} provides appropriate access privileges to this GitHub action ([See](https://help.github.com/en/articles/virtual-environments-for-github-actions#github_token-secret)). If that is not enough, you need to pass ${{ secrets.GITHUB_PERSONAL_ACCESS_TOKEN }} by issuing personal access token with appropriate grants. |
| GITHUB_PROJECT_URL         | https://github.com/username/reponame/projects/1, https://github.com/orgname/reponame/projects/1, https://github.com/orgs/orgname/projects/1 | A GitHub Project URL you want to use                                                                                                                                                                                                                                                                                                                                                                                                           |
| GITHUB_PROJECT_COLUMN_NAME | Anything (e.g: To Do)                                                                                                                       | A GitHub Project column name you want to place issues/pull-requests                                                                                                                                                                                                                                                                                                                                                                            |
| DEBUG                      | Anything (e.g: true)                                                                                                                        | A flag to produce debug messages for this GitHub Actions if this environment variable exists                                                                                                                                                                                                                                                                                                                                                   |

### Condition with contexts

You can easily detect [event contexts](https://help.github.com/en/articles/contexts-and-expression-syntax-for-github-actions#github-context) and use them in if statements. Here are some lists of the useful contexts for this GitHub action.

| Property name       | Values                                                                                                                                                                              | Description                                                                                                                                                                                                      |
| ------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| github.event.action | opened, closed, edited, and so on                                                                                                                                                   | The name of actions (references for [issues](https://developer.github.com/v3/activity/events/types/#issuesevent) and for [pull_request](https://developer.github.com/v3/activity/events/types/#pullrequestevent) |
| github.event_name   | [issues](https://developer.github.com/v3/activity/events/types/#webhook-event-name-19), [pull_quests](https://developer.github.com/v3/activity/events/types/#webhook-event-name-33) | The name of the event that triggered the workflow run                                                                                                                                                            |

## Development

### Build Docker image and update DockerHub

Change `IMAGE_NAME`, `DOCKER_REPO` and `TAG_NAME` in `Makefile` based on your DockerHub settings.

```bash
make
```

Except for `sandbox` tag, [`takanabe/add-new-issues-to-project-column`](https://hub.docker.com/r/takanabe/add-new-issues-to-project-column/tags) lists production ready Docker images matching [GitHub release tag](https://github.com/takanabe/add-new-issues-to-project-column/releases).

## License

[Apache 2.0](https://github.com/takanabe/add-new-issues-to-project-column/blob/master/LICENSE)
