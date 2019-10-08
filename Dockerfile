FROM takanabe/github-actions-automate-projects:latest

LABEL "com.github.actions.name"="Automate projects"
LABEL "com.github.actions.description"="GitHub Actions adding GitHub Issues & Pull requests to the specified GitHub Project column automatically ♻️"
LABEL "com.github.actions.icon"="terminal"
LABEL "com.github.actions.color"="purple"

LABEL "repository"="https://github.com/takanabe/github-actions-automate-projects"
LABEL "homepage"="https://github.com/takanabe/github-actions-automate-projects"
LABEL "maintainer"="Takayuki Watanabe <takanabe.w@gmail.com>"

CMD ["/app/main"]
