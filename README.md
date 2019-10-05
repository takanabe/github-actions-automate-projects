# add-issues-to-project-column

GitHub Actions adding new issues to a designate project column automatically :recycle:

## Usage

Create `.github/workflows/issues.yml` file on your repository and edit like below.

```yml
name: Add a new GitHub issue to a designate project column
on: issues
jobs:
  add-new-issues-to-project-column:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: add-new-issues-to-project-column
      uses: takanabe/add-new-issues-to-project-column@v0.0.3
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        GITHUB_PROJECT_URL: https://github.com/takanabe/add-new-issues-to-project-column/projects/1
        GITHUB_PROJECT_COLUMN_NAME: To do
```

You need to change `GITHUB_PROJECT_URL` and `GITHUB_PROJECT_COLUMN_NAME` depending on your GitHub Project URL and column name to which you want to add new cards.

## Development

To develop GitHub Actions in your local environment, use [act](https://github.com/nektos/act).

VSCode debug config example is [here](https://github.com/takanabe/add-new-issues-to-project-column/blob/master/.vscode/launch.json.example). `GITHUB_TOKEN` is necessary to access your repository and  project.


## Build Docker image and update DockerHub

```
$ docker build -f Dockerfile.build . -t add_issue_to_project
$ docker image tag add_issue_to_project takanabe/add-new-issues-to-project-column:0.1
$ docker login
$ docker push takanabe/add-new-issues-to-project-column:0.1
```

## License

[Apache 2.0](https://github.com/takanabe/add-new-issues-to-project-column/blob/master/LICENSE)
