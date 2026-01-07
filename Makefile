.PHONY: run-api migrate
run-api:
	cd server && go run .
migrate:
	psql $$DATABASE_URL -f infra/migrations/0001_init.sql && \
	psql $$DATABASE_URL -f infra/migrations/0002_indexes.sql && \
	psql $$DATABASE_URL -f infra/migrations/0003_task_queue_pg.sql
