#syntax=docker/dockerfile:1.4.0

FROM base

ARG TERRAFORM_GOOGLE_VERSION=4.83.0

RUN echo $TERRAFORM_GOOGLE_VERSION

RUN apk add --no-cache --update aws-cli

# Install google provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-google/${TERRAFORM_GOOGLE_VERSION}/terraform-provider-google_${TERRAFORM_GOOGLE_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-aws_${TERRAFORM_GOOGLE_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/*
