FROM goreleaser/goreleaser:v1.22.1 AS builder
ARG GOOS=linux
ARG GOARCH=amd64

COPY . /code/
WORKDIR /code/
RUN GOOS=${GOOS} GOARCH=${GOARCH} goreleaser build --single-target --output /code/dist/mach-composer --skip=before --skip=validate
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /code/dist/mach-composer /mach-composer
ENTRYPOINT ["/mach-composer"]
