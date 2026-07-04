FROM golang:1.25.11-alpine@sha256:cf92d27916083273afc9eac40a6150723fc3cc82c264957326e5fa34f010e2da AS builder
RUN apk --no-cache add ca-certificates git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
  go build -o /provider ./cmd/provider/

FROM alpine:3.23.4
RUN apk --no-cache add ca-certificates bash

ARG TARGETOS
ARG TARGETARCH

ENV USER_ID=65532
USER ${USER_ID}
EXPOSE 8080

COPY --from=builder /provider /usr/local/bin/provider
ENTRYPOINT ["provider"]
