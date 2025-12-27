# -------- Build stage --------
FROM golang:1.22.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build ONLY the main package
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o scheduler ./cmd/main

# -------- Runtime stage --------
FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/scheduler .

EXPOSE 50051

ENTRYPOINT ["./scheduler"]