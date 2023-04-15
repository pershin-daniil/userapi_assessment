lint:
	gofumpt -w .
	go mod tidy
	golangci-lint run

run:
	go run cmd/userapi/main.go

test:
	go test -v ./tests/userapi_test.go

build:
	docker build -f ./deploy/local/Dockerfile -t userapi .

rund: build
	docker run --rm --name userapi -p 8080:3333 userapi
