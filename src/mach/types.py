import re
from dataclasses import dataclass, field
from enum import Enum
from pathlib import Path
from typing import Any, Dict, List, Optional

from dataclasses_json import dataclass_json

PROTOCOL_RE = re.compile(r"^(http(s)?://)")

TerraformVariables = Dict[str, Any]
StoreVariables = Dict[str, TerraformVariables]
StoreSecretVariables = Dict[str, StoreVariables]
LocalizedString = Dict[str, str]


@dataclass_json
@dataclass
class AzureTFState:
    resource_group: str
    storage_account: str
    container_name: str
    state_folder: str


@dataclass_json
@dataclass
class AWSTFState:
    bucket: str
    key_prefix: str
    role_arn: Optional[str] = None
    region: str = "eu-west-1"
    lock_table: Optional[str] = None
    encrypt = True


@dataclass_json
@dataclass
class TerraformConfig:
    azure_remote_state: Optional[AzureTFState] = None
    aws_remote_state: Optional[AWSTFState] = None


@dataclass_json
@dataclass
class SentryConfig:
    dsn: str


@dataclass_json
@dataclass
class FrontDoorSettings:
    resource_group: str
    dns_zone: str
    ssl_key_vault_name: str
    ssl_key_vault_secret_name: str
    ssl_key_vault_secret_version: str


@dataclass_json
@dataclass
class AlertGroup:
    name: str
    alert_emails: List[str] = field(default_factory=list)
    webhook_url: Optional[str] = None
    logic_app: Optional[str] = None

    @property
    def logic_app_name(self) -> Optional[str]:
        return self.logic_app.split(".")[1] if self.logic_app else None

    @property
    def logic_app_resource_group(self) -> Optional[str]:
        return self.logic_app.split(".")[0] if self.logic_app else None


@dataclass_json
@dataclass
class AzureConfig:
    front_door: Optional[FrontDoorSettings] = None
    resources_prefix: Optional[str] = ""
    tenant_id: Optional[str] = ""
    subscription_id: Optional[str] = ""
    region: Optional[str] = ""
    service_object_ids: Dict[str, str] = field(default_factory=dict)


@dataclass_json
@dataclass
class AWSConfig:
    code_repository: str


@dataclass_json
@dataclass
class ContentfulConfig:
    cma_token: str
    organization_id: str


@dataclass_json
@dataclass
class AWSProvider:
    name: str
    region: str


class CloudOption(Enum):
    AWS = "aws"
    AZURE = "azure"


@dataclass_json
@dataclass
class GeneralConfig:
    """Config this is shared across sites."""

    environment: str
    terraform_config: TerraformConfig
    cloud: CloudOption
    sentry: Optional[SentryConfig] = None
    azure: Optional[AzureConfig] = None
    aws: Optional[AWSConfig] = None
    contentful: Optional[ContentfulConfig] = None


@dataclass_json
@dataclass
class ComponentConfig:
    name: str
    source: str
    version: str
    short_name: Optional[str] = ""
    integrations: List[str] = field(default_factory=list)
    has_public_api: Optional[bool] = False
    health_check_path: Optional[str] = ""

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
class AzureProvider:
    resource_group: str


@dataclass_json
@dataclass
class Store:
    name: LocalizedString
    key: str
    languages: List[str] = field(default_factory=list)
    distribution_channels: List[str] = field(default_factory=list)


@dataclass_json
@dataclass
class CommercetoolsChannel:
    key: str
    roles: List[str]
    name: Optional[LocalizedString] = None
    description: Optional[LocalizedString] = None


@dataclass_json
@dataclass
class CommercetoolsTax:
    country: str
    amount: float
    name: str


@dataclass_json
@dataclass
class CommercetoolsSettings:
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
class ContentfulSettings:
    space: str
    default_locale: str = "en-US"
    cma_token: str = ""
    organization_id: str = ""

    def merge(self, config: ContentfulConfig):
        self.cma_token = self.cma_token or config.cma_token
        self.organization_id = self.organization_id or config.organization_id


@dataclass_json
@dataclass
class Component:
    name: str
    variables: TerraformVariables = field(default_factory=dict)
    secrets: TerraformVariables = field(default_factory=dict)
    short_name: Optional[str] = ""
    health_check_path: Optional[str] = ""

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
class SiteAWSSettings:
    account_id: int
    region: str
    deploy_role: Optional[str] = None
    api_gateway: Optional[str] = ""
    extra_providers: Optional[List[AWSProvider]] = field(default_factory=list)
    code_repository: Optional[str] = ""  # Can overwrite values from AWSConfig

    def merge(self, config: AWSConfig):
        self.code_repository = config.code_repository


@dataclass_json
@dataclass
class SiteAzureSettings:
    service_object_ids: Dict[str, str] = field(default_factory=dict)
    front_door: Optional[FrontDoorSettings] = None
    alert_group: Optional[AlertGroup] = None
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
class Site:
    identifier: str
    base_url: Optional[str] = ""
    commercetools: Optional[CommercetoolsSettings] = None
    contentful: Optional[ContentfulSettings] = None
    azure: Optional[SiteAzureSettings] = None
    aws: Optional[SiteAWSSettings] = None
    components: List[Component] = field(default_factory=list)

    @property
    def public_api_components(self) -> List[Component]:
        return [c for c in self.components if c.has_public_api]

    def __post_init__(self):
        """Ensure short_name is set."""
        self.base_url = PROTOCOL_RE.sub("", self.base_url)


@dataclass_json
@dataclass
class MachConfig:
    general_config: GeneralConfig
    sites: List[Site]
    components: List[ComponentConfig] = field(default_factory=list)
    output_path: str = "deployments"
    file: str = ""

    @property
    def deployment_path(self) -> Path:
        return Path(self.output_path)

    def get_component(self, name: str) -> Optional[ComponentConfig]:
        for comp in self.components:
            if comp.name == name:
                return comp
        return None
