from dataclasses import dataclass, field
from typing import TYPE_CHECKING, Any, Dict, List, Optional, Union

from dataclasses_json import config, dataclass_json
from dataclasses_jsonschema import JsonSchemaMixin
from mach import utils
from marshmallow import ValidationError

from . import fields
from .general_config import (
    AmplienceConfig,
    AzureConfig,
    ContentfulConfig,
    FrontdoorSettings,
    SentryConfig,
)
from .shared import ComponentAzureConfig, ServicePlan

TerraformVariables = Dict[str, Any]
StoreVariables = Dict[str, TerraformVariables]
LocalizedString = Dict[str, str]
Tags = Dict[str, str]

if TYPE_CHECKING:
    from .components import ComponentConfig

__all__ = [
    "AlertGroup",
    "ApolloFederationSettings",
    "AWSProvider",
    "Endpoint",
    "CommercetoolsStore",
    "CommercetoolsChannel",
    "CommercetoolsTax",
    "CommercetoolsTaxCategory",
    "CommercetoolsProjectSettings",
    "CommercetoolsFrontendSettings",
    "CommercetoolsSettings",
    "ContentfulSettings",
    "AmplienceSettings",
    "SentryDsn",
    "Component",
    "SiteAWSSettings",
    "SiteAzureSettings",
    "Site",
]


@dataclass_json
@dataclass
class AlertGroup(JsonSchemaMixin):
    """Alert group configuration."""

    name: str
    alert_emails: List[str] = fields.list_()
    webhook_url: Optional[str] = fields.none()
    logic_app: Optional[str] = fields.none()

    @property
    def logic_app_name(self) -> Optional[str]:
        return self.logic_app.split(".")[1] if self.logic_app else None

    @property
    def logic_app_resource_group(self) -> Optional[str]:
        return self.logic_app.split(".")[0] if self.logic_app else None


@dataclass_json
@dataclass
class AWSProvider(JsonSchemaMixin):
    """AWS provider configuration."""

    name: str
    region: str
    default_tags: Tags = fields.dict_()


@dataclass_json
@dataclass
class AzureEndpoint(JsonSchemaMixin):
    session_affinity_enabled: Optional[bool] = fields.default(False)
    session_affinity_ttl_seconds: Optional[int] = fields.default(0)
    waf_policy_id: Optional[str] = fields.none()
    internal_name: Optional[str] = fields.none()


@dataclass_json
@dataclass
class AWSEndpoint(JsonSchemaMixin):
    throttling_burst_limit: Optional[int] = fields.none()
    throttling_rate_limit: Optional[int] = fields.none()
    enable_cdn: Optional[bool] = fields.default(False)


@dataclass_json
@dataclass
class Endpoint(JsonSchemaMixin):
    url: str
    key: str = field(metadata=config(exclude=lambda x: True))
    zone: Optional[str] = fields.none()
    aws: Optional[AWSEndpoint] = fields.none()
    azure: Optional[AzureEndpoint] = fields.none()

    # To be set by the parser
    components: List["Component"] = fields.list_()

    @property
    def contains_defaults(self):
        """Indicate if this endpoint contains just default values.

        Other then the `url` attribute.
        If only defaults, we can serialize the endpoints by just
        rendering the url, not the entire object.

        At this moment, we don't have any additional options, so it's always default.
        This can be extended in the future.
        """
        return True

    @property
    def subdomain(self) -> str:
        if not self.url:
            return ""

        return utils.subdomain_from_url(self.url)

    @property
    def is_root_domain(self):
        return self.url == self.zone

    def __post_init__(self):
        """Ensure endpoints have protocol stripped."""
        self.url = utils.strip_protocol(self.url)

        if not self.zone and self.url:
            try:
                self.zone = utils.dns_zone_from_url(self.url)
            except ValueError as e:
                raise ValidationError(f"Could not determine DNS zone: {e}")


@dataclass_json
@dataclass
class CommercetoolsStore(JsonSchemaMixin):
    """commercetools store definition."""

    key: str
    name: LocalizedString = fields.dict_()
    managed: bool = True
    languages: Optional[List[str]] = fields.list_()
    distribution_channels: Optional[List[str]] = fields.list_()
    supply_channels: Optional[List[str]] = fields.list_()

    def __post_init__(self):
        if self.managed and not self.name:
            raise ValidationError("name is required")


@dataclass_json
@dataclass
class CommercetoolsChannel(JsonSchemaMixin):
    """commercetools channel definition."""

    key: str
    roles: List[str]
    name: Optional[LocalizedString] = fields.none()
    description: Optional[LocalizedString] = fields.none()


