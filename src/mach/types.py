from dataclasses import dataclass, field
from enum import Enum
from pathlib import Path
from typing import Any, Dict, List, Optional, Union

from dataclasses_json import config, dataclass_json
from dataclasses_jsonschema import FieldEncoder, JsonSchemaMixin
from mach import utils
from marshmallow import ValidationError, fields

TerraformVariables = Dict[str, Any]
StoreVariables = Dict[str, TerraformVariables]
StoreSecretVariables = Dict[str, StoreVariables]
LocalizedString = Dict[str, str]


class StringEnum(str, Enum):
    def __eq__(self, value):
        if isinstance(value, str):
            return self.value == value
        return super().__eq__(value)

    def __ne__(self, value):
        if isinstance(value, str):
            return self.value != value
        return super().__ne__(value)


# Define a none value as a custom dataclasses field so that
# null values get excluded in a dataclass dump
_none = lambda: field(default=None, metadata=config(exclude=lambda x: x is None))
_default = lambda value: field(
    default_factory=lambda: value, metadata=config(exclude=lambda x: x == value)
)
_list = lambda: field(default_factory=list, metadata=config(exclude=lambda x: not x))
_dict = lambda: field(default_factory=dict, metadata=config(exclude=lambda x: not x))


@dataclass_json
@dataclass
class Endpoint:
    url: str
    key: str = field(metadata=config(exclude=lambda x: True))
    zone: Optional[str] = _none()
    throttling_burst_limit: Optional[int] = _none()
    throttling_rate_limit: Optional[int] = _none()

    # To be set by the parser
    components: Optional[List["Component"]] = _list()

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

    def __post_init__(self):
        """Ensure endpoints have protocol stripped."""
        self.url = utils.strip_protocol(self.url)

        if not self.zone and self.url:
            try:
                self.zone = utils.dns_zone_from_url(self.url)
            except ValueError as e:
                raise ValidationError(f"Could not determine DNS zone: {e}")


class EndpointsField(fields.Dict):
    def _serialize(self, value, attr, obj, **kwargs):
        result = {}
        for endpoint in value:
            if endpoint.contains_defaults:
                result[endpoint.key] = endpoint.url
            else:
                result[endpoint.key] = endpoint.to_dict()

        return super()._deserialize(result, attr, obj, **kwargs)

    def _deserialize(self, value, attr, data, **kwargs):
        value = super()._deserialize(value, attr, data, **kwargs)
        result = []
        for k, v in value.items():
            if isinstance(v, str):
                result.append(Endpoint(key=k, url=v))
            elif isinstance(v, dict):
                v["key"] = k
                result.append(Endpoint.schema(infer_missing=True).load(v))
            else:
                raise ValidationError(f"Unexpected value found for endpoint {k}")

        return result


class EndpointEncoder(FieldEncoder):
    @property
    def json_schema(self):
        return {"type": "string"}


JsonSchemaMixin.register_field_encoders({Endpoint: EndpointEncoder()})


@dataclass_json
@dataclass
class AzureTFState(JsonSchemaMixin):
    """Azure storage account state backend configuration."""

    resource_group: str
    storage_account: str
    container_name: str
    state_folder: str


@dataclass_json
@dataclass
class AWSTFState(JsonSchemaMixin):
    """AWS S3 bucket state backend configuration."""

    bucket: str
    key_prefix: str
    role_arn: Optional[str] = _none()
    region: str = _default("eu-west-1")
    lock_table: Optional[str] = _none()
    encrypt = True


@dataclass_json
@dataclass
class TerraformProviders(JsonSchemaMixin):
    """Terraform provider version overwrites."""

    aws: Optional[str] = _none()
    azure: Optional[str] = _none()
    commercetools: Optional[str] = _none()
    sentry: Optional[str] = _none()
    contentful: Optional[str] = _none()


@dataclass_json
@dataclass
class TerraformConfig(JsonSchemaMixin):
    """Terraform configuration."""

    azure_remote_state: Optional[AzureTFState] = _none()
    aws_remote_state: Optional[AWSTFState] = _none()
    providers: Optional[TerraformProviders] = _none()


@dataclass_json
@dataclass
class SentryConfig(JsonSchemaMixin):
    """Global Sentry configuration."""

    dsn: Optional[str] = _none()
    rate_limit_window: Optional[int] = _none()
    rate_limit_count: Optional[int] = _none()

    auth_token: Optional[str] = _none()
    base_url: Optional[str] = _none()
    project: Optional[str] = _none()
    organization: Optional[str] = _none()

    @property
    def managed(self):
        """Indicate if the Sentry DSN should be managed by MACH."""
        return bool(self.auth_token)


