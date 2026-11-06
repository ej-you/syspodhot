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

## Настройка автообновления SSL-сертификата

Переходим в контейнер с веб-сервером:

```shell
docker compose exec -it web-server sh
```

Раскомментируем `#include /etc/nginx/conf.d/https-site.conf;` в файле `/etc/nginx/nginx.conf`:

```shell
vi /etc/nginx/nginx.conf
```

Проверим конфиг nginx и перечитаем его:

```shell
nginx -t
nginx -s reload
```

Выйдем из контейнера `web-server`:

```shell
exit
```

Удалим контейнер `certbot`:

```shell
docker compose rm certbot
```

Для автообновления сертификата настроим задачу в `cron` (раз в месяц):

```shell
0 0 1 */1 * docker-compose -f /full/path/to/docker-compose.yml run --rm certbot renew --force-renewal
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
