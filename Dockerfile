# Multistage Build

### CREATE DOCKERMASTER USER
FROM alpine:3.6 AS alpine
RUN adduser -D -u 10001 dockmaster
RUN apk --update add ca-certificates

## MAIN IMAGE
FROM scratch
LABEL Name=s3-region-stats
LABEL Author=davyj0nes

COPY --from=alpine /etc/passwd /etc/passwd
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV AWS_DEFAULT_REGION=eu-west-1

ADD s3-region-stats_static /s3-region-stats
USER dockmaster

CMD ["./s3-region-stats"]
