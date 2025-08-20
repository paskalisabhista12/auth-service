# ---------- Build stage ----------
FROM golang:1.24.5 AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the whole project
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/main.go

# ---------- Runtime stage ----------
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/auth-service .

# Run as non-root for security
RUN adduser -D appuser
USER appuser

EXPOSE 8080
CMD ["./auth-service"]
