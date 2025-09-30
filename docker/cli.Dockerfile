FROM scratch

ARG GOOS=linux
ARG GOARCH=amd64_v1

COPY ./dist/mach-composer_${GOOS}_${GOARCH}/bin/mach-composer /mach-composer
ENTRYPOINT ["/mach-composer"]
