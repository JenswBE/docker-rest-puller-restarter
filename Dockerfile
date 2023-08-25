# Also update GitHub Actions workflow when bumping
# Based on https://www.docker.com/blog/faster-multi-platform-builds-dockerfile-cross-compilation-guide/
FROM --platform=${BUILDPLATFORM} docker.io/library/golang:1.21 AS builder
WORKDIR /src/
RUN GOARCH=amd64 go install golang.org/x/vuln/cmd/govulncheck@latest
COPY . .
RUN govulncheck ./...
ARG TARGETOS TARGETARCH TARGETVARIANT
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOARM=${TARGETVARIANT#v} go build -ldflags='-extldflags=-static' -o /bin/service

FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/docs/index.html /docs/index.html
COPY --from=builder /src/docs/openapi.yml /docs/openapi.yml
COPY --from=builder /bin/service /service
EXPOSE 8080
ENTRYPOINT ["./service"]
