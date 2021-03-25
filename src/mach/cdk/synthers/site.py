import logging
from typing import TYPE_CHECKING, Dict

import cdktf_cdktf_provider_aws as aws
from mach.cdk.imports import commercetools
from mach.utils import slugify

from .abstract import Synther

if TYPE_CHECKING:
    from mach import types

logger = logging.getLogger(__name__)


class SiteSynther(Synther):
    def synth(self, obj: "types.Site"):
        self.subsynth(SiteAWSSettingsSynther, obj.aws)
        self.synth_endpoints(obj)

        if (
            self.config.general_config.sentry
            and self.config.general_config.sentry.managed
        ):
            pass

        if obj.commercetools:
            self.synth_commercetools(obj.commercetools)

    def synth_endpoints(self, obj):
        if obj.aws:
            self.subsynth(AWSEndpointsSynther, obj)

    def synth_commercetools(self, ct: "types.CommercetoolsSettings"):
        commercetools.CommercetoolsProvider(
            self.stack,
            "commercetools",
            client_id=ct.client_id,
            client_secret=ct.client_secret,
            project_key=ct.project_key,
            scopes=ct.scopes,
            token_url=ct.token_url,
            api_url=ct.api_url,
        )

        commercetools.ProjectSettings(
            self.stack,
            "project",
            name=ct.project_key,
            countries=ct.countries,
            currencies=ct.currencies,
            languages=ct.languages,
            messages={
                "enabled": "true" if ct.messages_enabled else "false",
            },
        )


class SiteAWSSettingsSynther(Synther):
    def synth(self, obj: "types.SiteAWSSettings"):
        aws_extra_kwargs = {}
        if obj.deploy_role_name:
            aws_extra_kwargs["assume_role"] = aws.AwsProviderAssumeRole(
                role_arn=f"arn:aws:iam::{obj.account_id}:role/{obj.deploy_role_name}"
            )

        aws.AwsProvider(self.stack, "aws", region=obj.region, **aws_extra_kwargs)
        for provider in obj.extra_providers:
            aws.AwsProvider(
                self.stack,
                "aws",
                alias=provider.name,
                region=provider.region,
                **aws_extra_kwargs,
            )


class AWSEndpointsSynther(Synther):
    def synth(self, obj: "types.Site"):
        if not obj.used_endpoints:
            logger.debug("No used endpoints for site. Will skip")
            return

        zones = {
            zone: aws.DataAwsRoute53Zone(
                self.stack,
                slugify(zone),
                name=zone,
            )
            for zone in obj.dns_zones
        }

        for endpoint in obj.used_endpoints:
            self.synth_endpoint(endpoint, zones)

        # {% include 'partials/endpoints/aws_url_locals.tf' %}

    def synth_endpoint(
        self, endpoint: "types.Endpoint", zones: Dict[str, aws.DataAwsRoute53Zone]
    ):
        ep_slug = slugify(endpoint.key)
        gateway = aws.Apigatewayv2Api(
            self.stack,
            f"{ep_slug}_gateway",
            name=f"{self.site.identifier}-{ep_slug}-api",
            protocol_type="HTTP",
        )

        aws.Apigatewayv2Route(
            self.stack,
            f"{ep_slug}_application",
            api_id=gateway.id,
            route_key="$default",
        )

        #
        # api gateway stage
        #
        stage_kwargs = dict(
            name="$default",
            description="Stage for default release",
            api_id=gateway.id,
            auto_deploy=True,
            # TODO: reference components (modules)
            # depends_on = [
            #     {% for component in endpoint.components %}
            #         module.{{ component.name }},
            #     {% endfor %}
            # ]
        )
        default_route_settings = {}
        if endpoint.throttling_burst_limit:
            default_route_settings[
                "throttling_burst_limit"
            ] = endpoint.throttling_burst_limit
        if endpoint.throttling_rate_limit:
            default_route_settings[
                "throttling_rate_limit"
            ] = endpoint.throttling_rate_limit
        if default_route_settings:
            stage_kwargs["default_route_settings"] = default_route_settings

        stage = aws.Apigatewayv2Stage(self.stack, f"{ep_slug}_default", **stage_kwargs)

        if endpoint.url:
            acm_cert = aws.AcmCertificate(
                self.stack,
                ep_slug,
                domain_name=endpoint.url,
                validation_method="DNS",
            )

            # TODO: for_each should be supported in cdktf before we can implement this
            # aws.Route53Record(
            #     self.stack,
            #     f"{ep_slug}_acm_validation",
            #     # for_each = {
            #     #     for dvo in aws_acm_certificate.{{ endpoint.key|slugify }}.domain_validation_options : dvo.domain_name => {
            #     #     name   = dvo.resource_record_name
            #     #     record = dvo.resource_record_value
            #     #     type   = dvo.resource_record_type
            #     #     }
            #     # }

            #     # allow_overwrite = true
            #     # name            = each.value.name
            #     # records         = [each.value.record]
            #     # ttl             = 60
            #     # type            = each.value.type
            #     # zone_id         = data.aws_route53_zone.{{ endpoint.zone|slugify }}.zone_id
            # )

            # Route53 mappings
            # TODO: Resources with the same name (altho of different types) doesn't seem to
            # be supported yet in the CDK.
            # Related to https://github.com/hashicorp/terraform-cdk/pull/329
            gw_domain = aws.Apigatewayv2DomainName(
                self.stack,
                ep_slug,
                domain_name=endpoint.url,
                domain_name_configuration=[
                    aws.Apigatewayv2DomainNameDomainNameConfiguration(
                        certificate_arn=acm_cert.arn,
                        endpoint_type="REGIONAL",
                        security_policy="TLS_1_2",
                    ),
                ],
            )

            aws.Route53Record(
                self.stack,
                ep_slug,
                name=gw_domain.domain_name,
                type="A",
                zone_id=zones[endpoint.zone].id,
                alias=[
                    aws.Route53RecordAlias(
                        name=gw_domain.domain_name_configuration[0].target_domain_name,
                        zone_id=gw_domain.domain_name_configuration[0].hosted_zone_id,
                        evaluate_target_health=False,
                    ),
                ],
            )

            aws.Apigatewayv2ApiMapping(
                self.stack,
                ep_slug,
                api_id=gateway.id,
                stage=stage.id,
                domain_name=endpoint.url,
            )
