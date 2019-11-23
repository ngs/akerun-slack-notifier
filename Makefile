build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/callback callback/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/check check/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/redirect redirect/main.go

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
