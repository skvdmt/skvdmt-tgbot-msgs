CREATE TABLE IF NOT EXISTS messages (
    id UUID NOT NULL DEFAULT uuidv7(),
    telegram_user_id BIGINT NOT NULL,
    telegram_user_name VARCHAR(32) DEFAULT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

COMMENT ON TABLE messages IS 'Таблица сообщений отправленных в telegram бот';
COMMENT ON COLUMN messages.id IS 'Идентификатор сообщения';
COMMENT ON COLUMN messages.telegram_user_id IS 'Идентификатор telegram пользователя';
COMMENT ON COLUMN messages.telegram_user_name IS 'Имя telegram пользователя';
COMMENT ON COLUMN messages.text IS 'Текст сообщения';
COMMENT ON COLUMN messages.created_at IS 'дата и время создания записи';

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 1');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 2');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 3');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 4');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 5');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 6');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 7');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 8');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 9');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 10');

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 11');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 12');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 13');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 14');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 15');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 16');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 17');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 18');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 19');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 20');

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 21');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 22');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 23');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 24');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 25');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 26');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 27');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 28');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 29');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 30');

INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 31');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 32');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 33');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 34');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 35');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 36');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 37');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 38');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 39');
INSERT INTO messages(telegram_user_id, telegram_user_name, text) VALUES(111, 'username', 'Hello 40');
