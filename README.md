# skvdmt-tgbot-msgs

### Translations:
[Русский](./README_ru.md)

# Description
Telegram bot collecting messages from users.

Example: [t.me/skidanovdima_msgs_bot](https://t.me/skidanovdima_msgs_bot).

Created without third-party libraries for working with [Telegram Bot API](https://core.telegram.org/bots/api).


The bot performs authorization by analyzing the characters entered from the image. The [skvdmt/captcha](https://github.com/skvdmt/captcha) library is used to generate captchas and additional fonts located in the [fonts](./fonts) folder

# Download sources
```sh
git clone https://github.com/skvdmt/skvdmt-tgbot-msgs .
```

# Download docker image
```sh
docker pull skvdmt/skvdmt-tgbot-msgs
```

# Run in docker container
```sh
docker run -d \
  --name msgs \
  --env MODE=prod \
  --env TGBOT_TOKEN=bot_token \
  --env DB_PASSWORD=postgres_password \
  --restart unless-stopped \
  skvdmt/skvdmt-tgbot-msgs
```

## Environment variable requirements:

application:
- MODE — Operating mode (available values: dev | prod);
- TGBOT_TOKEN — Telegram bot authorization token;
- DB_PASSWORD — Password for the postgres user;

postgres:
- POSTGRES_USER — Username postgres;
- POSTGRES_PASSWORD — Password for the postgres user;
- POSTGRES_DB — Postgres database name;

## Instructions
Minimal environment setup in [docker-compose.yaml](docker-compose.yaml)

Instructions for building an image in [Dockerfile](Dockerfile)

Application settings in [config file](./config/config.yaml)

## Links:
- [Docker image](https://hub.docker.com/r/skvdmt/skvdmt-tgbot-msgs) — Docker image on docker hub.
- [Example](https://t.me/skidanovdima_msgs_bot) An example of a working Telegram bot.
- [Sources](https://github.com/skvdmt/skvdmt-msgs-front/) — Frontend application for displaying messages sent to the bot.
- [Sources](https://github.com/skvdmt/skvdmt-msgs-back/) — Backend API application for displaying messages sent to the bot.
- [Author](https://skvdmt.ru) — Skidanov Dmitry.
