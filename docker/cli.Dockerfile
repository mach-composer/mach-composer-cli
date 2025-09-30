FROM scratch
COPY ./dist/mach-composer_linux_amd64_v1/bin/mach-composer /mach-composer
ENTRYPOINT ["/mach-composer"]
