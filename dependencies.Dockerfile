FROM golang:1.20 AS dependencies

WORKDIR /chess

COPY go.mod ./
COPY go.sum ./

RUN go mod download
# Add the shared packages.
# COPY ./data /go/src/app/data
# COPY ./util /go/src/app/util