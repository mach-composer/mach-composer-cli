FROM alpine:3 AS certs

FROM scratch

ARG GOOS=linux
ARG GOARCH=amd64_v1

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY dist/mach-composer_${GOOS}_${GOARCH}/bin/mach-composer /usr/local/bin/mach-composer

ENTRYPOINT ["/usr/local/bin/mach-composer"]
