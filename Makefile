.PHONY: test.cleancache
test.cleancache:
	go clean -testcache

.PHONY: test.cover
test.cover: test.cleancache
	go test -v -race ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func coverage.out

.PHONY: migration
migration:
	migrate create -ext sql -dir db/migrations/$(module) $(name)

.PHONY: migrate
migrate:
	migrate -path db/migrations -database "$(url)?sslmode=disable" -verbose up
