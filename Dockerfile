FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /todo-list ./main.go

FROM ubuntu:latest

COPY --from=builder /todo-list /todo-list

COPY --from=builder /app/web /web

EXPOSE 7540

ENV TODO_PASSWORD=qwerty95

CMD ["/todo-list"]