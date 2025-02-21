# Development Stage
FROM golang:1.24-alpine AS dev

# Install required packages
RUN apk upgrade --no-cache & apk add --no-cache \
  git \
  bash \
  gcc \
  musl-dev \
  wget tar 

# Install Just cli for task automation
RUN wget https://github.com/casey/just/releases/download/1.36.0/just-1.36.0-x86_64-unknown-linux-musl.tar.gz \
  && tar -xzf just-1.36.0-x86_64-unknown-linux-musl.tar.gz \
  && mv just /usr/local/bin/ \
  && chmod +x /usr/local/bin/just \
  && rm just-1.36.0-x86_64-unknown-linux-musl.tar.gz \
  && just --version

# Install Go tools
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
  go install github.com/air-verse/air@latest

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first, to cache dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Expose ports
EXPOSE 51000 41000
