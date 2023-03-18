# syntax=docker/dockerfile:1.4-labs

ARG GOLANG_VERSION="1.20"
ARG BUILD_IMAGE="golang:${GOLANG_VERSION}-alpine"
ARG GOLANGCI_LINT_IMAGE="golangci/golangci-lint:latest-alpine"
ARG RELEASE_IMAGE="scratch"

# =============================================================================
FROM ${BUILD_IMAGE} as base

SHELL ["/usr/bin/env", "/bin/sh", "-e", "-u" ,"-o", "pipefail", "-o", "errexit", "-o", "nounset", "-c"]

WORKDIR /src/port-scanner-go

ARG GO111MODULE="on"
ARG CGO_ENABLED="0"
ARG GOARCH="amd64"
ARG GOOS="linux"
ARG APP_VERSION="docker"
ENV GO111MODULE="${GO111MODULE}" \
    CGO_ENABLED="${CGO_ENABLED}"  \
    GOARCH="${GOARCH}" \
    GOOS="${GOOS}" \
    APP_VERSION="${APP_VERSION}"

COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# =============================================================================
FROM ${GOLANGCI_LINT_IMAGE} AS lint-base

# =============================================================================
FROM base AS lint

RUN --mount=target=. \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/.cache/golangci-lint \
    golangci-lint run --timeout 10m0s ./... \
    && echo "golint is finished" > /lint.txt

# =============================================================================
FROM base AS test

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go test -v -coverprofile=/cover.out ./...

# =============================================================================
FROM base as build

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -tags musl -ldflags "-X 'main.AppVersion=${APP_VERSION}'" -o /port-scanner-go \
    && /port-scanner-go --version

RUN --mount=from=lint,src=/lint.txt,target=/tmp/lint.txt \
    --mount=from=test,src=/cover.out,target=/tmp/cover.out \
    cat /tmp/lint.txt

# ============================================================================= 
FROM ${RELEASE_IMAGE} as release
COPY --from=build /port-scanner-go /port-scanner-go
ENTRYPOINT [ "/port-scanner-go" ]
CMD [ "--version" ]
