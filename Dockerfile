FROM golang:latest

WORKDIR /go/src/
COPY ./cmd ./cmd
COPY ./go.mod ./
COPY ./.env ./

RUN go mod tidy
RUN go build cmd/main.go
CMD ["./main", "postgres"]
