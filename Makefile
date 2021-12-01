postgres:
	docker-compose up --build postgres
redis:
	docker-compose up --build redis

stop:
	docker-compose down
kek:
	docker-compose up --build

.PHONY: postgres redis stop kek
