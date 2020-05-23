FROM golang:alpine as builder

WORKDIR /go/src/app

RUN apk add --no-cache git

COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/bin/results /app
COPY ./tmpl/ /tmpl/

ENTRYPOINT /app

LABEL Name=results Version=0.0.1
EXPOSE 8080
