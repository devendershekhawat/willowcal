# Multi-stage build for efficient image size

# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/web

# Copy package files
COPY web/package*.json ./

# Install dependencies
RUN npm ci

# Copy frontend source
COPY web/ ./

# Build frontend
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Install git (required for Go modules)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o willowcal ./cmd/willowcal

# Stage 3: Final runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates git nodejs npm

WORKDIR /app

# Copy binary from builder
COPY --from=backend-builder /app/willowcal /app/willowcal

# Copy frontend build from frontend-builder
COPY --from=frontend-builder /app/web/dist /app/web/dist

# Copy startup script
COPY docker-entrypoint.sh /app/docker-entrypoint.sh
RUN chmod +x /app/docker-entrypoint.sh

# Create workspace directory
RUN mkdir -p /app/workspace

# Expose ports
EXPOSE 8080

# Set environment variables
ENV PORT=8080
ENV WORKSPACE_DIR=/app/workspace
ENV STATIC_DIR=/app/web/dist

# Run the application
ENTRYPOINT ["/app/docker-entrypoint.sh"]
