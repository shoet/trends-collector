.PHONY: build clean deploy
.DEFAULT_GOAL := help

build:
	env GOOS=linux go build -ldflags="-s -w" -o cmd/bin/main cmd/main.go

build-local:
	env go build -ldflags="-s -w" -o cmd/bin/main cmd/main.go

clean:
	rm -rf ./bin

deploy: clean build
	# TODO

help: ## Show options
	@grep -E '^[a-zA-Z_]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
