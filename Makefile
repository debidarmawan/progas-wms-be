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
	docker compose logs -f app
