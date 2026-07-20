# === ЭТАП 1: Сборка бинарника ===
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Копируем оба нужных модуля — todo-proto как зависимость, todo-tasks как сам сервис
COPY todo-proto/ ./todo-proto/
COPY todo-tasks/ ./todo-tasks/

WORKDIR /build/todo-tasks
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# === ЭТАП 2: Финальный образ ===
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/todo-tasks/app .
EXPOSE 8081
EXPOSE 9091
CMD ["./app"]