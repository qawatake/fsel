BINDIR := $(CURDIR)/bin
golangci-lint := go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
gofmtmd := go run github.com/po3rin/gofmtmd/cmd/gofmtmd@latest

fmtmd:
	$(gofmtmd) -r README.md

test:
	go test ./... -shuffle=on -race

lint:
	go vet  ./...
	$(golangci-lint) run

test.cover:
	go test -race -shuffle=on -coverprofile=coverage.txt -covermode=atomic ./...

# it is expected to fail
test.fail:
	go test -tags falsepositive ./internal/...

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
