FROM golang:latest

WORKDIR /app

COPY *.go /app
COPY go.mod /app
COPY species.json /app
COPY das/ /app/das

RUN go mod tidy
RUN go get .

RUN pwd

RUN go build -o svr *.go

