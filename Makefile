.PHONY: run-api run-web migrate

run-api:
	@cd api && go run .

run-web:
	@cd web && npm run dev

migrate:
	@psql $$DATABASE_URL -f infra/migrations/0001_init.sql && \
		psql $$DATABASE_URL -f infra/migrations/0002_indexes.sql
