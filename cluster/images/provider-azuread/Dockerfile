FROM golang:1.25 AS builder
WORKDIR /workspace

ARG TARGETOS=linux
ARG TARGETARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -o /out/provider ./cmd/provider

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /out/provider /usr/local/bin/provider
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/provider"]
