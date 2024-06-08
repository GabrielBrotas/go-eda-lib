# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21.9
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /src

# Cache and bind go.mod and go.sum
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# Build the application
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=1 go build -o /bin/server ./cmd/walletcore

# Use an Alpine image with glibc for the final image
FROM frolvlad/alpine-glibc:latest AS final

# Install necessary packages including librdkafka
RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        librdkafka-dev \
        && \
        update-ca-certificates

# Copy the built binary and set permissions
COPY --from=build /bin/server /bin/server
RUN chmod +x /bin/server

# Create a non-root user
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    appuser
USER appuser

EXPOSE 8080

ENTRYPOINT ["/bin/server"]
