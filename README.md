
# Trello-app

  
Trello is a project management tool that uses a Kanban-style board to help teams organize and prioritize tasks. An integration with another tool allows the two systems to communicate and share data with each other, making it easier for teams to manage their workflows.

  

One popular integration for Trello is with Jira, a project management tool for software development teams. This integration allows teams to link Jira issues with Trello cards, and see the status of development tasks directly on their Trello boards. Additionally, team members can use the integration to create and update Jira issues directly from Trello.

  

Another popular integration for Trello is with Slack, a messaging platform for teams. This integration allows teams to receive notifications about Trello activity in a Slack channel, and take actions on Trello cards directly from Slack messages.

  

There are a lot of trello integration option like Google Drive, Microsoft teams, Evernote, DropBox, Mailchimp, etc.

  

Trello offers a wide range of integrations with various tools to allow teams to customize their workflow and make the most of the platform.

  

## Getting started

you should clone the repository, below is an example of how to do it

```bash
npm install
git clone https://github.com/juliotorresmoreno/trello-app.git
```

we create or edit the .env file and make sure it has the following content, then the trello documentation to get the key and token parameters.

[introduction](https://developer.atlassian.com/cloud/trello/guides/rest-api/api-introduction/)

The TRELLO_BOARD_ID parameter is the one that appears when you click on the board you are going to use.

  

```
PORT=3000
TRELLO_KEY=
TRELLO_TOKEN=
TRELLO_SERVER=https://api.trello.com
TRELLO_BOARD_ID=
ENV=development
```

## Run tests
```bash
 go test ./... -v
```
### Coverage
```bash
go clean -testcache && go test ./... -coverprofile=test/coverage.out
go tool cover -html=test/coverage.out -o test/coverage.html
browse test/coverage.html # only for unix like
```

## Run project on Docker
Required docker and docker-compose, please check:
* [Docker](https://docs.docker.com/get-docker/)
* [Docker-compose](https://docs.docker.com/compose/install/)
```bash
docker-compose up -d
```

## Links of interest

[API](http://localhost:3000)

[doc](http://localhost:3000/api/v1/docs)

