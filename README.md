# add-issues-to-project-column

![Docker Automated build](https://img.shields.io/docker/automated/takanabe/add-new-issues-to-project-column.svg?logo=docker)
![License](https://img.shields.io/github/license/takanabe/add-new-issues-to-project-column.svg)

GitHub Actions adding new GitHub Issues to a specified GitHub Project column automatically :recycle:.

## Usage

GitHub Projects belong to organizations, repositories, and users. This GitHub action currently does not support user-based GitHub Project. Create `.github/workflows/issues.yml` file on your repository and edit like below.

For any type of GitHub Projects, you need to change `GITHUB_PROJECT_URL` and `GITHUB_PROJECT_COLUMN_NAME` depending on your GitHub Project URL and column name to which you want to add new cards.

### Repository-based project

```yml
name: Add a new GitHub Project card linked to a GitHub issue to a specified project column
on: issues
jobs:
  add-new-issues-to-project-column:
    runs-on: ubuntu-latest
    steps:
    - name: add-new-issues-to-repository-based-project-column
      uses: docker://takanabe/add-new-issues-to-project-column:v0.0.1
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        GITHUB_PROJECT_URL: https://github.com/takanabe/add-new-issues-to-project-column/projects/1
        GITHUB_PROJECT_COLUMN_NAME: To do
```

### Organization-based project

```yml
name: Add a new GitHub issue to a designate project column
on: issues
jobs:
  add-new-issues-to-project-column:
    runs-on: ubuntu-latest
    steps:
    - name: add-new-issues-to-organization-based-project-column
      uses: docker://takanabe/add-new-issues-to-project-column:v0.0.1
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_PERSONAL_TOKEN_TO_ADD_PROJECT }}
        GITHUB_PROJECT_URL: https://github.com/orgs/organization_name/projects/1
        GITHUB_PROJECT_COLUMN_NAME: test
```

1. Replace the URL set on `GITHUB_PROJECT_URL` to the URL of your repository project to place issues
1. Replace the URL set on `GITHUB_PROJECT_COLUMN_NAME` to the string which your repository project has and want to place issues
1. Replace the secret set on ${{ secrets.GITHUB_PERSONAL_TOKEN_TO_ADD_PROJECT }} to your personal GitHub token
   1. Create a new personal access token from https://github.com/settings/tokens
   1. Create a new personal access token from https://github.com/organization_name/repository_name/settings/secrets with the value of personal access token you created above
   1. Replace the pesonal token name from ${{ secrets.GITHUB_PERSONAL_TOKEN_TO_ADD_PROJECT }} to ${{ secrets.YOUR_NEW_PERSONAL_TOKEN }}

### User-based project

User-based project is not supported yet

### Environment variables

| Environment variable       | Value                                                                                                                                       | Description                                                                                                                                                                                                                                                                                                                                                                                                                                |
| :------------------------- | :------------------------------------------------------------------------------------------------------------------------------------------ | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| GITHUB_TOKEN               | ${{ secrets.GITHUB_TOKEN }}, ${{ secrets.GITHUB_PERSONAL_ACCESS_TOKEN }}                                                                    | An Access token to access your issues and projects. if you use repository-based projects, ${{ secrets.GITHUB_TOKEN }} provides appropriate access privileges to this GitHub action ([See](https://help.github.com/en/articles/virtual-environments-for-github-actions#github_token-secret)). If that is not enough, you need to pass ${{ secrets.GITHUB_PERSONAL_ACCESS_TOKEN }} by issuing personal access token with appropriate grants. |
| GITHUB_PROJECT_URL         | https://github.com/username/reponame/projects/1, https://github.com/orgname/reponame/projects/1, https://github.com/orgs/orgname/projects/1 | A GitHub Project URL you want to use                                                                                                                                                                                                                                                                                                                                                                                                       |
| GITHUB_PROJECT_COLUMN_NAME | Anything (e.g: To Do)                                                                                                                       | A GitHub Project column name you want to place new issues                                                                                                                                                                                                                                                                                                                                                                                  |
| DEBUG                      | Anything (e.g: true)                                                                                                                        | A flag to produce debug messages for this GitHub Actions if this environment variable exists                                                                                                                                                                                                                                                                                                                                               |

## Development

### Build Docker image and update DockerHub

Change `IMAGE_NAME`, `DOCKER_REPO` and `TAG_NAME` in `Makefile` based on your DockerHub settings.

```bash
make
```

Except for `sandbox` tag, [`takanabe/add-new-issues-to-project-column`](https://hub.docker.com/r/takanabe/add-new-issues-to-project-column/tags) lists production ready Docker images matching [GitHub release tag](https://github.com/takanabe/add-new-issues-to-project-column/releases).

## License

[Apache 2.0](https://github.com/takanabe/add-new-issues-to-project-column/blob/master/LICENSE)
