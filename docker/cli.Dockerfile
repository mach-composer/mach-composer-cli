FROM goreleaser/goreleaser:v1.15.2 AS builder
ARG GORELEASER_ARGS
ARG GOOS=linux
ARG GOARCH=amd64

COPY . /code/
WORKDIR /code/
RUN GOOS=${GOOS} GOARCH=${GOARCH} goreleaser build --single-target --output /code/dist/mach-composer --skip-validate ${GORELEASER_ARGS}

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /code/dist/mach-composer /mach-composer
ENTRYPOINT ["/mach-composer"]
