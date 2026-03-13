FROM golang:1.23-alpine AS build
ARG BUILD_VERSION=dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.version=${BUILD_VERSION}" -o multipass-mcp .

FROM alpine:latest
ARG BUILD_VERSION=dev
LABEL org.opencontainers.image.version="${BUILD_VERSION}"
COPY --from=build /app/multipass-mcp /usr/local/bin/
ENTRYPOINT ["multipass-mcp"]
