FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
RUN apk add --no-cache make
RUN apk add --no-cache ca-certificates
WORKDIR /src
RUN git clone https://github.com/labd/mach-composer.git /src
RUN make build

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/bin/mach-composer /bin/mach-composer
ENTRYPOINT ["/bin/mach-composer"]
