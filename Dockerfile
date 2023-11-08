# syntax=docker/dockerfile:1

# Comments are provided throughout this file to help you get started.
# If you need more help, visit the Dockerfile reference guide at
# https://docs.docker.com/engine/reference/builder/

################################################################################
# Create a stage for building the application.
ARG GO_VERSION=1.21
FROM golang:${GO_VERSION} AS build
WORKDIR /src

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage bind mounts to go.sum and go.mod to avoid having to copy them into
# the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# Build the application.
# Leverage a cache mount to /go/pkg/mod/ to speed up subsequent builds.
# Leverage a bind mount to the current directory to avoid having to copy the
# source code into the container.
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    go build -o /bin/server -ldflags '-extldflags "-static"' src/main.go
COPY src/assets /app/src/assets
COPY src/templates /app/src/templates

################################################################################
# Create a new stage for running the application that contains the minimal
# runtime dependencies for the application. This often uses a different base
# image from the build stage where the necessary files are copied from the build
# stage.
ARG GO_VERSION=1.21
FROM alpine AS final

# Copy the executable from the "build" stage.
COPY --from=build /bin/server /bin/

COPY --from=build /app/src/assets /go/src/assets
COPY --from=build /app/src/templates /go/src/templates

# Expose the port that the application listens on.
EXPOSE 8080

WORKDIR /go

# What the container should run when it is started.
ENTRYPOINT [ "/bin/server" ]
