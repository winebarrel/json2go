.PHONY: all
all: vet test build

.PHONY: build
build:
	go build ./cmd/json2go

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: gen
gen:
	go generate

.PHONY: redash-setup
redash-setup:
	psql -U postgres -h localhost -p 15432 -f etc/redash.sql
	$(MAKE) redash-upgrade-db

.PHONY: redash-upgrade-db
redash-upgrade-db:
	docker compose run --rm server manage db upgrade

.PHONY: redash-create-db
redash-create-db:
	docker compose run --rm server create_db

.PHONY: pg-dump
pg-dump:
	pg_dump -U postgres -h localhost -p 15432  -c --if-exists  > etc/redash.sql
