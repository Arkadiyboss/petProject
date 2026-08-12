#!/bin/sh

# Запускаем PostgreSQL
/usr/local/bin/docker-entrypoint.sh postgres &

# Ждём БД
sleep 5

# Запускаем приложение
/usr/local/bin/myapp