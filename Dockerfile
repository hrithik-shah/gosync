# --- build stage ---
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /gosync ./cmd/server

# --- final stage ---
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /gosync .

EXPOSE 8080

ENTRYPOINT ["./gosync"]
CMD ["start"]