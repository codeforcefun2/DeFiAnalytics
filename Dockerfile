# Stage 1: Build the Go application
FROM golang:1.18 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o defi-analytics ./cmd/server

# Stage 2: Create a minimal runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/defi-analytics .
EXPOSE 8080
CMD ["./defi-analytics"]
