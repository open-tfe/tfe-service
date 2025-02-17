BINARY_NAME=tfe-service
# Go variables
GOCMD=go
GOBUILD=$(GOCMD) build
# Build directories
BUILD_DIR=build
LINUX_AMD64_DIR=$(BUILD_DIR)/linux_amd64
DARWIN_AMD64_DIR=$(BUILD_DIR)/darwin_amd64
DARWIN_ARM64_DIR=$(BUILD_DIR)/darwin_arm64

.PHONY: build run test swagger migrate

# Build for Linux (amd64)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(LINUX_AMD64_DIR)/$(BINARY_NAME) ./cmd/server

# Build for macOS (both amd64 and arm64)
build-darwin:
	# CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DARWIN_AMD64_DIR)/$(BINARY_NAME) ./cmd/server
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DARWIN_ARM64_DIR)/$(BINARY_NAME) ./cmd/server


test:
	$(GOCMD) test -v ./...

swagger:
	swag init -g cmd/server/main.go -o swagger

lint:
	golangci-lint run

deps:
	$(GOCMD) mod download

clean:
	rm -rf bin/ 

# Run go fmt against code
fmt:
	$(GOCMD) mod tidy
	$(GOCMD) fmt ./...

start-db:
	docker run --rm -d --name tf-test -p 3306:3306 -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=tfe -p 5432:5432 postgres:16

stop-db:
	docker stop tf-test

migrate:
	atlas schema apply \
	-u "postgres://postgres:pass@localhost:5432/tfe?sslmode=disable" \
	--to "file://migrations/" \
	--dev-url "docker://postgres/16/dev" \
	--auto-approve

run:
	$(GOCMD) run cmd/server/main.go
