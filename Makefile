

.PHONY: all
all: setuppostgres

.PHONY: setuppostgres
setuppostgres:
	docker run --name gosupport_db --env-file .env postgres

.PHONY: postgres
postgres:
	docker run gosupport_db
