FROM golang:1.17-alpine3.16 AS builder
ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o userapi cmd/userapi/main.go

FROM alpine:3.16
COPY --from=builder ["/app/userapi", "/app/userapi"]

ENTRYPOINT ["/app/userapi"]