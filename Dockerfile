FROM alpine:3.19
RUN apk add --no-cache ca-certificates && \
    adduser -D -g '' appuser
COPY mbzr /mbzr
USER appuser
ENTRYPOINT ["/mbzr"]
