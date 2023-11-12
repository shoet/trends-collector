.DEFAULT_GOAL := help

DOCKER_TAG := latest

.PHONY: build
build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/health functions/health/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/topic functions/topic/main.go

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: deploy
deploy: clean build
	sls deploy --verbose

.PHONY: build-crawler
build-crawler:
	env GOOS=linux go build -ldflags="-s -w" -o crawler/cmd/bin/main crawler/cmd/main.go

.PHONY: build-crawler-local
build-crawler-local:
	go build -ldflags="-s -w" -o crawler/cmd/bin/main crawler/cmd/main.go

.PHONY: build-container-crawler
build-container-crawler:
	docker build -t trands-collector-crawler:${DOCKER_TAG} \
		--platform linux/amd64 \
		--target deploy \
		-f crawler/Dockerfile \
		.

.PHONY: build-container-crawler-local
build-container-crawler-local:
	docker build -t trands-collector-crawler:${DOCKER_TAG} \
		--target deploy \
		-f crawler/Dockerfile \
		--no-cache \
		.

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
