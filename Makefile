build-container:
	docker build -t filequeue:latest .

run-container:
	docker run -it --name file-pubsub filequeue:latest

tools: install-mockery

install-mockery:
	go install github.com/vektra/mockery/v2@latest

generate:
	go generate ./...

test:
	go test -short -race ./...

lint:
	golangci-lint run ./...

.PHONY: build-container run-container tools generate lint