FROM golang:alpine as builder

WORKDIR /go/src/app

RUN apk add --no-cache git

COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/worker /app

ENTRYPOINT /app

LABEL Name=Worker Version=0.0.1
LABEL MAINTAINER="mihai.oprea@gmail.com"