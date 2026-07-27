# Build stage: contains the Go compiler.
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app main.go

# Runtime stage: contains only the compiled binary and a small base image.
FROM alpine:3.22
RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /app /app
USER app
EXPOSE 8081
ENTRYPOINT ["/app"]
