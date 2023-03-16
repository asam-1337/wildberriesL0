PHONY: test
test:
	go test ./...

.PHONY: run
run:
	go run ./cmd/app/main.go

.PHONY: publish
publish:
	go run cmd/publish/publish.go

PHONY: goose-up
goose-up:
	goose -dir ./migrations postgres "user=postgres dbname=postgres password=qwerty host=127.0.0.1  port=5432" up

PHONY: goose-reset
goose-reset:
	goose -dir ./migrations postgres "user=postgres dbname=postgres password=qwerty host=127.0.0.1 port=5432" reset

PHONY: docker-compose-up
docker-compose-up:
	docker compose up