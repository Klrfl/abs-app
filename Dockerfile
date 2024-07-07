FROM golang:1.22-alpine as base

WORKDIR usr/src/app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . ./
RUN go build -v -o /usr/local/bin/app

EXPOSE 8080

CMD ["app"]
