.PHONY: build run migrate seed clean docker-migrate docker-seed all help
# Docker compose command
DOCKER_COMPOSE=docker-compose

run-dev:
	$(DOCKER_COMPOSE) up

migrate:
	@echo "Running migrations in docker..."
	$(DOCKER_COMPOSE) run db_migrate

seed:
	@echo "Seeding database in docker..."
	$(DOCKER_COMPOSE) run db_seed

docker-down: 
	@echo "Stopping docker containers..."
	$(DOCKER_COMPOSE) down

docker-db-status:
	@echo "Checking database status..."
	docker exec -it postgres_quiz_db psql -U postgres -d quizdb -c "\
	SELECT 'Quizzes' as table_name, COUNT(*) as count FROM quizzes UNION ALL \
	SELECT 'Questions', COUNT(*) FROM questions UNION ALL \
	SELECT 'Answers', COUNT(*) FROM answers;"

# Full setup: Start, migrate, and seed
docker-setup: docker-up docker-migrate docker-seed
	@echo "Application setup complete!"
