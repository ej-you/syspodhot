# Сайт-визитка для музыкальной группы Системный Подход

## Конфигурация

### Сервер

Файл конфигурации приложения - [config.yml](./config.yml)

## Запуск и настройка

После настройки можно запустить приложение. Для этого используется `docker compose`:

```shell
docker compose up -d --build
```

Проверить запущенный сервер можно командой:

```shell
docker compose ps
```

Чтобы опустить сервер, используйте:

```shell
docker compose down
```

## Администрирование

Для управления пользователями используется менеджер:

```shell
docker exec -it syspodhot_server /bin/sh -c "/app/main manager"
```

После изменения параметров [config.yml](./config.yml) необходимо перезапустить сервер:

```shell
    docker compose restart server
```
