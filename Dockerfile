FROM golang:1.23.2 AS build
WORKDIR /go/src/app
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOPROXY=direct
RUN go build -v -o app .

# Instead of scratch, use a minimal Alpine image
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates && update-ca-certificates

# Create a non-root user to run the application
RUN adduser -D appuser

# Copy the compiled application from the build stage
COPY --from=build /go/src/app/app /usr/local/bin/app

# Switch to the non-root user for better security
USER appuser

# Set the entry point to your application
ENTRYPOINT ["/usr/local/bin/app"]
