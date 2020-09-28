from dataclasses import dataclass, field
from enum import Enum
from pathlib import Path
from typing import Any, Dict, List, Optional

from dataclasses_json import dataclass_json

TerraformVariables = Dict[str, Any]
StoreVariables = Dict[str, TerraformVariables]
StoreSecretVariables = Dict[str, StoreVariables]
LocalizedString = Dict[str, str]


@dataclass_json
@dataclass
class AzureTFState:
    resource_group_name: str
    storage_account_name: str
    container_name: str
    state_folder: str


@dataclass_json
@dataclass
class TerraformConfig:
    azure_remote_state: AzureTFState


@dataclass_json
@dataclass
class SentryConfig:
    dsn: str


@dataclass_json
@dataclass
class FrontDoorSettings:
    resource_group_name: str
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


class Environment(Enum):
    DEV = "development"
    TEST = "test"
    PROD = "production"

    def __str__(self):
        return self.value


@dataclass_json
@dataclass
class GeneralConfig:
    """Config this is shared across sites."""

    environment: Environment
    terraform_config: TerraformConfig
    sentry: Optional[SentryConfig] = None
    azure: Optional[AzureConfig] = None


@dataclass_json
@dataclass
class ComponentConfig:
    name: str
    source: str
    version: str
    short_name: str = ""

    def __post_init__(self):
        """Ensure short_name is set."""
        self.short_name = self.short_name or self.name


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
class Component:
    name: str
    variables: TerraformVariables = field(default_factory=dict)
    secrets: TerraformVariables = field(default_factory=dict)
    is_software_component: Optional[bool] = True
    has_public_api: Optional[bool] = False
    health_check_path: Optional[str] = ""
    short_name: Optional[str] = ""

    @property
    def definition(self) -> ComponentConfig:
        return self._definition

    @definition.setter
    def definition(self, definition: ComponentConfig):
        self._definition = definition


@dataclass_json
@dataclass
class SiteAzureSettings:
    service_object_ids: Dict[str, str] = field(default_factory=dict)
    front_door: Optional[FrontDoorSettings] = None
    alert_group: Optional[AlertGroup] = None
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
        )

    def merge(self, config: AzureConfig):
        self.front_door = self.front_door or config.front_door
        self.tenant_id = self.tenant_id or config.tenant_id
        self.subscription_id = self.subscription_id or config.subscription_id
        self.region = self.region or config.region


@dataclass_json
@dataclass
class Site:
    identifier: str
    commercetools: Optional[CommercetoolsSettings] = None
    azure: Optional[SiteAzureSettings] = None
    components: List[Component] = field(default_factory=list)

    @property
    def public_api_components(self):
        return [c for c in self.components if c.has_public_api]


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
