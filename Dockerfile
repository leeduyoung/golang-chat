########################################################################
# Base install 
########################################################################
FROM golang:alpine AS base

RUN apk add --no-cache file
RUN apk add ca-certificates tzdata

RUN apk add git openssh-client wget
RUN apk update

RUN apk add --no-cache build-base

########################################################################
# Build
########################################################################
FROM base AS build
ARG authtoken

WORKDIR /work
COPY . .

# Configure Go
RUN export GO111MODULE=on && \
    export GIT_TERMINAL_PROMPT=1 && \
    export GOPRIVATE=github.com && \
    export GONOPROXY=none && \
    export CGO_ENABLED=1 && \
    git config --global --add url."https://$authtoken:x-oauth-basic@github.com/".insteadOf "https://github.com/"

RUN go mod tidy
RUN go build -a -ldflags '-s' -installsuffix cgo -o app

########################################################################
# alpine Stage layer 실행 
########################################################################
FROM alpine:latest

RUN apk add --no-cache file
RUN apk add ca-certificates tzdata

COPY --from=build /work/app .
CMD ["./app"]
