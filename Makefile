.PHONY : test test-lexer

all: test build

build:
	go build ./...

run:
	go run ./... $(FILE)

update-snapshots:
	UPDATE_SNAPS=true make test

test:
	go test -coverpkg=./... -coverprofile=coverage.txt ./...

coverage: test
	go tool cover -html=coverage.cov