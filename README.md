# add-issues-to-project-column

GitHub Actions adding new issues to a designate project column automatically :recycle:

## Usage

Create `.github/main.workflow` file on your repository and edit like below.

```hcl
workflow "Add a new GitHub issue to the designate project column" {
  resolves = ["add-new-issues-to-project-column"]
  on = "issues"
}

action "add-new-issues-to-project-column" {
  uses = "takanabe/add-new-issues-to-project-column@master"
  env = {
    PROJECT_NAME = "PROJECT_NAME"
    PROJECT_COLUMN_NAME = "PROJECT_COLUMN_NAME"
  }
  secrets = ["GITHUB_TOKEN"]
}
```

You need to change `PROJECT_NAME` and `PROJECT_COLUMN_NAME` depending on your GitHub Project name and its column name.

## Development

To develop GitHub Actions in your local environment, use [act](https://github.com/nektos/act).

VSCode debug config example is [here](https://github.com/takanabe/add-new-issues-to-project-column/blob/master/.vscode/launch.json.example). `GITHUB_TOKEN` is necessary to access your repository and  project.

## License

[Apache 2.0](https://github.com/takanabe/add-new-issues-to-project-column/blob/master/LICENSE)
