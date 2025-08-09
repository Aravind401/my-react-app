# ---------- Stage 1: Build ----------
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest
ENV PATH="/go/bin:${PATH}"

# Generate Swagger docs
RUN swag init

# Build the binary
RUN go build -o main .

# ---------- Stage 2: Final Image ----------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

EXPOSE 8080

ENTRYPOINT ["./main"]
