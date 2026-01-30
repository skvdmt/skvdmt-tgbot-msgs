CREATE TABLE IF NOT EXISTS messages (
    id UUID NOT NULL DEFAULT uuidv7(),
    telegram_user_id BIGINT NOT NULL,
    telegram_user_name VARCHAR(32) DEFAULT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

COMMENT ON TABLE messages IS 'Таблица сообщений отправленных в telegram бот';
COMMENT ON COLUMN messages.id IS 'Идентификатор сообщения';
COMMENT ON COLUMN messages.telegram_user_id IS 'Идентификатор telegram пользователя';
COMMENT ON COLUMN messages.telegram_user_name IS 'Имя telegram пользователя';
COMMENT ON COLUMN messages.message IS 'Текст сообщения';
COMMENT ON COLUMN messages.created_at IS 'дата и время создания записи';
