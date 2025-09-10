# ---------- Build stage ----------
FROM golang:1.24.5 AS builder

WORKDIR /app

# Copy go.mod & go.sum first (dependency cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# ✅ Force static build (no libc dependency)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o auth-service -ldflags="-w -s" ./cmd/main.go

# ---------- Runtime stage ----------
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder stage only
COPY --from=builder /app/auth-service .

# Run as non-root for safety
RUN adduser -D appuser
USER appuser

EXPOSE 8080
CMD ["./auth-service"]
