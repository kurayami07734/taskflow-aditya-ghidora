.PHONY: dev build run test

export DB_HOST := localhost

dev:
	air

build:
	go build -o ./tmp/main .

run:
	go run .

test:
	go test ./tests/... -v -p=4 -count=1
