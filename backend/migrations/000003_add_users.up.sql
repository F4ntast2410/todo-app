-- 1. Создаем таблицу пользователей
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY, -- Сюда будем писать Telegram ID (он большой, поэтому BIGINT)
    username VARCHAR(255) DEFAULT '',
    first_name VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW()
);

-- 2. Добавляем колонку user_id в существующую таблицу tasks
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS user_id BIGINT;

-- 3. Накатываем CONSTRAINT (внешний ключ), который связывает tasks и users
ALTER TABLE tasks 
ADD CONSTRAINT fk_tasks_user 
FOREIGN KEY (user_id) 
REFERENCES users(id) 
ON DELETE CASCADE;