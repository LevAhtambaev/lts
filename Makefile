PWD = ${CURDIR}
# Название сервиса
SERVICE_NAME = lts

# компиляция сервиса
.PHONY: build
build:
	go build -o bin/$(SERVICE_NAME)  $(PWD)/cmd/$(SERVICE_NAME)  -config ./configs/config.yaml

# Запуск сервиса
.PHONY: run
run:
	go run $(PWD)/cmd/$(SERVICE_NAME)

# Запуск миграций
.PHONY: migrate
migrate:
	go run $(PWD)/cmd/migrate

# Генерация сваггера
.PHONY: swagger
swagger:
	swag init -g /cmd/lts/main.go