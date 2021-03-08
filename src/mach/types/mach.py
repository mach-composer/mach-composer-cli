from dataclasses import dataclass

from dataclasses_json import dataclass_json
from dataclasses_jsonschema import JsonSchemaMixin

__all__ = ["MachComposerConfig"]


@dataclass_json
@dataclass
class MachComposerConfig(JsonSchemaMixin):
    version: str
