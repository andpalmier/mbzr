# The binary is copied in rather than compiled here, so the image ships the
# exact artifact the release archives contain.
#
# GoReleaser lays the build context out as <os>/<arch>/mbzr and buildx sets
# TARGETPLATFORM to match, which is how one Dockerfile serves both platforms.
FROM alpine:3.24
RUN apk add --no-cache ca-certificates && \
    adduser -D -g '' appuser

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/mbzr /mbzr

USER appuser
ENTRYPOINT ["/mbzr"]
