# Build stage
FROM golang:1.25-alpine AS builder
ARG TARGETARCH
WORKDIR /app
RUN apk add --no-cache git make
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -ldflags="-s -w" -o mbzr .

# Final stage
FROM alpine:latest
WORKDIR /
RUN apk add --no-cache ca-certificates && \
    adduser -D -g '' appuser
COPY --from=builder /app/mbzr /mbzr
USER appuser
ENTRYPOINT ["/mbzr"]
