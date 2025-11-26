########################
# Builder Stage
########################
FROM golang:1.25.4 AS build

WORKDIR /app

# Copy source code
COPY . .

# Download Go dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/uzaraka

########################
# Final Stage
########################
FROM alpine:3.22.0

ENV TZ=UTC

# Install runtime dependencies and set up directories
RUN apk --update --no-cache add tzdata ca-certificates curl

# Copy built binary
COPY --from=build /app/uzaraka /usr/local/bin/uzaraka
RUN chmod a+x /usr/local/bin/uzaraka && ln -s /usr/local/bin/uzaraka /uzaraka

# Set working directory
WORKDIR /app
COPY views /app/views
# Expose HTTP Port
EXPOSE 8080

ENTRYPOINT ["/uzaraka"]