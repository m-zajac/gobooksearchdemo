.PHONY: build test clean lint lint-more image generate

build: generate
	go build ./cmd/cli
	go build ./cmd/server

test:
	go test -race ./...

clean:
	rm -f ./cli
	rm -f ./server
	rm -rf ./data

lint: $(shell go env GOPATH)/bin/golint
	@$(shell go env GOPATH)/bin/golint -set_exit_status `go list ./... | grep -v /vendor/`

lint-more: $(shell go env GOPATH)/bin/golangci-lint
	@$(shell go env GOPATH)/bin/golangci-lint run ./...

image:
	ln -sf ./build/Dockerfile .
	docker build -t gobooksearchdemo .
	rm ./Dockerfile

generate: $(shell go env GOPATH)/bin/swag
	GOFLAGS="-mod=readonly" go generate ./...

$(shell go env GOPATH)/bin/golint:
	@GOFLAGS="-mod=readonly" go get golang.org/x/lint/golint

$(shell go env GOPATH)/bin/golangci-lint:
	@GOFLAGS="-mod=readonly" go get github.com/golangci/golangci-lint/cmd/golangci-lint

$(shell go env GOPATH)/bin/swag:
	@GOFLAGS="-mod=readonly" go get github.com/swaggo/swag/cmd/swag@v1.6.6-0.20200323071853-8e21f4cefeea
