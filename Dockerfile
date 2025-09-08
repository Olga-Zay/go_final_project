FROM golang:1.23 AS builder
WORKDIR /app
COPY . .
RUN go mod download
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o main

#создадим второй образ более облегченный на основе предыдущего и добавим в него файл базы данных с локального компьютера
FROM alpine:latest
WORKDIR /app
ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=12345
COPY ./web ./web
COPY --from=builder /app/main .
CMD ./main
