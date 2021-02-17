from dataclasses import dataclass
from typing import Optional

from dataclasses_json import dataclass_json
from dataclasses_jsonschema import JsonSchemaMixin

from . import fields

__all__ = ["ComponentAzureConfig"]


@dataclass_json
@dataclass
class ComponentAzureConfig(JsonSchemaMixin):
    service_plan: Optional[str] = fields.none()
