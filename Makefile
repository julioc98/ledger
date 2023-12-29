.PHONY: run migrate-up docker-run migrate-tests test


docker-run:
	docker-compose up -d --build

migrate-up: docker-run
	docker run -v ./gateway/pg/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up

run: docker-run migrate-up
	go run ./app/service/main.go

migrate-tests: docker-run
	docker run -v ./gateway/pg/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database 'postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable' up

test: docker-run migrate-tests 
	go test -v ./...
	

