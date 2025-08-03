# Subscriptions API

REST API для управления онлайн-подписками пользователей, реализованный на Go с использованием Gin, PostgreSQL, Swagger и Docker.

## Описание

- CRUDL-операции для подписок (создание, чтение, обновление, удаление, список)
- Подсчёт суммарной стоимости подписок за выбранный период
- Swagger-документация
- Конфигурация через .yaml

## Технологии

- Go + Gin
- PostgreSQL
- Goose (миграции)
- Swagger (документация)
- Logrus (логирование)
- Docker + Docker Compose
<!-- - GoMock + mockgen (моки в тестах)
- GitHub Actions (CI: тесты, линтер и сборка) -->

## Запуск контейнера

```bash
docker compose -f build/docker-compose.yml up -d --build
```


Swagger-документация доступна по адресу

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

<!-- #### `make build`

Собирает приложение в bin/marketplace.

#### `make test`

Запускает тесты.

#### `make lint`

Запускает линтер. -->