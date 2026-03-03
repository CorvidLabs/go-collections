.PHONY: build test lint vet fmt clean bench

build:
	go build ./...

test:
	go test -v -race -count=1 ./...

bench:
	go test -bench=. -benchmem ./...

lint: vet fmt

vet:
	go vet ./...

fmt:
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed on:" && gofmt -l . && exit 1)

clean:
	go clean -testcache
