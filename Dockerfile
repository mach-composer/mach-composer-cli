ARG PYTHON_VERSION="3.8.11"

FROM golang:1.15.6 AS go-builder
# RUN go get -d -v golang.org/x/net/html
RUN GO111MODULE=on go get -u go.mozilla.org/sops/v3/cmd/sops@v3.7.1 && \
    cd $GOPATH/pkg/mod/go.mozilla.org/sops/v3@v3.7.1 && \
    make install


FROM python:${PYTHON_VERSION}-alpine
COPY --from=go-builder /go/bin/sops /usr/bin/

ENV AZURE_CLI_VERSION=2.26.1
ENV TERRAFORM_VERSION=0.14.5
ENV TERRAFORM_EXTERNAL_VERSION=1.2.0
ENV TERRAFORM_AZURE_VERSION=2.86.0
ENV TERRAFORM_AWS_VERSION=3.70.0
ENV TERRAFORM_NULL_VERSION=2.1.2
ENV TERRAFORM_COMMERCETOOLS_VERSION=0.29.3
ENV TERRAFORM_CONTENTFUL_VERSION=0.1.0
ENV TERRAFORM_AMPLIENCE_VERSION=0.2.2
ENV TERRAFORM_SENTRY_VERSION=0.6.0

RUN apk add --no-cache --virtual .build-deps g++ libffi-dev openssl-dev wget unzip make \
    && apk add bash ca-certificates git libc6-compat openssl openssh-client jq curl

# Install Azure CLI
RUN pip --no-cache-dir install azure-cli==${AZURE_CLI_VERSION}

# Pre-install Terreform plugins
ENV TF_PLUGIN_CACHE_DIR=/root/.terraform.d/plugin-cache
ENV TERRAFORM_PLUGINS_PATH=/root/.terraform.d/plugins/linux_amd64
RUN mkdir -p ${TF_PLUGIN_CACHE_DIR}
RUN mkdir -p ${TERRAFORM_PLUGINS_PATH}

# Install terraform
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip -n terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*



# Install null provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-null/${TERRAFORM_NULL_VERSION}/terraform-provider-null_${TERRAFORM_NULL_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-null_${TERRAFORM_NULL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install external provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-external/${TERRAFORM_EXTERNAL_VERSION}/terraform-provider-external_${TERRAFORM_EXTERNAL_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-external_${TERRAFORM_EXTERNAL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install aws provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-aws/${TERRAFORM_AWS_VERSION}/terraform-provider-aws_${TERRAFORM_AWS_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-aws_${TERRAFORM_AWS_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install azure provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-azurerm/${TERRAFORM_AZURE_VERSION}/terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install commercetools provider
RUN cd /tmp && \
    wget https://github.com/labd/terraform-provider-commercetools/releases/download/v${TERRAFORM_COMMERCETOOLS_VERSION}/terraform-provider-commercetools_${TERRAFORM_COMMERCETOOLS_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-commercetools_${TERRAFORM_COMMERCETOOLS_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install contentful provider
RUN cd /tmp && \
    wget https://github.com/labd/terraform-provider-contentful/releases/download/v${TERRAFORM_CONTENTFUL_VERSION}/terraform-provider-contentful_${TERRAFORM_CONTENTFUL_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-contentful_${TERRAFORM_CONTENTFUL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install amplience provider
RUN cd /tmp && \
    wget https://github.com/labd/terraform-provider-amplience/releases/download/v${TERRAFORM_AMPLIENCE_VERSION}/terraform-provider-amplience_${TERRAFORM_AMPLIENCE_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-amplience_${TERRAFORM_AMPLIENCE_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

# Install sentry provider
RUN cd /tmp && \
    wget https://github.com/jianyuan/terraform-provider-sentry/releases/download/v${TERRAFORM_SENTRY_VERSION}/terraform-provider-sentry_${TERRAFORM_SENTRY_VERSION}_linux_amd64.zip && \
    unzip -n terraform-provider-sentry_${TERRAFORM_SENTRY_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH} && \
    rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

RUN mkdir /code
RUN mkdir /deployments
WORKDIR /code

ADD requirements.txt .
RUN pip install -r requirements.txt
COPY src /code/src/
ADD MANIFEST.in .
ADD setup.cfg .
ADD setup.py .
RUN python setup.py bdist_wheel && pip install dist/mach_composer-$(python setup.py --version)-py3-none-any.whl

RUN apk del .build-deps

ENTRYPOINT ["mach"]
