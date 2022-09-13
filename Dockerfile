FROM golang:1.18

WORKDIR /usr/src/app

COPY go.mod go.sum ./