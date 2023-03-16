FROM dependencies AS build

WORKDIR /chess
# Copy the application source code.
COPY bot_service/*.go bot_service/
COPY leaderboard/*.go leaderboard/
COPY memory/*.go memory/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o /go/bin/server ./bot_service

ENTRYPOINT [ "/go/bin/server" ]

FROM golang:alpine

COPY --from=build /go/bin/server /go/bin/server
ENTRYPOINT [ "/go/bin/server" ]
