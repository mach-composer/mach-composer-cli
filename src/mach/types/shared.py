from dataclasses import dataclass
from typing import Optional

from dataclasses_json import dataclass_json
from dataclasses_jsonschema import JsonSchemaMixin

from . import fields

__all__ = ["ComponentAzureConfig", "ServicePlan"]


@dataclass_json
@dataclass
class ComponentAzureConfig(JsonSchemaMixin):
    service_plan: Optional[str] = fields.none()
    short_name: Optional[str] = fields.none()

    def merge(self, config: "ComponentAzureConfig"):
        self.service_plan = self.service_plan or config.service_plan
        self.short_name = self.short_name or config.short_name


@dataclass_json
@dataclass
class ServicePlan(JsonSchemaMixin):
    kind: str
    tier: str
    size: str
    capacity: Optional[int] = fields.none()
    dedicated_resource_group: bool = fields.default(False)
    per_site_scaling: bool = fields.default(False)
