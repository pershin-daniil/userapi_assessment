# Refactoring User API in Go 💙

This is a small API that handles User entity, where the storage is a JSON file.

## Requirements 😏

- The storage must remain a JSON file.
- User structure should not be reduced.
- The application should not lose any existing functionality.

## Installation ✅

Clone the repository and download dependencies.

```shell
git clone git@github.com:pershin-daniil/userapi_assessment.git
cd userapi_assessment
go mod tidy
```

Use `make run` - to start service locally on `localhost:3333` OR `make rund` - to start service in docker
on `localhost:8080`.

If you want to test service use command below and then use scripts in API Endpoints.

Also, you can run service locally and test it with http requests [here](./http)

```shell
make rund
```

## API Endpoints 🖊

### GET /users

Returns a list of all users in the storage.

#### Request

```shell
curl --location 'http://localhost:8080/api/v1/users'
```

### Response

```json
{
  "increment": 2,
  "list": {
    "1": {
      "id": 1,
      "displayName": "Ivan",
      "email": "test@mail.com",
      "created": "2023-04-16T19:04:53.525308741Z",
      "updated": "2023-04-16T19:04:53.525308741Z"
    },
    "2": {
      "id": 2,
      "displayName": "Ivan",
      "email": "test@mail.com",
      "created": "2023-04-16T19:07:52.723930891Z",
      "updated": "2023-04-16T19:07:52.723930891Z"
    }
  }
}
```

### POST /users

Adds a new user to the storage.

#### Request

```shell
curl --location 'localhost:8080/api/v1/users' \
--header 'Content-Type: application/json' \
--data-raw '{
  "displayName": "Ivan",
  "email": "test@mail.com"
}'
```

#### Response

```json
{
  "id": 1,
  "displayName": "Ivan",
  "email": "test@mail.com",
  "created": "2023-04-16T19:04:53.525308741Z",
  "updated": "2023-04-16T19:04:53.525308741Z"
}
```

### GET /users/{id}

Returns a specific user by ID.

#### Request

```shell
curl --location 'localhost:8080/api/v1/users/1'
```

#### Response

```json
{
  "id": 1,
  "displayName": "Ivan",
  "email": "test@mail.com",
  "created": "2023-04-16T19:04:53.525308741Z",
  "updated": "2023-04-16T19:04:53.525308741Z"
}
```

### PATCH /users/{id}

Updates an existing user by ID. You can change only displayName.

#### Request

```shell
curl --location --request PATCH 'localhost:8080/api/v1/users/1' \
--header 'Content-Type: application/json' \
--data '{
  "displayName": "MASHA"
}'
```

#### Response

```json
{
  "id": 1,
  "displayName": "MASHA",
  "email": "test@mail.com",
  "created": "2023-04-16T19:04:53.525308741Z",
  "updated": "2023-04-16T19:28:45.540040396Z"
}
```

### DELETE /users/{id}

Deletes a user from the storage by ID.

#### Request

```shell
curl --location --request DELETE 'localhost:8080/api/v1/users/4'
```

#### Response

Status: 204 No Content

## Useful information 🤔

👉 Full task text [here](./docs/task.md)

👉 Make file commands

```makefile
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
```