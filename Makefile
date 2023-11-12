.DEFAULT_GOAL := help

DOCKER_IMAGE := trends-collector-crawler
DOCKER_TAG := latest

.PHONY: build
build: ## Build Lambda functions binary
	env GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/health functions/health/main.go
	env GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/topic functions/topic/main.go

.PHONY: clean
clean: ## Clean Lambda functions binary
	rm -rf ./bin

.PHONY: deploy
deploy: clean build ## Deploy by Serverless Framework
	sls deploy --verbose

.PHONY: build-crawler
build-crawler: ## Build crawler binary
	cd crawler && \
		env GOOS=linux go build -trimpath -ldflags="-s -w" -o crawler/cmd/bin/main cmd/crawltask/main.go

.PHONY: build-crawler-local
build-crawler-local: ## Build crawler binary on Arm64
	cd crawler && \
		go build -trimpath -ldflags="-s -w" -o crawler/cmd/bin/main cmd/crawltask/main.go

.PHONY: build-container-crawler
build-container-crawler: ## Build crawler container image
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} \
		--platform linux/amd64 \
		--target deploy \
		-f crawler/Dockerfile \
		.

.PHONY: build-container-crawler-local
build-container-crawler-local: ## Build crawler container image on Arm64
	docker build -t ${DOCKER_IMAGE}:${DOCKER_TAG} \
		--target deploy \
		-f crawler/Dockerfile \
		--no-cache \
		.

.PHONY: push-container-crawler
push-container-crawler: ## Push crawler container image
	bash ./container_push.sh

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
