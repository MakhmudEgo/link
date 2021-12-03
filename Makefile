
runPostgres: postgres
	./db.sh postgres
	docker-compose up --build link

runRedis: redis
	./db.sh redis
	docker-compose up --build link

postgres:
	docker-compose up --build -d postgres

redis:
	docker-compose up --build -d redis

stop:
	docker-compose down

.PHONY: postgres redis stop runPostgres runRedis
