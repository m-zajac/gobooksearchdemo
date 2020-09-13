.PHONY: build test clean image loadtest

build:
	go build ./cmd/cli

test:
	go test -race ./...

clean:
	rm -f ./cli
	rm -f ./bookdata.db


lint: $(shell go env GOPATH)/bin/golint
	@$(shell go env GOPATH)/bin/golint -set_exit_status `go list ./... | grep -v /vendor/`

lint-more: $(shell go env GOPATH)/bin/golangci-lint
	@$(shell go env GOPATH)/bin/golangci-lint run ./...

# image:
# 	ln -sf ./build/Dockerfile .
# 	docker build -t goprojectdemo .
# 	rm ./Dockerfile

# loadtest:
# 	wrk --latency -d 15m -t 2 -c 15 -s scripts/loadtest.lua http://localhost:8080

$(shell go env GOPATH)/bin/golint:
	@GOFLAGS="-mod=readonly" go get golang.org/x/lint/golint

$(shell go env GOPATH)/bin/golangci-lint:
	@GOFLAGS="-mod=readonly" go get github.com/golangci/golangci-lint/cmd/golangci-lint