# Этап сборки
FROM golang:1.25-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей
COPY bot/go.mod bot/go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY bot/ ./

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bot .

# Финальный этап - минимальный образ
FROM alpine:latest

# Устанавливаем необходимые пакеты для работы с SSL/TLS
RUN apk --no-cache add ca-certificates tzdata

# Создаем пользователя для запуска приложения
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Копируем бинарный файл из этапа сборки
COPY --from=builder /app/bot .

# Меняем владельца файла
RUN chown appuser:appuser /app/bot

# Переключаемся на непривилегированного пользователя
USER appuser

# Указываем переменные окружения (можно переопределить при запуске)
ENV TELEGRAM_BOT_TOKEN=""
ENV COINMARKETCAP_API_KEY=""

# Запускаем приложение
CMD ["./bot"]
