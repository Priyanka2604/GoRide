# Stage 1: Build
FROM golang:1.24 AS builder

WORKDIR /app

# Download dependencies first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the service
RUN go build -o booking_svc ./cmd

# Stage 2: Run (minimal image)
FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/booking_svc .
EXPOSE 8080
CMD ["./booking_svc"]