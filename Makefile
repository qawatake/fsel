BINDIR := $(CURDIR)/bin

test:
	go mod tidy
	go test ./... -shuffle=on -race

lint:
	go mod tidy
	go vet  ./...

test.cover:
	go mod tidy
	go test -race -shuffle=on -coverprofile=coverage.txt -covermode=atomic ./...

# For local environment
cov:
	go test -cover -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

build:
	@go build -o $(BINDIR)/fsel ./cmd/fsel

test.vet:
	-@go vet -vettool=$(BINDIR)/fsel ./internal/...


# Run the example
run.example:
	@$(MAKE) build
	@$(MAKE) test.vet
