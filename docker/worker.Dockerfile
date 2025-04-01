FROM golang:1.24 as builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy the entire project and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker_service ./cmd/worker

# Use a minimal base image for the final container
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/worker_service /app/worker_service

# Set execution permissions and define the entrypoint
CMD ["/app/worker_service"]
