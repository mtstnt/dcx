FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o executor .

FROM debian:stable-20250317-slim
WORKDIR /app
COPY --from=builder /app/executor .
EXPOSE 8080
CMD ["./executor"]