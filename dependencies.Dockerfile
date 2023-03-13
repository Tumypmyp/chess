FROM golang:1.20 AS dep
# Add the module files and download dependencies.

# ENV GO111MODULE=on

WORKDIR /chess

COPY go.mod ./
COPY go.sum ./

RUN go mod download
# Add the shared packages.
# COPY ./data /go/src/app/data
# COPY ./util /go/src/app/util