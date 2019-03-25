# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:1.12

ARG LOG_DIR=/logs

# Create Log Directory
RUN mkdir -p ${LOG_DIR}

# Environment Variables
ENV APP_LOG_FILE=${LOG_DIR}/app_log.log

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

ENV SLACK_TOKEN=${SLACK_TOKEN}
ENV SLACK_BOT_ACCESS_TOKEN=${SLACK_BOT_ACCESS_TOKEN}
ENV GITHUB_TOKEN=${GITHUB_TOKEN}
ENV ORGANIZATION=${ORGANIZATION}

RUN env

ENV HOST_ADDRESS=${HOST_ADDRESS}
ENV PORT=${PORT}

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...


# This container exposes port 8080 to the outside world
EXPOSE 8080

WORKDIR $GOPATH/src/prot/cmd

RUN go build .
# Run the executable
RUN ./cmd
