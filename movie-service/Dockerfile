# Устанавливаем базовый образ
FROM golang:1.22.4-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем зависимости и устанавливаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Сборка исполняемого файла
RUN go build -o /movie-service ./cmd/movie/main.go

# Создаем минимальный образ для запуска
FROM alpine:latest
WORKDIR /root/

# Копируем скомпилированное приложение из builder'а
COPY --from=builder /movie-service .

# Экспорт порта
EXPOSE 8082

# Команда для запуска приложения
CMD ["./movie-service"]
