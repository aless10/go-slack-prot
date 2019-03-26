# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:1.12 AS build-env

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/prot

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY ./prot .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
# RUN go get -d -v ./...

# Install the package
# RUN go install -v ./...

WORKDIR $GOPATH/src/prot/cmd

RUN ls -la

RUN GOOS=linux go build .
# RUN go build .

RUN echo "Build Complete"

FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk*
RUN apk add --no-cache bash
WORKDIR /prot

ARG LOG_DIR=/prot-logs

# Create Log Directory
RUN mkdir -p ${LOG_DIR}
RUN touch ${LOG_DIR}/app_log.log

COPY --from=build-env /go/src/prot/ /prot

EXPOSE 8080