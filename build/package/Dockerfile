# Build module
FROM golang:1.20.7-alpine3.18 AS builder

# Copy source
WORKDIR /src
COPY go.mod go.sum .
RUN go mod download

COPY . ./
RUN go build -o bin/main cmd/main.go

# Prod module
FROM alpine:3.18 AS prod

# Copy binary
WORKDIR /app
COPY --from=builder /src/bin/main .

ENTRYPOINT ./main