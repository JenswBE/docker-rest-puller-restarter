FROM --platform=${BUILDPLATFORM} golang:1.18 AS builder-base

FROM builder-base AS builder-amd64
ENV GOOS=linux
ENV GOARCH=amd64

FROM builder-base AS builder-armv6
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=6

FROM builder-base AS builder-armv7
ENV GOOS=linux
ENV GOARCH=arm
ENV GOARM=7

FROM builder-base AS builder-arm64
ENV GOOS=linux
ENV GOARCH=arm64

FROM builder-${TARGETARCH}${TARGETVARIANT} AS builder
WORKDIR /src/
COPY . .
RUN CGO_ENABLED=0 go build -ldflags='-extldflags=-static' -o /bin/service

FROM alpine
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/docs/index.html /docs/index.html
COPY --from=builder /src/docs/openapi.yml /docs/openapi.yml
COPY --from=builder /bin/service /service
EXPOSE 8080
ENTRYPOINT ["./service"]
