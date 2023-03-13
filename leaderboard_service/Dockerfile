FROM dependencies AS build

WORKDIR /chess
# Copy the application source code.
COPY ./ ./

RUN go build ./leaderboard
RUN go build -o /server ./leaderboard_service
ENTRYPOINT [ "/server" ]

# FROM alpine:latest
# COPY --from=build /server /server
# CMD [ "/server" ]
