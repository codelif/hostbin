# syntax=docker/dockerfile:1.7

FROM golang:1.26.1-bookworm AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
	go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/hostbin ./cmd/server

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app

ENV LISTEN_ADDR=0.0.0.0:8080 \
	DB_PATH=/var/lib/hostbin/data.db \
	LOG_LEVEL=info

COPY --from=build /out/hostbin /usr/local/bin/hostbin

VOLUME ["/var/lib/hostbin"]
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/hostbin"]
