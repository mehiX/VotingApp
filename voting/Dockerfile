
#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/app
RUN apk add --no-cache git

# copy this from original location (above apk add) to use caching
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/voting /app
# remember to copy the templates and other directories
COPY ./tmpl/ /tmpl/

ENTRYPOINT /app
LABEL Name=voting Version=0.0.1
LABEL MAINTAINER="mihai.oprea@gmail.com"
EXPOSE 8080
