FROM golang:1.11-alpine

ENV CGO_ENABLED=0
ENV GOPATH=/go
ENV TMPDIR=/tmp

RUN apk add --update git

WORKDIR /app

ADD ./vendor /app/

ADD . .
RUN go build .
