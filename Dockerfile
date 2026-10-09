# syntax=docker/dockerfile:1

ARG GO_VERSION=1.27.2

# Cross-compiles on the build host, so multi-arch images don't need QEMU
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags "-s -w -buildid=" -o /out/netstatus-api-go .

FROM scratch

COPY --from=build /out/netstatus-api-go /usr/local/bin/netstatus-api-go

# Mount config.json at /etc/netstatus-api-go/config.json, defaults are used without one
USER 65532:65532
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/netstatus-api-go"]
