# syntax=docker/dockerfile:1

FROM golang:1.23 AS build

WORKDIR $GOPATH/src/github.com/brotherlogic/fokus

COPY go.mod ./
COPY go.sum ./

RUN mkdir proto
COPY proto/*.go ./proto/

RUN mkdir server
COPY server/*.go ./server/

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /fokus

##
## Deploy
##
FROM ubuntu:22.04
USER root:root

WORKDIR /
COPY --from=build /fokus /fokus

EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["/fokus"]