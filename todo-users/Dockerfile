# === ЭТАП 1: Сборка бинарника ===
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Копируем оба нужных модуля — todo-proto как зависимость, todo-users как сам сервис
COPY todo-proto/ ./todo-proto/
COPY todo-users/ ./todo-users/

WORKDIR /build/todo-users
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

# === ЭТАП 2: Финальный образ ===
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /build/todo-users/app .
EXPOSE 8080
EXPOSE 9090
CMD ["./app"]