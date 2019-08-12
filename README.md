# add-issues-to-project-column

GitHub Actions adding new issues to a designate project column automatically :recycle:

## Usage

Create `.github/workflows/issues.yml` file on your repository and edit like below.

```yml
on: issues
name: Add a new GitHub issue to a designate project column
jobs:
  add-new-issues-to-project-column:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: add-new-issues-to-project-column
      uses: ./
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        PROJECT_COLUMN_NAME: To do
        PROJECT_NAME: test project
```

You need to change `PROJECT_NAME` and `PROJECT_COLUMN_NAME` depending on your GitHub Project name and its column name.

## Development

To develop GitHub Actions in your local environment, use [act](https://github.com/nektos/act).

VSCode debug config example is [here](https://github.com/takanabe/add-new-issues-to-project-column/blob/master/.vscode/launch.json.example). `GITHUB_TOKEN` is necessary to access your repository and  project.

## License

[Apache 2.0](https://github.com/takanabe/add-new-issues-to-project-column/blob/master/LICENSE)
