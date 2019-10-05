FROM takanabe/add-new-issues-to-project-column:0.1

LABEL "com.github.actions.name"="Add new issues to a designate project column"
LABEL "com.github.actions.description"="GitHub Actions adding new issues to a specified project column automatically"
LABEL "com.github.actions.icon"="terminal"
LABEL "com.github.actions.color"="purple"

LABEL "repository"="https://github.com/takanabe/add-new-issues-to-project-column"
LABEL "homepage"="https://github.com/takanabe/add-new-issues-to-project-column"
LABEL "maintainer"="Takayuki Watanabe <takanabe.w@gmail.com>"

CMD ["/app/main"]