@dataclass_json
@dataclass
class FrontDoorSettings(JsonSchemaMixin):
    """Frontdoor settings."""

    resource_group: str

    # Undocumented option to workaround some tenacious issues
    # with using Frontdoor in the Azure Terraform provider
    suppress_changes: bool = _default(False)


@dataclass_json
@dataclass
class AlertGroup(JsonSchemaMixin):
    """Alert group configuration."""

    name: str
    alert_emails: List[str] = _list()
    webhook_url: Optional[str] = _none()
    logic_app: Optional[str] = _none()

    @property
    def logic_app_name(self) -> Optional[str]:
        return self.logic_app.split(".")[1] if self.logic_app else None

    @property
    def logic_app_resource_group(self) -> Optional[str]:
        return self.logic_app.split(".")[0] if self.logic_app else None


@dataclass_json
@dataclass
class AzureConfig(JsonSchemaMixin):
    """Azure configuration."""

    tenant_id: str
    subscription_id: str
    region: str
    frontdoor: Optional[FrontDoorSettings] = _none()
    resources_prefix: Optional[str] = ""
    service_object_ids: Dict[str, str] = field(default_factory=dict)


@dataclass_json
@dataclass
class ContentfulConfig(JsonSchemaMixin):
    """Generic Contenful configuration."""

    cma_token: str
    organization_id: str


@dataclass_json
@dataclass
class AmplienceConfig(JsonSchemaMixin):
    """Generic Amplience configuration."""

    client_id: str
    client_secret: str


@dataclass_json
@dataclass
class AWSProvider(JsonSchemaMixin):
    """AWS provider configuration."""

    name: str
    region: str


class CloudOption(StringEnum):
    AWS = "aws"
    AZURE = "azure"


@dataclass_json
@dataclass
class GeneralConfig(JsonSchemaMixin):
    """Config this is shared across sites."""

    environment: str
    terraform_config: TerraformConfig
    cloud: CloudOption
    sentry: Optional[SentryConfig] = _none()
    azure: Optional[AzureConfig] = _none()
    contentful: Optional[ContentfulConfig] = _none()
    amplience: Optional[AmplienceConfig] = _none()


@dataclass_json
@dataclass
class ComponentConfig(JsonSchemaMixin):
    """Component definition."""

    name: str
    source: str
    version: str
    short_name: Optional[str] = _none()
    integrations: List[str] = _list()
    endpoints: Dict[str, str] = _dict()
    health_check_path: Optional[str] = _none()

    # Development options
    branch: Optional[str] = _none()
    package_script: Optional[str] = _none()
    package_filename: Optional[str] = _none()

    def __post_init__(self):
        """Ensure short_name is set."""
        self.short_name = self.short_name or self.name

    @property
    def use_version_reference(self):
        """Indicate if the module should be referenced with the version.

        This will be mainly used for development purposes when referring
        to a local directory; versioning is not possible, but we should
        still be able to define a version in our component for the actual
        function deployment itself.
        """
        return self.source.startswith("git")


@dataclass_json
@dataclass
class Store(JsonSchemaMixin):
    """commercetools store definition."""

    name: LocalizedString
    key: str
    languages: Optional[List[str]] = _list()
    distribution_channels: Optional[List[str]] = _list()
    supply_channels: Optional[List[str]] = _list()


@dataclass_json
@dataclass
class CommercetoolsChannel(JsonSchemaMixin):
    """commercetools channel definition."""

    key: str
    roles: List[str]
    name: Optional[LocalizedString] = _none()
    description: Optional[LocalizedString] = _none()


@dataclass_json
@dataclass
class CommercetoolsTax(JsonSchemaMixin):
    """commercetools tax definition."""

    country: str
    amount: float
    name: str


