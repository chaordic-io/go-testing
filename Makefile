.PHONY: all
all: proto-lint proto-generate test lint

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go clean -testcache
	go test ./... -race -covermode=atomic -coverprofile=coverage.out

.PHONY: proto-lint
proto-lint: 
	buf lint

.PHONY: proto-install
proto-install: 
	go install \
    	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    	google.golang.org/protobuf/cmd/protoc-gen-go \
    	google.golang.org/grpc/cmd/protoc-gen-go-grpc

.PHONY: proto-generate
proto-generate:
	buf generate
