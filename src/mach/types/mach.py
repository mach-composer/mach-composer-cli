from dataclasses import dataclass
from typing import Optional

from dataclasses_json import dataclass_json
from dataclasses_jsonschema import JsonSchemaMixin

__all__ = ["MachComposerConfig"]


@dataclass_json
@dataclass
class MachComposerConfig(JsonSchemaMixin):
    version: str
    variables_file: Optional[str] = None
