# Variables
BINARY_NAME=redis-vs-nats
BUILD_DIR=bin
MAIN_FILE=main.go
BUILDTIMESTAMP=$(shell date -u +%Y%m%d%H%M%S)

# append / to this if set
DOCKER_REGISTRY=
DOCKER_REPO=phillezi
DOCKER_IMAGE=$(BINARY_NAME)
DOCKER_TAG=latest

# Targets
.PHONY: all clean build run dcoker/local docker/push compose/bench compose/bench/% lint

all: build

build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 go build -ldflags "-X main.buildTimestamp=$(BUILDTIMESTAMP)" -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete."

run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

docker/local:
	@echo "Building docker image..."
	@docker buildx build -t $(DOCKER_REGISTRY)$(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) --load

docker/push:
	@echo "Building and pushing docker image..."
	@docker buildx build -t $(DOCKER_REGISTRY)$(DOCKER_REPO)/$(DOCKER_IMAGE):$(DOCKER_TAG) --push

compose/bench/%:
	@echo "🔹 Running Benchmark [$*]..."
	# todo make this neater
	@docker compose --file $*.docker-compose.yml up --force-recreate --build --abort-on-container-exit 2>&1 | grep --line-buffered "benchmark_runner"
	@docker compose --file $*.docker-compose.yml down --remove-orphans

compose/bench: compose/bench/redis compose/bench/nats compose/bench/mono
	@echo "Benchmarking completed for Redis, NATS and a monolithic channel based impl.

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

lint:
	@./scripts/util/check-lint.sh