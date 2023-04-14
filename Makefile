lint:
	gofumpt -w .
	go mod tidy
	golangci-lint run

run:
	go run cmd/userapi/main.go

test:
	go test -v ./tests/userapi_test.go
