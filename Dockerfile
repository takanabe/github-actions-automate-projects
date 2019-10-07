FROM takanabe/github-actions-automate-projects:sandbox

LABEL "com.github.actions.name"="Add new issues to a designate project column"
LABEL "com.github.actions.description"="GitHub Actions adding new issues to a specified project column automatically"
LABEL "com.github.actions.icon"="terminal"
LABEL "com.github.actions.color"="purple"

LABEL "repository"="https://github.com/takanabe/github-actions-automate-projects"
LABEL "homepage"="https://github.com/takanabe/github-actions-automate-projects"
LABEL "maintainer"="Takayuki Watanabe <takanabe.w@gmail.com>"

CMD ["/app/main"]
