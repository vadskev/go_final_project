include ./env/.env

# RUN Local

build:
	go build -o bin/task_app ./cmd/main.go

run:
	go run ./cmd/main.go -config-path=./env/.env

test:
	go test ./tests -count=1

# RUN Docker

img-build:
	docker build -t task_app .

img-run:
	docker run --env-file ./env/.env -d -p ${TODO_PORT}:${TODO_PORT} --name task-app-api task_app

img-stop:
	docker stop task-app-api
