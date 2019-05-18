FROM golang:1.11

LABEL "com.github.actions.name"="Add new issues to a designate project column"
LABEL "com.github.actions.description"="GitHub Actions place new issues to a designate GitHub project column"
LABEL "com.github.actions.icon"="terminal"
LABEL "com.github.actions.color"="purple"

LABEL "repository"="https://github.com/takanabe/add-new-issues-to-project-column"
LABEL "homepage"="https://github.com/takanabe/add-new-issues-to-project-column"
LABEL "maintainer"="Takayuki Watanabe <takanabe.w@gmail.com>"

# Force the go compiler to use modules
ENV GO111MODULE=on

RUN mkdir /app
COPY . /app/
WORKDIR /app

RUN go mod download
RUN go build -o main .

ENTRYPOINT ["/app/main"]