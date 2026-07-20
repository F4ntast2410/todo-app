# === ЭТАП 1: Сборка бинарника ===
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Копируем оба нужных модуля — todo-proto как зависимость, todo-telegram как сам сервис
COPY todo-proto/ ./todo-proto/
COPY todo-telegram/ ./todo-telegram/

WORKDIR /build/todo-telegram
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# === ЭТАП 2: Финальный образ ===
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /build/todo-telegram/app .
EXPOSE 8082
EXPOSE 9092
CMD ["./app"]