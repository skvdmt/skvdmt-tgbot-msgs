# skvdmt-tgbot-msgs

### Переводы:
[English](./README.md)

# Описание
Телеграм бот собирающий сообщения от пользователей.

Пример: [t.me/skidanovdima_msgs_bot](https://t.me/skidanovdima_msgs_bot).

Создан без сторонних библиотек по работе с [Telegram Bot API](https://core.telegram.org/bots/api).

Бот проводит авторизацию путем анализа введенных символов с картинки. Для генерирования каптчи используется библиотека [skvdmt/captcha](https://github.com/skvdmt/captcha) и дополнительные шрифты, находящиеся в папке [fonts](./fonts)

# Скачать исходный код
```sh
git clone https://github.com/skvdmt/skvdmt-tgbot-msgs .
```
# Скачать Docker image
```sh
docker pull skvdmt/skvdmt-tgbot-msgs
```

## Требования переменных окружения:

для приложения:
- MODE — Режим работы (доступные значения: dev | prod);
- TGBOT_TOKEN — Токен авторизации телеграм бота;
- DB_PASSWORD — Пароль от пользователя postgres;

для postgres:
- POSTGRES_USER — Имя пользователя postgres;
- POSTGRES_PASSWORD — Пароль от пользователя postgres;
- POSTGRES_DB — Имя базы данных postgres;

## Инструкции
Минимальная настройка окружения в [docker-compose.yaml](docker-compose.yaml)

Инструкции по сборке образа в [Dockerfile](Dockerfile)

Настройки приложения в [файле конфигурации](./config/config.yaml)

## Ссылки:
- [Docker image](https://hub.docker.com/r/skvdmt/skvdmt-tgbot-msgs) — Собраный образ приложения на docker hub.
- [Example](https://t.me/skidanovdima_msgs_bot) — Пример работающего телеграм бота.
- [Исходник](https://github.com/skvdmt/skvdmt-msgs-front/) — Фронтенд приложение для отображения сообщений, отправленных боту.
- [Исходник](https://github.com/skvdmt/skvdmt-msgs-back/) — Бекенд API приложение для отображения сообщений, отправленных боту.
- [Author](https://skvdmt.ru) — Skidanov Dmitry.
