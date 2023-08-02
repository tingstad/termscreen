ARG GO_VERSION

FROM golang:${GO_VERSION}-alpine AS build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod go.sum termscreen*.go LICENSE Makefile Dockerfile ./
COPY cmd/ ./cmd/

RUN go build -o main cmd/main.go

FROM alpine:3.18.2

WORKDIR /dist

COPY --from=build /build/*.* /build/LICENSE /build/Makefile /build/Dockerfile .
COPY --from=build /build/cmd/ ./cmd/

CMD ["/dist/main"]

