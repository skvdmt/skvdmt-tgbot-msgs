CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL DEFAULT uuidv7(),
    telegram_user_id INTEGER NOT NULL,
    message_created_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id),
    UNIQUE(telegram_user_id)
);

COMMENT ON TABLE users IS 'Таблица telegram пользователей';
COMMENT ON COLUMN users.id IS 'Идентификатор пользователя';
COMMENT ON COLUMN users.telegram_user_id IS 'Идентификатор telegram пользователя';
COMMENT ON COLUMN users.message_created_at IS ' пользователя';
COMMENT ON COLUMN users.created_at IS 'дата и время создания записи';
COMMENT ON COLUMN users.updated_at IS 'дата и время последнего обновления записи';

CREATE TABLE IF NOT EXISTS messages (
    id UUID NOT NULL DEFAULT uuidv7(),
    user_id INTEGER NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

COMMENT ON TABLE messages IS 'Таблица сообщений отправленных в telegram бот';
COMMENT ON COLUMN messages.id IS 'Идентификатор сообщения';
COMMENT ON COLUMN messages.user_id IS 'Идентификатор telegram пользователя';
COMMENT ON COLUMN messages.message IS 'Текст сообщения';
COMMENT ON COLUMN messages.created_at IS 'дата и время создания записи';
COMMENT ON COLUMN messages.updated_at IS 'дата и время последнего обновления записи';
