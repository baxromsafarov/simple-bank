# ---------- Stage 1: Build ----------
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Устанавливаем git (часто нужен для go mod download)
RUN apk add --no-cache git

# Копируем файлы модулей отдельно (для кеширования)
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download && go mod verify

# Копируем остальной исходный код
COPY . .

# Собираем бинарник
RUN go build -o main main.go


# ---------- Stage 2: Run ----------
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Для работы с базой может понадобиться ca-certificates
RUN apk add --no-cache ca-certificates

# Копируем бинарник из builder-стейджа
COPY --from=builder /app/main .

# (Не обязательно, но можно) копировать app.env внутрь образа:
# COPY app.env .

# Указываем порт, который слушает приложение
EXPOSE 8080

# Запускаем приложение
CMD ["/app/main"]
