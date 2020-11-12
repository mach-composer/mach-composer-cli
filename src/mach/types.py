import re
from dataclasses import dataclass, field
from enum import Enum
from pathlib import Path
from typing import Any, Dict, List, Optional, Union

from dataclasses_json import config, dataclass_json
from dataclasses_jsonschema import JsonSchemaMixin

PROTOCOL_RE = re.compile(r"^(http(s)?://)")

TerraformVariables = Dict[str, Any]
StoreVariables = Dict[str, TerraformVariables]
StoreSecretVariables = Dict[str, StoreVariables]
LocalizedString = Dict[str, str]


# Define a none value as a custom dataclasses field so that
# null values get excluded in a dataclass dump
_none = lambda: field(default=None, metadata=config(exclude=lambda x: x is None))


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
    region: str = "eu-west-1"
    lock_table: Optional[str] = _none()
    encrypt = True


@dataclass_json
@dataclass
class TerraformConfig(JsonSchemaMixin):
    """Terraform configuration."""

    azure_remote_state: Optional[AzureTFState] = _none()
    aws_remote_state: Optional[AWSTFState] = _none()


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
    dns_zone: str
    ssl_key_vault_name: str
    ssl_key_vault_secret_name: str
    ssl_key_vault_secret_version: str


@dataclass_json
@dataclass
class AlertGroup(JsonSchemaMixin):
    """Alert group configuration."""

    name: str
    alert_emails: List[str] = field(default_factory=list)
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

    front_door: Optional[FrontDoorSettings] = _none()
    resources_prefix: Optional[str] = ""
    tenant_id: Optional[str] = ""
    subscription_id: Optional[str] = ""
    region: Optional[str] = ""
    service_object_ids: Dict[str, str] = field(default_factory=dict)


@dataclass_json
@dataclass
class ContentfulConfig(JsonSchemaMixin):
    """Generic Contenful configuration."""

    cma_token: str
    organization_id: str


@dataclass_json
@dataclass
class AWSProvider(JsonSchemaMixin):
    """AWS provider configuration."""

    name: str
    region: str


class CloudOption(Enum):
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


@dataclass_json
@dataclass
class ComponentConfig(JsonSchemaMixin):
    """Component definition."""

    name: str
    source: str
    version: str
    short_name: Optional[str] = _none()
    integrations: List[str] = field(default_factory=list)
    has_public_api: Optional[bool] = False
    health_check_path: Optional[str] = _none()
    branch: Optional[str] = _none()

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
    languages: List[str] = field(default_factory=list)
    distribution_channels: List[str] = field(default_factory=list)


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
    token_url: Optional[str] = "https://auth.europe-west1.gcp.commercetools.com"
    api_url: Optional[str] = "https://api.europe-west1.gcp.commercetools.com"
    # CT settings
    currencies: List[str] = field(default_factory=list)
    languages: List[str] = field(default_factory=list)
    countries: List[str] = field(default_factory=list)
    messages_enabled: Optional[bool] = True
    channels: Optional[List[CommercetoolsChannel]] = field(default_factory=list)
    taxes: Optional[List[CommercetoolsTax]] = field(default_factory=list)
    stores: List[Store] = field(default_factory=list)
    create_frontend_credentials: bool = True


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
class SentryDsn(JsonSchemaMixin):
    """Specific sentry DSN settings."""

    dsn: Optional[str] = _none()
    rate_limit_window: Optional[int] = _none()
    rate_limit_count: Optional[int] = _none()

    @classmethod
    def from_config(cls, config: SentryConfig) -> "SentryDsn":
        return cls(dsn=config.dsn)

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
    variables: TerraformVariables = field(default_factory=dict)
    secrets: TerraformVariables = field(default_factory=dict)
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
    def is_software_component(self) -> bool:
        return "aws" in self.integrations or "azure" in self.integrations

    @property
    def has_public_api(self):
        return self.definition.has_public_api


@dataclass_json
@dataclass
class SiteAWSSettings(JsonSchemaMixin):
    """Site-specific AWS settings."""

    account_id: int
    region: str
    deploy_role: Optional[str] = _none()
    extra_providers: Optional[List[AWSProvider]] = field(default_factory=list)
    route53_zone_name: Optional[str] = ""


@dataclass_json
@dataclass
class SiteAzureSettings(JsonSchemaMixin):
    """Site-specific Azure settings."""

    service_object_ids: Dict[str, str] = field(default_factory=dict)
    front_door: Optional[FrontDoorSettings] = _none()
    alert_group: Optional[AlertGroup] = _none()
    resource_group: Optional[str] = ""
    tenant_id: Optional[str] = ""  # Can overwrite values from AzureConfig
    subscription_id: Optional[str] = ""  # Can overwrite values from AzureConfig
    region: Optional[str] = ""  # Can overwrite values from AzureConfig

    @classmethod
    def from_config(cls, config: AzureConfig):
        return cls(
            front_door=config.front_door,
            tenant_id=config.tenant_id,
            subscription_id=config.subscription_id,
            region=config.region,
            service_object_ids=config.service_object_ids,
        )

    def merge(self, config: AzureConfig):
        self.front_door = self.front_door or config.front_door
        self.tenant_id = self.tenant_id or config.tenant_id
        self.subscription_id = self.subscription_id or config.subscription_id
        self.region = self.region or config.region
        self.service_object_ids = self.service_object_ids or config.service_object_ids


@dataclass_json
@dataclass
class Site(JsonSchemaMixin):
    """Site definition."""

    identifier: str
    base_url: Optional[str] = _none()
    commercetools: Optional[CommercetoolsSettings] = _none()
    contentful: Optional[ContentfulSettings] = _none()
    azure: Optional[SiteAzureSettings] = _none()
    aws: Optional[SiteAWSSettings] = _none()
    components: List[Component] = field(default_factory=list)
    sentry: Optional[SentryDsn] = _none()

    @property
    def public_api_components(self) -> List[Component]:
        return [c for c in self.components if c.has_public_api]

    def __post_init__(self):
        """Ensure base_url has protocol stripped"""
        if self.base_url:
            self.base_url = PROTOCOL_RE.sub("", self.base_url)


@dataclass_json
@dataclass
class MachConfig(JsonSchemaMixin):
    """Main MACH configuration object."""

    general_config: GeneralConfig
    sites: List[Site]
    components: List[ComponentConfig] = field(default_factory=list)
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
