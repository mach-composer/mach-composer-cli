from dataclasses import dataclass, field

from dataclasses_json import dataclass_json
from dataclasses_jsonschema import JsonSchemaMixin

__all__ = ["MachComposerConfig"]


@dataclass_json
@dataclass
class MachComposerConfig(JsonSchemaMixin):
    version: str
