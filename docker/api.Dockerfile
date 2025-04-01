# Build Stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project and build
COPY . .
RUN go build -o banking-api ./cmd/api

# Ensure binary is executable
RUN chmod +x banking-api

# Production Stage (Distroless)
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/banking-api /app/banking-api

EXPOSE 8080

# Run the binary
ENTRYPOINT ["/app/banking-api"]
