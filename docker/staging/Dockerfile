FROM golang:1.17-alpine as go-builder

RUN apk add --update alpine-sdk

RUN mkdir /app

ADD . /app

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main .

# running docker file
FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=go-builder /app/main /app/main

EXPOSE 8080