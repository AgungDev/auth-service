# syntax=docker/dockerfile:1
FROM golang:tip-alpine3.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth_service ./cmd/server

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/auth_service ./auth_service
COPY .env ./
COPY jwt_private.pem ./
COPY jwt_public.pem ./
EXPOSE 8080
CMD ["./auth_service"]