@dataclass_json
@dataclass
class CommercetoolsTax(JsonSchemaMixin):
    """commercetools tax definition."""

    country: str
    amount: float
    name: str
    included_in_price: Optional[bool] = True


@dataclass_json
@dataclass
class CommercetoolsTaxCategory(JsonSchemaMixin):
    """commercetools tax categories definition."""

    key: str
    name: str
    rates: Optional[List[CommercetoolsTax]] = fields.list_()


@dataclass_json
@dataclass
class CommercetoolsZoneLocation(JsonSchemaMixin):
    country: str
    state: str = ""


@dataclass_json
@dataclass
class CommercetoolsZone(JsonSchemaMixin):
    name: str
    description: str = ""
    locations: List[CommercetoolsZoneLocation] = fields.list_()


@dataclass_json
@dataclass
class CommercetoolsFrontendSettings(JsonSchemaMixin):
    create_credentials: bool = fields.default(True)
    permission_scopes: List[str] = fields.default(
        [
            "create_anonymous_token",
            "manage_my_profile",
            "manage_my_orders",
            "manage_my_shopping_lists",
            "manage_my_payments",
            "view_products",
            "view_project_settings",
        ]
    )


@dataclass_json
@dataclass
class CommercetoolsProjectSettings(JsonSchemaMixin):
    # Project settings
    currencies: Optional[List[str]] = fields.none()
    languages: Optional[List[str]] = fields.none()
    countries: Optional[List[str]] = fields.none()
    messages_enabled: bool = fields.default(True)


@dataclass_json
@dataclass
class CommercetoolsSettings(JsonSchemaMixin):
    """commercetools configuration."""

    # Authentication settings
    project_key: str
    client_id: str
    client_secret: str
    scopes: str

    token_url: Optional[str] = fields.default(
        "https://auth.europe-west1.gcp.commercetools.com"
    )
    api_url: Optional[str] = fields.default(
        "https://api.europe-west1.gcp.commercetools.com"
    )

    # Set to false to not manage commercetools settings via mach composer
    project_settings: Optional[CommercetoolsProjectSettings] = fields.none()

    channels: Optional[List[CommercetoolsChannel]] = fields.list_()
    taxes: Optional[List[CommercetoolsTax]] = fields.list_()
    tax_categories: Optional[List[CommercetoolsTaxCategory]] = fields.list_()
    stores: List[CommercetoolsStore] = fields.list_()
    zones: List[CommercetoolsZone] = fields.list_()

    # Extra credentials
    frontend: CommercetoolsFrontendSettings = fields.none()

    def __post_init__(self):
        if not self.frontend:
            self.frontend = CommercetoolsFrontendSettings()

    @property
    def managed_stores(self):
        return [store for store in self.stores if store.managed]


@dataclass_json
@dataclass
class ContentfulSettings(JsonSchemaMixin):
    """Contentful settings."""

    space: str
    default_locale: str = "en-US"
    cma_token: str = ""
    organization_id: str = ""

    def merge(self, config: ContentfulConfig):
        self.cma_token = self.cma_token or config.cma_token
        self.organization_id = self.organization_id or config.organization_id


@dataclass_json
@dataclass
class AmplienceSettings(JsonSchemaMixin):
    """Amplience settings."""

    hub_id: str
    client_id: str = ""
    client_secret: str = ""

    def merge(self, config: AmplienceConfig):
        self.client_id = self.client_id or config.client_id
        self.client_secret = self.client_secret or config.client_secret


@dataclass_json
@dataclass
class SentryDsn(JsonSchemaMixin):
    """Specific sentry DSN settings."""

    dsn: Optional[str] = fields.none()
    project: Optional[str] = fields.none()
    rate_limit_window: Optional[int] = fields.none()
    rate_limit_count: Optional[int] = fields.none()

    @classmethod
    def from_config(cls, config: SentryConfig) -> "SentryDsn":
        return cls(
            dsn=config.dsn,
            project=config.project,
            rate_limit_window=config.rate_limit_window,
            rate_limit_count=config.rate_limit_count,
        )

    def merge(self, config: Union[SentryConfig, "SentryDsn"]):
        if not self.dsn:
            self.dsn = config.dsn
        if not self.project:
            self.project = config.project
        if not self.rate_limit_window:
            self.rate_limit_window = config.rate_limit_window
        if not self.rate_limit_count:
            self.rate_limit_count = config.rate_limit_count


@dataclass_json
@dataclass
class ApolloFederationSettings(JsonSchemaMixin):
    """Apollo Federation settings."""

    api_key: str
    graph: str
    graph_variant: str


