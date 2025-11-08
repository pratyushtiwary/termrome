.PHONY : test test-lexer

all: test build

build:
	go build ./...

run:
	go run ./... $(FILE)

update-snapshots:
	UPDATE_SNAPS=true make test

test:
	go test -coverpkg=./... -coverprofile=coverage.cov ./...

coverage: test
	go tool cover -html=coverage.cov