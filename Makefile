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
