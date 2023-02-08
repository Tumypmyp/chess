FROM golang

WORKDIR /chess

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /server

CMD [ "/server" ]