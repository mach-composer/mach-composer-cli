FROM scratch

ARG GOOS=linux
ARG GOARCH=amd64_v1

COPY . /mach-composer
RUN ln -s /mach-composer/dist/mach-composer__${GOOS}_${GOARCH}/bin/mach-composer /usr/local/bin/mach-composer


ENTRYPOINT ["mach-composer"]
