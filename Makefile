.PHONY: env up down build logs

env:
	@test -f .env || cp .env.example .env
	@echo ".env ready — review and update secrets before production use"

up: env
	docker compose up --build -d

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
