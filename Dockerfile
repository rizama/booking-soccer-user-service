# Stage 1: Builder - Menggunakan Go image untuk compile aplikasi
FROM golang:1.24.4-alpine3.21 AS builder

# Update package list Alpine Linux
RUN apk update
# Install dependencies yang diperlukan untuk build:
# - git: untuk download dependencies dari Git repositories
# - openssh: untuk SSH connections
# - tzdata: untuk timezone data
# - build-base: untuk C compiler (gcc, make, dll)
# - python3: jika ada dependencies yang memerlukan Python
# - net-tools: untuk network utilities
RUN apk add git openssh tzdata build-base python3 net-tools

# Set working directory di dalam container
WORKDIR /app

# Copy file .env.example dan rename menjadi .env
# Ini untuk menyediakan environment variables default
COPY .env.example .env
# Copy semua file dari current directory ke container
COPY . .

# Install gin framework (mungkin untuk hot reload development)
RUN go install github.com/buu700/gin@latest
# Download dan install semua Go module dependencies
RUN go mod tidy

# Build aplikasi menggunakan Makefile
# Ini akan mengcompile Go code menjadi binary executable
RUN make build

# Stage 2: Production - Menggunakan Alpine Linux yang lebih kecil untuk runtime
FROM alpine:latest

# Update package list, upgrade existing packages, dan install runtime dependencies:
# - tzdata: untuk timezone support
# - curl: untuk health checks atau HTTP requests
# - mkdir /app: buat directory untuk aplikasi
RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    apk --no-cache add curl && \
    mkdir /app

# Set working directory untuk production container
WORKDIR /app

# Expose port 8081 untuk aplikasi Go
# Port ini harus match dengan port yang digunakan di aplikasi
EXPOSE 8081

# Copy compiled binary dan files dari builder stage ke production stage
# Ini mengambil hasil build dari stage pertama
COPY --from=builder /app /app

# Set entry point untuk menjalankan aplikasi
# Ketika container start, akan menjalankan binary user-service
ENTRYPOINT [ "/app/user-service" ]