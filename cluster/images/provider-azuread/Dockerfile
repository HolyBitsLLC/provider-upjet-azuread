FROM golang:1.25 AS builder
WORKDIR /workspace

ARG TARGETOS=linux
ARG TARGETARCH=amd64

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -o /out/provider ./cmd/provider

FROM alpine:3.23.4
RUN apk --no-cache add ca-certificates bash

ENV USER_ID=65532
USER ${USER_ID}
EXPOSE 8080

COPY --from=builder /out/provider /usr/local/bin/provider
ENTRYPOINT ["provider"]
