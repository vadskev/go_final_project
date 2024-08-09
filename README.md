## Особенности проекта
Приложение представляет собой TODO лист с веб-интерфейсом и API поддерживающим CRUD операции. Данный проект был реализован с целью обучения, познания принципов SOLID и чистой архитектуры.

### Проект использует:
* Стандартные CRUD-операции таблицы базы данных
* RESTful endpoints в распространенном формате
* Аутентификацию на основе JWT токена
* Управление конфигурацией приложения, в зависимости от среды
* Структурированное логирование
* Обработку ошибок с генерацией ответов
* Проверку валидности данных
* Контейнеризацию в Docker образе

### API поддерживает следующие операции:
* добавить задачу
* получить список задач
* удалить задачу
* получить параметры задачи
* изменить параметры задачи
* отметить задачу как выполненную

### Используемые в проекте пакеты Go:
* Routing: [go-chi](https://github.com/go-chi/chi)
* Fluent SQL generator for Go: [Squirrel](https://github.com/Masterminds/squirrel)
* General purpose extensions to golang's database/sql: [sqlx](https://github.com/jmoiron/sqlx)
* Loads environment variables from .env files: [GoDotEnv](https://github.com/joho/godotenv)
* Sqlite3 driver for go using database/sql: [go-sqlite3](https://github.com/mattn/go-sqlite3)
* Package zap provides fast, structured, leveled logging: [Zap](go.uber.org/zap)

## Начало работы

Если вы впервые сталкиваетесь с Go, [установите Go по инструкции](https://golang.org/doc/install) на свой компьютер. Для проекта требуется Go 1.22.4 или выше.

[Docker](https://www.docker.com/get-started) можно также установить, если вы не хотите настраивать окружение для работы. Проекту требуется Docker 17.05 или выше.

Также для работы в терминале, рекомендую установить [Taskfile](https://taskfile.dev/installation/).

После установки Go, Docker и TaskFile выполните следующие команды, чтобы начать работу:

```shell
## Запуск TODO-проекта
# скачиваем репозиторий
git clone https://github.com/vadskev/go_final_project.git
cd go_final_project

# создаем конфигурацию .env:
$ nano env/.env {
LOG_LEVEL=info
TODO_HOST=localhost
TODO_PORT=7540
TODO_DBFILE=./scheduler.db
TODO_PASSWORD=secret_pass
}

# скачиваем зависимости:
go mod tidy

# запускаем сервер:
go run ./cmd
```

```shell
## Создание Docker-образа
img-build:
	docker build -t task_app .
	
## Запуск контейнера
img-run:
	docker run --env-file ./env/.env -d -p ${TODO_PORT}:${TODO_PORT} --name task-app-api task_app

## Остановка контейнера
img-stop:
	docker stop task-app-api
```

```shell
## Запуск тестов
go test ./tests
```

Проект будет доступен по адресу: `http://localhost:7540`. Сервер предоставляет следующие endpoint:

* `GET /api/task{id}`: возвращает задачу
* `POST /api/task{id}`: добавляет задачу
* `PUT /api/task`: обновляет задачу
* `DELETE /api/task{id}`: удаляет задачу
* `GET /api/tasks{search}`: возвращает список ближайших задач, возможна фильтрация с помощью параметра `search`
* `GET /api/nextdate{now}{date}{repeat}`: вычисляет следующую дату для задачи в соответствии с указанным правилом
* `POST /api/task/done`: делает задачу выполненной
* `POST /api/signin`: сверяет вводимый пароль на странице `/login.html` с хранимым в переменной окружения, в случае успеха - создает JWT-токен и записывает в Cookie `token`.

## Проект имеет следующую структуру:
```
go_final_project/
├── cmd/                  стартавая точка работы проекта
├── env/                  переменные окружения, пример .env файла
├── internal/             приватная директория приложения, библиотеки
│   ├── api/              Responce ответы
│   ├── app/              сборка приложения
│   └── config/           конфигурация приложения
│       └── env/          сбор переменных окружения
│   ├── handlers/         обработчики HTTP запросов
│   ├── logger/           создание logger
│   └── middleware/       middlewares
│       ├── auth/         middleware для проверки авторизации
│       └── logger/       middleware логирования запросов
│   ├── models/           слой моделей
│   ├── nextdate/         функция для работы с датой
│   ├── storage/          слой для работы с SQLite базой
├── tests/                тесты для проверки API
└── web/                  содержит файлы фронтенда
```

Каталоги первого уровня `cmd`, `internal` обычно встречаются в других популярных проектах Go, как описано 
в разделе [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

Пакет `internal` структурирован в соответствии с [screaming architecture](https://blog.cleancoder.com/uncle-bob/2011/09/30/Screaming-Architecture.html). 
Например, пакет `transport` содержит логику, связанную с передачей данных в другие объекты.

В каждом пакете функций код организован по уровням (API, service, repository) в соответствии с рекомендациями по зависимостям, 
описанными в [clean architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html).

## Переменные окружения
* `LOG_LEVEL=info`
* `TODO_HOST=localhost`
* `TODO_PORT=7540`
* `TODO_DBFILE=./scheduler.db`
* `TODO_PASSWORD=secret_pass`

## Управление конфигурацией

Конфигурация приложения представлена в директории `internal/config/*`.

При запуске приложение загружает конфигурацию из переменных окружения. Путь к переменным окружения должен быть расположен в по адресу `./env/.env`.

