FROM dependencies AS build

WORKDIR /chess
# Copy the application source code.
COPY leaderboard_service/*.go leaderboard_service/
COPY leaderboard/*.go leaderboard/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o /go/bin/server ./leaderboard_service

ENTRYPOINT [ "/go/bin/server" ]

FROM golang:alpine

COPY --from=build /go/bin/server /go/bin/server
ENTRYPOINT [ "/go/bin/server" ]
