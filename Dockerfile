FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

# Install Air for hot reloading
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go.mod and go.sum first for better caching

COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Set the command to run Air
CMD ["air", "-c", ".air.toml"]
