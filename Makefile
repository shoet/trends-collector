.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/health handlers/health/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/topic handlers/topic/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/health functions/health/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/topic functions/topic/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
