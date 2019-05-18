workflow "Add a new GitHub issue to a designate project column" {
  resolves = ["add-new-issues-to-project-column"]
  on = "issues"
}

action "add-new-issues-to-project-column" {
  uses = "./"
  // args = ["PROJECT NAME", "PROJECT_COLUMN_NAME"]
  args = ["test project", "To do"]
  secrets = ["GITHUB_TOKEN"]
}
