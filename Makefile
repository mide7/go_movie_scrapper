# Define the output binary name
BINARY_NAME := movie_scrapper

# Load environment variables from .env file
.PHONY: load-env
load-env:
	@export $$(grep -v '^#' ./config/.env | xargs) && \
	 echo "Loaded environment variables" && \
	 echo "MONGO_URI: $$MONGO_URI"


# Default target: Build the project
.PHONY: all
all: build

# Build the Go project
.PHONY: build
build:
	go build -o $(BINARY_NAME) ./cmd/movie_scrapper

# Run the Go project
.PHONY: run
run: build
	./$(BINARY_NAME)

# Run the Go project in development mode
.PHONY: run-d
run-d:
	go run ./cmd/movie_scrapper

# Test the Go project
.PHONY: test
test:
	go test ./...

# Clean the build files
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

# Update dependencies
.PHONY: deps
deps:
	go mod tidy

# Format the code
.PHONY: fmt
fmt:
	go fmt ./...

