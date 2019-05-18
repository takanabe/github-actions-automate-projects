workflow "Add a new GitHub issue to a designate project column" {
  resolves = ["add-new-issues-to-project-column"]
  on = "issues"
}

action "add-new-issues-to-project-column" {
  uses = "./"
  //env = {
  //  PROJECT_NAME  = "PROJECT_NAME"
  //  PROJECT_COLUMN_NAME = "PROJECT_COLUMN_NAME"
  //}
  env = {
    PROJECT_NAME  = "test project"
    PROJECT_COLUMN_NAME = "To do"
  }
  secrets = ["GITHUB_TOKEN"]
}
