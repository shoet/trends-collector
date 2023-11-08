.PHONY: build clean deploy
.DEFAULT_GOAL := help

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/health functions/health/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/topic functions/topic/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/collect functions/collect_trends/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

help: ## Show options
	@grep -E '^[a-zA-Z_]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
