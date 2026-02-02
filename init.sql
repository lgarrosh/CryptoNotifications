-- Инициализационный SQL скрипт для PostgreSQL
-- Этот файл выполняется автоматически при первом запуске контейнера PostgreSQL

-- Создание расширений (если нужны)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Пример создания таблицы для истории запросов (раскомментируйте при необходимости)
-- CREATE TABLE IF NOT EXISTS price_requests (
--     id SERIAL PRIMARY KEY,
--     user_id BIGINT NOT NULL,
--     username VARCHAR(255),
--     symbols TEXT NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- Создание индексов для оптимизации запросов
-- CREATE INDEX IF NOT EXISTS idx_price_requests_user_id ON price_requests(user_id);
-- CREATE INDEX IF NOT EXISTS idx_price_requests_created_at ON price_requests(created_at);
