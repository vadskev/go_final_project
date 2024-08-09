FROM golang:1.22.4

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE ${TODO_PORT}

RUN go build -o bin/task_app ./cmd/main.go

CMD ["./bin/task_app"]