# postgres
make_postgres:
	docker run --name my-pg \
	-e POSTGRES_USER=ROOT \
	-e POSTGRES_PASSWORD=123456 \
	-e POSTGRES_DB=MAN \
	-p 5432:5432 \
	-d postgres

# redis
make_redis:
	docker run --name=my-rd \
	-p 6379:6379 \
	-d redis

databaseURL="postgresql://ROOT:123456@localhost:5432/MAN?sslmode=disable"

migrate_up:
	migrate -path="database/migrate" -database=${databaseURL} up

migrate_drop:
	migrate -path="database/migrate" -database=${databaseURL} drop -f


make_db: make_postgres make_redis

start_pg:
	docker start my-pg

start_rd:
	docker start my-rd

start_db: start_pg start_rd