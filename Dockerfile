FROM golang:1.19-alpine

ADD . /app

WORKDIR /app

RUN go mod download
RUN go build -o /trello-app

EXPOSE 3000

ENTRYPOINT [ "/trello-app" ]
