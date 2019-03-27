# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:1.12-alpine

RUN apk add -U --no-cache git ca-certificates
# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/prot

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY ./prot .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

WORKDIR $GOPATH/src/prot/cmd

RUN go build -o main

RUN apk add --no-cache bash

ARG LOG_DIR=/prot-logs

# Create Log Directory
RUN mkdir -p ${LOG_DIR}
RUN touch ${LOG_DIR}/app_log.log

EXPOSE 8080

RUN rm -rf /prot/cmd

WORKDIR $GOPATH/bin

ENTRYPOINT ["cmd"]