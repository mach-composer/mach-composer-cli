from dataclasses import dataclass
from enum import Enum
from typing import Dict, Optional

from dataclasses_json import dataclass_json
from dataclasses_jsonschema import JsonSchemaMixin

from . import fields
from .shared import ServicePlan

__all__ = [
    "AWSTFState",
    "AmplienceConfig",
    "AzureConfig",
    "AzureTFState",
    "CloudOption",
    "ContentfulConfig",
    "FrontdoorSettings",
    "GlobalConfig",
    "SentryConfig",
    "TerraformConfig",
    "TerraformProviders",
]


class StringEnum(str, Enum):
    def __eq__(self, value):
        if isinstance(value, str):
            return self.value == value
        return super().__eq__(value)

    def __ne__(self, value):
        if isinstance(value, str):
            return self.value != value
        return super().__ne__(value)


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
    region: str
    role_arn: Optional[str] = fields.none()
    lock_table: Optional[str] = fields.none()
    encrypt = True


@dataclass_json
@dataclass
class TerraformProviders(JsonSchemaMixin):
    """Terraform provider version overwrites."""

    aws: Optional[str] = fields.none()
    azure: Optional[str] = fields.none()
    commercetools: Optional[str] = fields.none()
    sentry: Optional[str] = fields.none()
    contentful: Optional[str] = fields.none()
    amplience: Optional[str] = fields.none()


@dataclass_json
@dataclass
class TerraformConfig(JsonSchemaMixin):
    """Terraform configuration."""

    azure_remote_state: Optional[AzureTFState] = fields.none()
    aws_remote_state: Optional[AWSTFState] = fields.none()
    providers: Optional[TerraformProviders] = fields.none()


@dataclass_json
@dataclass
class SentryConfig(JsonSchemaMixin):
    """Global Sentry configuration."""

    dsn: Optional[str] = fields.none()
    rate_limit_window: Optional[int] = fields.none()
    rate_limit_count: Optional[int] = fields.none()

    auth_token: Optional[str] = fields.none()
    base_url: Optional[str] = fields.none()
    project: Optional[str] = fields.none()
    organization: Optional[str] = fields.none()

    @property
    def managed(self):
        """Indicate if the Sentry DSN should be managed by MACH."""
        return bool(self.auth_token)


@dataclass_json
@dataclass
class FrontdoorSslConfig(JsonSchemaMixin):
    name: str
    resource_group: str
    secret_name: str


@dataclass_json
@dataclass
class FrontdoorSettings(JsonSchemaMixin):
    """Frontdoor settings."""

    dns_resource_group: str
    ssl_key_vault: Optional[FrontdoorSslConfig] = fields.none()
    # Undocumented option to workaround some tenacious issues
    # with using Frontdoor in the Azure Terraform provider
    suppress_changes: bool = fields.default(False)


@dataclass_json
@dataclass
class AzureConfig(JsonSchemaMixin):
    """Azure configuration."""

    tenant_id: str
    subscription_id: str
    region: str
    frontdoor: Optional[FrontdoorSettings] = fields.none()
    resources_prefix: str = ""
    service_object_ids: Dict[str, str] = fields.dict_()
    service_plans: Dict[str, ServicePlan] = fields.dict_()


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


class CloudOption(StringEnum):
    AWS = "aws"
    AZURE = "azure"


@dataclass_json
@dataclass
class GlobalConfig(JsonSchemaMixin):
    """Config that is shared across sites."""

    environment: str
    terraform_config: TerraformConfig
    cloud: CloudOption
    sentry: Optional[SentryConfig] = fields.none()
    azure: Optional[AzureConfig] = fields.none()
    contentful: Optional[ContentfulConfig] = fields.none()
    amplience: Optional[AmplienceConfig] = fields.none()
