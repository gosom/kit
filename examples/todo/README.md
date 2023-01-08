# A simple todo app using event sourcing


Just to demostrate the `es` package

```
docker-compose up -d
go run cmd/app/main.go
```

Create a Todo:

```
curl -XPOST 'http://localhost:8080/todo' \
--header 'Content-Type: application/json' \
--data-raw '{
    "title": "do task 1"
}'
```

Mark it completed (replace the uuid from the one you get from above):

```
curl -XPATCH 'http://localhost:8080/todo/b656ad49-4c8b-449a-9787-7407dc73ff47' \
--header 'Content-Type: application/json' \
--data-raw '{
    "status": "completed"
}'
```

Get Aggregate:

```
curl 'http://localhost:8080/domain/aggregates/todo-b656ad49-4c8b-449a-9787-7407dc73ff47'
```

Get Events:

```
curl 'http://localhost:8080/domain/events/todo-b656ad49-4c8b-449a-9787-7407dc73ff47'
```

Get Command:

```
curl 'http://localhost:8080/domain/commands/01GP8X6PC3J6YKE87MA1YZ0TK7'
```
