FROM goreleaser/goreleaser:v2.12.3 AS builder

COPY . /mach-composer
RUN ln -s /mach-composer/dist/mach-composer_linux_amd64_v1/bin/mach-composer /usr/local/bin/mach

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /mach-composer/dist/mach-composer_linux_amd64_v1/bin/mach-composer /mach-composer
ENTRYPOINT ["/mach-composer"]
