FROM python:3.8.5-alpine

ENV AZURE_CLI_VERSION=2.5.1
ENV TERRAFORM_VERSION=0.12.28
ENV TERRAFORM_EXTERNAL_VERSION=1.2.0
ENV TERRAFORM_AZURE_VERSION=2.19.0
ENV TERRAFORM_NULL_VERSION=2.1.2
ENV TERRAFORM_COMMERCETOOLS_VERSION=0.23.0
ENV TERRAFORM_PLUGINS_PATH=/root/.terraform.d/plugins/linux_amd64
RUN mkdir -p ${TERRAFORM_PLUGINS_PATH}

RUN apk update && \
    apk add --no-cache --virtual .build-deps g++ python3-dev libffi-dev openssl-dev && \
    apk add --no-cache --update python3 && \
    apk add bash curl tar ca-certificates git libc6-compat openssl jq unzip wget openssh-client make

# Install terraform
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin

# Install null provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-null/${TERRAFORM_NULL_VERSION}/terraform-provider-null_${TERRAFORM_NULL_VERSION}_linux_amd64.zip && \
    unzip terraform-provider-null_${TERRAFORM_NULL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH}

# Install external provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-external/${TERRAFORM_EXTERNAL_VERSION}/terraform-provider-external_${TERRAFORM_EXTERNAL_VERSION}_linux_amd64.zip && \
    unzip terraform-provider-external_${TERRAFORM_EXTERNAL_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH}

# Install azure provider
RUN cd /tmp && \
    wget https://releases.hashicorp.com/terraform-provider-azurerm/${TERRAFORM_AZURE_VERSION}/terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip && \
    unzip terraform-provider-azurerm_${TERRAFORM_AZURE_VERSION}_linux_amd64.zip -d ${TERRAFORM_PLUGINS_PATH}


# Install commercetools provider
RUN cd /tmp && \
    wget https://github.com/labd/terraform-provider-commercetools/releases/download/${TERRAFORM_COMMERCETOOLS_VERSION}/terraform-provider-commercetools-${TERRAFORM_COMMERCETOOLS_VERSION}-linux-amd64.tar.gz && \
    tar -C ${TERRAFORM_PLUGINS_PATH} -xzf terraform-provider-commercetools-${TERRAFORM_COMMERCETOOLS_VERSION}-linux-amd64.tar.gz

RUN rm -rf /tmp/* && \
    rm -rf /var/cache/apk/* && \
    rm -rf /var/tmp/*

RUN pip --no-cache-dir install azure-cli==${AZURE_CLI_VERSION}

# TODO: use build containers to optimize this, for now this works ;^)
RUN mkdir /code
RUN mkdir /deployments
WORKDIR /code

ADD requirements.txt .
RUN pip install -r requirements.txt
COPY src /code/src/
ADD MANIFEST.in .
ADD setup.cfg . 
ADD setup.py . 
RUN python setup.py bdist_wheel && pip install dist/mach-0.0.0-py3-none-any.whl


ENTRYPOINT ["mach"]
