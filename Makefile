.PHONY: env up up-aiven down build logs restart ps health mysql-cli

env:
	@test -f .env || cp .env.example .env
	@echo ".env ready — review and update secrets before production use"

up: env
	docker compose up --build -d

# Backend + nginx only, DB from .env (e.g. Aiven)
up-aiven: env
	docker compose -f docker-compose.aiven.yml up --build -d

mysql-cli: env
	@set -a && . ./.env && set +a && \
	docker compose exec mysql mysql -u"$$MYSQL_USER" -p"$$MYSQL_PASSWORD" "$$MYSQL_DATABASE"

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

restart:
	docker compose restart

ps:
	docker compose ps

health:
	@curl -sf http://127.0.0.1:3131/api/v1/health && echo
