FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache ca-certificates git

# Cache Dependencies
COPY go.mod ./
RUN go mod download

COPY . .

# Build with generalized architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w"\
    -o /out/rwnd ./cmd/rwnd

# Premake the logs folder and set it's owner to the google distroless non-root so they can edit / make it here ahead of time
RUN mkdir -p /out/app && \
    cp /out/rwnd /out/app/rwnd && \
    mkdir -p /out/app/.rwnd/logs && \
    chown -R 65532:65532 /out/app


FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app

COPY --from=builder /out/app /app

ENTRYPOINT [ "/app/rwnd" ]
