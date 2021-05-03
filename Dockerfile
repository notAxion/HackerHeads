# syntax=docker/dockerfile:1.2

# base to use use it for build, test, lint
FROM --platform=${BUILDPLATFORM} golang:1.16.3 AS base
WORKDIR '/app'
COPY go.* .
RUN go mod download
COPY . .

# run from build
FROM base AS build
ARG TARGETOS
ARG TARGETARCH
# RUN --mount=target=.,rw \
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main .
CMD ["./main"]

# linter base image 
FROM golangci/golangci-lint:v1.39 AS lint-base

# linter to check for good practice code
FROM base AS lint
RUN --mount=target=. \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    golangci-lint run --timeout 10m0s ./...