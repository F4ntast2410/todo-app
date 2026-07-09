CREATE TABLE user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Уникальный токен сессии
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
-- Индекс для быстрого поиска сессии при каждом запросе
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON user_sessions(expires_at);