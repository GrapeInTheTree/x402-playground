.PHONY: build test test-integration lint run-facilitator run-resource run-client run-demo run-demo-permit2 run-explorer run-learn run-dashboard clean

build:
	go build -o facilitator ./cmd/facilitator
	go build -o resource ./cmd/resource
	go build -o client ./cmd/client
	go build -o explorer ./cmd/explorer

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
	go run ./cmd/explorer --mode=practice --flow=eip3009

run-demo-permit2:
	go run ./cmd/explorer --mode=practice --flow=permit2

run-explorer:
	go run ./cmd/explorer

run-learn:
	go run ./cmd/explorer --mode=learn

run-dashboard:
	go run ./cmd/explorer --mode=dashboard

clean:
	rm -f facilitator resource client explorer
