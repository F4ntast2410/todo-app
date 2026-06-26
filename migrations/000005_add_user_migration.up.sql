-- 1. Удаляем старые таблицы каскадно (связи в tasks отвалятся сами)
DROP TABLE IF EXISTS users CASCADE;

-- 2. На всякий случай подчищаем колонку в tasks, если она осталась
ALTER TABLE tasks DROP COLUMN IF EXISTS user_id;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT,          -- Отображаемое имя в твоем приложении
    created_at TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE tasks ADD COLUMN IF NOT EXISTS user_id INT;

-- 3. Накатываем CONSTRAINT (внешний ключ), который связывает tasks и users
ALTER TABLE tasks 
ADD CONSTRAINT fk_tasks_user 
FOREIGN KEY (user_id) 
REFERENCES users(id) 
ON DELETE CASCADE;

CREATE TABLE user_passwords (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL -- Хэш пароля (никогда не храним в чистом виде!)
);

CREATE TABLE user_telegram (
    tg_id BIGINT PRIMARY KEY, -- ID пользователя внутри Telegram
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    username TEXT             -- TG-юзернейм для удобства
);