.PHONY: run test verify container
run:
	go run ./cmd/server
test:
	go test -race ./...
verify:
	test -z "$$(gofmt -l .)"
	go test -race ./...
container:
	docker build -t quota-lab .
