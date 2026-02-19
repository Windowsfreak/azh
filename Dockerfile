# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o azh ./cmd/app

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/azh .
EXPOSE 8080
CMD ["./azh"]
