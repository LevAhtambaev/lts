# LTS - Leo`s Travel Stories

## Запуск сервиса

### Запуск локально:
```
docker compose up -d db
make migrate
make run
```

### Запуск в docker контейнере:
```
docker compose build
docker compose up
```