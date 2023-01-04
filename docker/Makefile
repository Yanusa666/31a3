#!make

.PHONY: start-containers
start-containers: ## Создает необходимые контейнеры
	docker-compose --project-name crud -f ./docker-compose.yml down
	docker-compose --project-name crud -f ./docker-compose.yml pull
	docker-compose --project-name crud -f ./docker-compose.yml up -d --build

.PHONY: stop-containers
stop-containers: ## Останавливает контейнеры
	docker-compose --project-name crud  -f ./docker-compose.yml down