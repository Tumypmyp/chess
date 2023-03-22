FROM dependencies AS build

WORKDIR /chess
# Copy the application source code.
COPY player_service/*.go player_service/
COPY helpers/*.go helpers/
COPY memory/*.go memory/
COPY leaderboard/*.go leaderboard/
COPY game/*.go game/
COPY player/*.go player/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o /go/bin/server ./player_service

ENTRYPOINT [ "/go/bin/server" ]

FROM golang:alpine

COPY --from=build /go/bin/server /go/bin/server
ENTRYPOINT [ "/go/bin/server" ]
