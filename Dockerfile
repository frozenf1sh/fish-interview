FROM --platform=$BUILDPLATFORM golang:1.26-alpine3.22 AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY content ./content
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags='-s -w' -o /out/fish-interview ./cmd/fish-interview

FROM alpine:3.22.1

RUN addgroup -S -g 10001 app && adduser -S -D -u 10001 -G app app
WORKDIR /app
COPY --from=build /out/fish-interview /app/fish-interview
COPY --from=build /src/content /app/content
USER 10001:10001
EXPOSE 8080
ENTRYPOINT ["/app/fish-interview"]
