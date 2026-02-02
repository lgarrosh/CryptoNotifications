# Telegram бот для получения котировок криптовалют

Telegram-бот на Go, который предоставляет актуальные котировки криптовалют через CoinMarketCap API.

## Возможности

- Получение котировок одной криптовалюты по команде `/price`
- Получение котировок нескольких криптовалют одновременно по команде `/prices`
- Отображение цены, изменения за 24 часа, рыночной капитализации и объема торгов

## Требования

- Go 1.25.4 или выше
- Telegram бот токен (получить у [@BotFather](https://t.me/BotFather))
- CoinMarketCap API ключ (получить на [coinmarketcap.com/api](https://coinmarketcap.com/api/))

## Установка

1. Клонируйте репозиторий или перейдите в директорию проекта:
```bash
cd bot
```

2. Установите зависимости:
```bash
go mod download
```

3. Создайте файл `.env` на основе `.env.example`:
```bash
cp .env.example .env
```

4. Заполните переменные окружения в файле `.env`:
   - `TELEGRAM_BOT_TOKEN` - токен вашего Telegram бота
   - `COINMARKETCAP_API_KEY` - API ключ CoinMarketCap

## Запуск

### Вариант 1: Использование переменных окружения напрямую

```bash
export TELEGRAM_BOT_TOKEN="your_token_here"
export COINMARKETCAP_API_KEY="your_api_key_here"
go run .
```

### Вариант 2: Использование .env файла (требует установки go-dotenv)

Если вы хотите использовать `.env` файл, установите пакет:
```bash
go get github.com/joho/godotenv
```

И добавьте в начало `main.go`:
```go
import _ "github.com/joho/godotenv/autoload"
```

### Вариант 3: Запуск с Docker Compose (рекомендуется)

1. Создайте файл `.env` в корне проекта:
```bash
# Telegram Bot
TELEGRAM_BOT_TOKEN=your_token_here
COINMARKETCAP_API_KEY=your_api_key_here

# PostgreSQL Database
POSTGRES_DB=
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_PORT=

# Database Connection String (опционально, можно использовать для подключения из приложения)
DATABASE_URL=postgres://cryptobot:changeme@postgres:5432/cryptobot?sslmode=disable
```

2. Запустите контейнеры (PostgreSQL и бот):
```bash
docker-compose up -d
```

3. Просмотр логов:
```bash
# Все сервисы
docker-compose logs -f

# Только бот
docker-compose logs -f crypto-bot

# Только PostgreSQL
docker-compose logs -f postgres
```

4. Остановка контейнеров:
```bash
docker-compose down
```

5. Остановка с удалением volumes (удалит данные БД):
```bash
docker-compose down -v
```

6. Перезапуск после изменений:
```bash
docker-compose up -d --build
```

**Примечание:** Данные PostgreSQL сохраняются в Docker volume `postgres_data` и не будут удалены при остановке контейнеров (если не использовать флаг `-v`).

**Подключение к PostgreSQL:**
```bash
# Из хоста (если порт проброшен)
psql -h localhost -p 5432 -U cryptobot -d cryptobot

# Из другого контейнера в той же сети
psql -h postgres -U cryptobot -d cryptobot

# Пароль можно указать через переменную окружения
PGPASSWORD=changeme psql -h localhost -p 5432 -U cryptobot -d cryptobot
```

**Просмотр данных БД:**
```bash
# Список volumes
docker volume ls

# Информация о volume
docker volume inspect cryptonotifications_postgres_data

# Резервное копирование БД
docker-compose exec postgres pg_dump -U cryptobot cryptobot > backup.sql

# Восстановление БД
docker-compose exec -T postgres psql -U cryptobot cryptobot < backup.sql
```

### Вариант 4: Запуск с Docker напрямую

1. Соберите Docker образ:
```bash
docker build -t crypto-bot .
```

2. Запустите контейнер с переменными окружения:
```bash
docker run -d \
  --name crypto-bot \
  -e TELEGRAM_BOT_TOKEN="your_token_here" \
  -e COINMARKETCAP_API_KEY="your_api_key_here" \
  --restart unless-stopped \
  crypto-bot
```

Или используйте файл с переменными окружения:
```bash
docker run -d \
  --name crypto-bot \
  --env-file .env \
  --restart unless-stopped \
  crypto-bot
```

3. Просмотр логов:
```bash
docker logs -f crypto-bot
```

4. Остановка контейнера:
```bash
docker stop crypto-bot
docker rm crypto-bot
```

## Использование

После запуска бота, отправьте ему команду `/start` в Telegram.

### Доступные команды:

- `/start` - начать работу с ботом
- `/help` - показать справку по командам
- `/price <символ>` - получить котировку одной криптовалюты
  - Пример: `/price BTC`
- `/prices <символ1,символ2,...>` - получить котировки нескольких криптовалют
  - Пример: `/prices BTC,ETH,BNB`

### Примеры использования:

```
/price BTC
/prices BTC,ETH,BNB,ADA
```

## Структура проекта

```
CryptoNotifications/
├── Dockerfile           # Dockerfile для сборки образа
├── docker-compose.yml   # Конфигурация Docker Compose
├── .dockerignore        # Игнорируемые файлы для Docker
├── .env                 # Файл с переменными окружения (создать самостоятельно)
├── init.sql             # SQL скрипт для инициализации БД (опционально)
└── bot/
    ├── main.go              # Основной файл бота
    ├── coinmarketcap.go     # Модуль для работы с CoinMarketCap API
    ├── go.mod               # Файл зависимостей Go
    ├── go.sum               # Файл контрольных сумм зависимостей
    └── README.md            # Этот файл
```

## Получение API ключей

### Telegram Bot Token

1. Откройте Telegram и найдите [@BotFather](https://t.me/BotFather)
2. Отправьте команду `/newbot`
3. Следуйте инструкциям для создания бота
4. Скопируйте полученный токен

### CoinMarketCap API Key

1. Перейдите на [coinmarketcap.com/api](https://coinmarketcap.com/api/)
2. Зарегистрируйтесь или войдите в аккаунт
3. Перейдите в раздел API Keys
4. Создайте новый API ключ (Basic план бесплатный)
5. Скопируйте API ключ

## Лицензия

MIT
