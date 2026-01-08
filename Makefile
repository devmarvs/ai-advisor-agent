.PHONY: run-api migrate
run-api:
	cd server && go run .
migrate:
	@DB_URL=$${DATABASE_URL}; \
	if [ -z "$$DB_URL" ]; then \
		DB_URL="postgresql://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$${DB_PORT:-5432}/$$DB_NAME?sslmode=$${DB_SSLMODE:-require}$${DB_CHANNEL_BINDING:+&channel_binding=$$DB_CHANNEL_BINDING}"; \
	fi; \
	psql "$$DB_URL" -f api/migrations/0001_init.sql && \
	psql "$$DB_URL" -f api/migrations/0002_indexes.sql && \
	psql "$$DB_URL" -f api/migrations/0003_task_queue_pg.sql
