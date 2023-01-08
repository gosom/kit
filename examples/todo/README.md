# A simple todo app using event sourcing


Just to demostrate the `es` package

```
docker-compose up -d
go run cmd/app/main.go
```

Create a Todo:

```
curl --location --request POST 'http://localhost:8080/todo/commands' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "CreateTodo",
    "payload": {
        "id": "11186428-8f6c-11ed-bde4-13557563d9d6",
        "title": "do task 1"
    }
}'
```

Mark it completed:

```
curl --location --request POST 'http://localhost:8080/todo/commands' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "UpdateTodoStatus",
    "payload": {
        "id": "11186428-8f6c-11ed-bde4-13557563d9d6",
        "status": "completed"
    }
}'
```

Get Aggregate:

```
curl 'http://localhost:8080/todo/aggregates/todo-11186428-8f6c-11ed-bde4-13557563d9d6'
```

Get Events:

```
curl 'http://localhost:8080/domain/events/todo-11186428-8f6c-11ed-bde4-13557563d9d6'
```

Get Command:

```
curl 'http://localhost:8080/domain/commands/01GP8X6PC3J6YKE87MA1YZ0TK7'
```
