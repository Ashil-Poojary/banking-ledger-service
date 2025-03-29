FROM golang:1.24-alpine


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o banking-ledger-service ./cmd/api

EXPOSE 8080

CMD ["./banking-ledger-service"]
