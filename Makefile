.PHONY: build test test-integration lint run-facilitator run-resource run-client run-demo run-demo-permit2 clean

build:
	go build -o facilitator ./cmd/facilitator
	go build -o resource ./cmd/resource
	go build -o client ./cmd/client

test:
	go test ./... -v -count=1

test-integration:
	go test ./test/integration/... -v -count=1 -tags=integration

lint:
	golangci-lint run ./...

run-facilitator:
	go run ./cmd/facilitator

run-resource:
	go run ./cmd/resource

run-client:
	go run ./cmd/client

run-demo:
	go run ./cmd/demo

run-demo-permit2:
	ASSET_TRANSFER_METHOD=permit2 go run ./cmd/demo

clean:
	rm -f facilitator resource client
