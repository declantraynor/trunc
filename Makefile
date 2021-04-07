.PHONY: run stop test

run:
	@docker compose up --build -d web

stop:
	@docker compose down

test:
	@docker compose build test
	@docker compose run --rm test