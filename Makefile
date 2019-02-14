

.PHONY: all
all: setuppostgres

.PHONY: setuppostgres
setuppostgres:
	docker run --name gosupport_db -v "`pwd`/db-data/postgres:/var/lib/postgresql/data" --env-file .env postgres

.PHONY: postgres
postgres:
	docker run -p "5433:5432" postgres
