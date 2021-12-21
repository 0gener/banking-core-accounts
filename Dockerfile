# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR ${GOPATH}/github.com/0gener/banking-core-accounts

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go .
COPY data data
COPY proto proto

RUN go build -o /banking-core-accounts

EXPOSE 5000

CMD [ "/banking-core-accounts" ]