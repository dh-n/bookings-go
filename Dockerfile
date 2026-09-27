# --- STAGE 1: BUILD THE BINARY ---
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy dependency files first (enables Docker caching for faster subsequent builds)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code
COPY . .

# Compile the Go web app into a static binary named "server"
# CGO_ENABLED=0 ensures the binary is self-contained and runs anywhere
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# --- STAGE 2: RUN THE BINARY ---
FROM alpine:3.19

WORKDIR /app

# Install security certificates (required if your app makes HTTPS requests to external APIs)
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from Stage 1
COPY --from=builder /app/server .

# Copy your web assets (templates, html, css, js) so the binary can render them
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

# Expose the port your Go web application runs on
EXPOSE 8080

# Run the web server
CMD ["./server"]