@dataclass_json
@dataclass
class Component(JsonSchemaMixin):
    """Component configuration."""

    name: str
    variables: TerraformVariables = fields.dict_()
    secrets: TerraformVariables = fields.dict_()
    store_variables: StoreVariables = fields.dict_()
    store_secrets: StoreVariables = fields.dict_()
    health_check_path: Optional[str] = fields.none()
    sentry: Optional[SentryDsn] = fields.none()
    azure: Optional[ComponentAzureConfig] = fields.none()

    @property
    def definition(self) -> "ComponentConfig":
        return self._definition

    @definition.setter
    def definition(self, definition: "ComponentConfig"):
        self._definition = definition
        self.health_check_path = self.health_check_path or definition.health_check_path

    @property
    def integrations(self) -> List[str]:
        return self.definition.integrations

    @property
    def has_cloud_integration(self) -> bool:
        return "aws" in self.integrations or "azure" in self.integrations

    @property
    def endpoints(self) -> Dict[str, str]:
        return self.definition.endpoints


@dataclass_json
@dataclass
class SiteAWSSettings(JsonSchemaMixin):
    """Site-specific AWS settings."""

    account_id: int
    region: str
    deploy_role_name: Optional[str] = fields.none()
    default_tags: Tags = fields.dict_()
    extra_providers: Optional[List[AWSProvider]] = fields.list_()


@dataclass_json
@dataclass
class SiteAzureSettings(JsonSchemaMixin):
    """Site-specific Azure settings."""

    frontdoor: Optional[FrontdoorSettings] = fields.none()
    alert_group: Optional[AlertGroup] = fields.none()
    resource_group: str = ""
    tenant_id: str = ""  # Can overwrite values from AzureConfig
    subscription_id: str = ""  # Can overwrite values from AzureConfig
    region: str = ""  # Can overwrite values from AzureConfig
    service_object_ids: Dict[str, str] = fields.dict_()
    service_plans: Dict[str, ServicePlan] = fields.dict_()

    @classmethod
    def from_config(cls, config: AzureConfig):
        return cls(
            frontdoor=config.frontdoor,
            tenant_id=config.tenant_id,
            subscription_id=config.subscription_id,
            region=config.region,
            service_object_ids=config.service_object_ids,
            service_plans=config.service_plans,
        )

    def merge(self, config: AzureConfig):
        self.frontdoor = self.frontdoor or config.frontdoor
        self.tenant_id = self.tenant_id or config.tenant_id
        self.subscription_id = self.subscription_id or config.subscription_id
        self.region = self.region or config.region
        self.service_object_ids = self.service_object_ids or config.service_object_ids
        self.service_plans = {
            **(config.service_plans or {}),
            **(self.service_plans or {}),
        }


@dataclass_json
@dataclass
class Site(JsonSchemaMixin):
    """Site definition."""

    identifier: str
    endpoints: List[Endpoint] = field(
        default_factory=list,
        metadata=config(mm_field=fields.EndpointsField(), exclude=lambda x: not x),
    )
    commercetools: Optional[CommercetoolsSettings] = fields.none()
    contentful: Optional[ContentfulSettings] = fields.none()
    amplience: Optional[AmplienceSettings] = fields.none()
    apollo_federation: Optional[ApolloFederationSettings] = fields.none()
    azure: Optional[SiteAzureSettings] = fields.none()
    aws: Optional[SiteAWSSettings] = fields.none()
    components: List[Component] = fields.list_()
    sentry: Optional[SentryDsn] = fields.none()

    @property
    def cloud_components(self) -> List[Component]:
        """Return components with cloud platform integration."""
        return [c for c in self.components if c.has_cloud_integration]

    @property
    def used_endpoints(self) -> List[Endpoint]:
        """Return only the endpoints that are actually used by the components."""
        return [ep for ep in self.endpoints if ep.components]

    @property
    def has_cdn_endpoint(self) -> bool:
        """Check if there is an endpoint with a cdn enabled."""
        return any(ep.aws.enable_cdn for ep in self.used_endpoints)

    @property
    def used_custom_endpoints(self) -> List[Endpoint]:
        """Return custom endpoints that are used by the components."""
        return [ep for ep in self.used_endpoints if ep.url]

    @property
    def dns_zones(self) -> List[str]:
        return list({e.zone for e in self.used_endpoints if e.zone})

    @classmethod
    def json_schema(cls, *args, **kwargs):
        result = super().json_schema(*args, **kwargs)
        endpoints = result["Site"]["properties"]["endpoints"]
        result["Site"]["properties"]["endpoints"] = {
            "type": "object",
            "additionalProperties": endpoints["items"],
            "default": {},
        }
        return result
