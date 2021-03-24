from typing import TYPE_CHECKING

import cdktf_cdktf_provider_aws as aws
from cdktf import App, TerraformStack

from .imports import commercetools

if TYPE_CHECKING:
    from .types import MachConfig, Site


class SiteStack(TerraformStack):
    def __init__(self, app, *, config: "MachConfig", site: "Site"):
        super().__init__(app, site.identifier)

        if site.aws:
            aws_extra_kwargs = {}
            if site.aws.deploy_role_name:
                aws_extra_kwargs["assume_role"] = aws.AwsProviderAssumeRole(
                    role_arn=f"arn:aws:iam::{site.aws.account_id}:role/{site.aws.deploy_role_name}"
                )

            aws.AwsProvider(self, "aws", region=site.aws.region, **aws_extra_kwargs)

        for provider in site.aws.extra_providers:
            aws.AwsProvider(
                self,
                "aws",
                alias=provider.name,
                region=provider.region,
                **aws_extra_kwargs,
            )

        if config.general_config.sentry and config.general_config.sentry.managed:
            pass

        if site.commercetools:
            commercetools.CommercetoolsProvider(
                self,
                "commercetools",
                client_id=site.commercetools.client_id,
                client_secret=site.commercetools.client_secret,
                project_key=site.commercetools.project_key,
                scopes=site.commercetools.scopes,
                token_url=site.commercetools.token_url,
                api_url=site.commercetools.api_url,
            )

            commercetools.ProjectSettings(
                self,
                "project",
                name=site.commercetools.project_key,
                countries=site.commercetools.countries,
                currencies=site.commercetools.currencies,
                languages=site.commercetools.languages,
                messages={
                    "enabled": "true"
                    if site.commercetools.messages_enabled
                    else "false",
                },
            )


def generate(config: "MachConfig", site: "Site", outdir: str):
    app = App(outdir=outdir)
    SiteStack(app, config=config, site=site)
    app.synth()
