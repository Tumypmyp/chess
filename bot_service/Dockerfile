FROM dependencies AS build

WORKDIR /chess
# Copy the application source code.
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o /go/bin/server ./bot_service

ENTRYPOINT [ "/go/bin/server" ]

FROM golang:alpine
RUN apk --no-cache add ca-certificates

COPY --from=build /go/bin/server /go/bin/server
ENTRYPOINT [ "/go/bin/server" ]