@dataclass_json
@dataclass
class CommercetoolsSettings(JsonSchemaMixin):
    """commercetools configuration."""

    project_key: str
    client_id: str
    client_secret: str
    scopes: str
    currencies: List[str]
    languages: List[str]
    countries: List[str]
    token_url: Optional[str] = _default(
        "https://auth.europe-west1.gcp.commercetools.com"
    )
    api_url: Optional[str] = _default("https://api.europe-west1.gcp.commercetools.com")
    # CT settings
    messages_enabled: Optional[bool] = _default(True)
    channels: Optional[List[CommercetoolsChannel]] = _list()
    taxes: Optional[List[CommercetoolsTax]] = _list()
    stores: List[Store] = _list()
    create_frontend_credentials: bool = _default(True)


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

    dsn: Optional[str] = _none()
    rate_limit_window: Optional[int] = _none()
    rate_limit_count: Optional[int] = _none()

    @classmethod
    def from_config(cls, config: SentryConfig) -> "SentryDsn":
        return cls(
            dsn=config.dsn,
            rate_limit_window=config.rate_limit_window,
            rate_limit_count=config.rate_limit_count,
        )

    def merge(self, config: Union[SentryConfig, "SentryDsn"]):
        if not self.dsn:
            self.dsn = config.dsn
        if not self.rate_limit_window:
            self.rate_limit_window = config.rate_limit_window
        if not self.rate_limit_count:
            self.rate_limit_count = config.rate_limit_count


@dataclass_json
@dataclass
class Component(JsonSchemaMixin):
    """Component configuration."""

    name: str
    variables: Optional[TerraformVariables] = _dict()
    secrets: Optional[TerraformVariables] = _dict()
    store_variables: Optional[StoreVariables] = _dict()
    store_secrets: Optional[StoreVariables] = _dict()
    short_name: Optional[str] = _none()
    health_check_path: Optional[str] = _none()
    sentry: Optional[SentryDsn] = _none()

    @property
    def definition(self) -> ComponentConfig:
        return self._definition

    @definition.setter
    def definition(self, definition: ComponentConfig):
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
    deploy_role_arn: Optional[str] = _none()
    extra_providers: Optional[List[AWSProvider]] = _list()


@dataclass_json
@dataclass
class SiteAzureSettings(JsonSchemaMixin):
    """Site-specific Azure settings."""

    service_object_ids: Dict[str, str] = field(default_factory=dict)
    frontdoor: Optional[FrontDoorSettings] = _none()
    alert_group: Optional[AlertGroup] = _none()
    resource_group: Optional[str] = ""
    tenant_id: Optional[str] = ""  # Can overwrite values from AzureConfig
    subscription_id: Optional[str] = ""  # Can overwrite values from AzureConfig
    region: Optional[str] = ""  # Can overwrite values from AzureConfig

    @classmethod
    def from_config(cls, config: AzureConfig):
        return cls(
            frontdoor=config.frontdoor,
            tenant_id=config.tenant_id,
            subscription_id=config.subscription_id,
            region=config.region,
            service_object_ids=config.service_object_ids,
        )

    def merge(self, config: AzureConfig):
        self.frontdoor = self.frontdoor or config.frontdoor
        self.tenant_id = self.tenant_id or config.tenant_id
        self.subscription_id = self.subscription_id or config.subscription_id
        self.region = self.region or config.region
        self.service_object_ids = self.service_object_ids or config.service_object_ids


@dataclass_json
@dataclass
class Site(JsonSchemaMixin):
    """Site definition."""

    identifier: str
    endpoints: Optional[List[Endpoint]] = field(
        default_factory=list,
        metadata=config(mm_field=EndpointsField(), exclude=lambda x: not x),
    )
    commercetools: Optional[CommercetoolsSettings] = _none()
    contentful: Optional[ContentfulSettings] = _none()
    amplience: Optional[AmplienceSettings] = _none()
    azure: Optional[SiteAzureSettings] = _none()
    aws: Optional[SiteAWSSettings] = _none()
    components: List[Component] = _list()
    sentry: Optional[SentryDsn] = _none()

    @property
    def used_endpoints(self) -> List[Endpoint]:
        """Return only the endpoints that are actually used by the components."""
        return [ep for ep in self.endpoints if ep.components]

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


@dataclass_json
@dataclass
class MachConfig(JsonSchemaMixin):
    """Main MACH configuration object."""

    general_config: GeneralConfig
    sites: List[Site]
    components: List[ComponentConfig] = _list()

    # Not used during a mach update or apply, but MACH
    # must be able to accept this attribute when parsing
    # encrypted configurations
    sops: Optional[Any] = _none()

    # Items that are not used in the configuration itself by set by the parser
    output_path: str = "deployments"
    file: Optional[str] = _none()

    @property
    def deployment_path(self) -> Path:
        return Path(self.output_path)

    def get_component(self, name: str) -> Optional[ComponentConfig]:
        for comp in self.components:
            if comp.name == name:
                return comp
        return None
