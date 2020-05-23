FROM golang:alpine as builder

WORKDIR /go/src/app

RUN apk add --no-cache git

COPY *.go .
RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/app /app

ENTRYPOINT ["/app"]

LABEL Name=generator Version=0.0.1
LABEL Maintainer="mihai.oprea@gmail.com"