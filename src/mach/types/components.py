from dataclasses import dataclass
from typing import Dict, List, Optional

from dataclasses_json import dataclass_json
from dataclasses_jsonschema import JsonSchemaMixin

from . import fields
from .shared import ComponentAzureConfig

__all__ = [
    "ComponentConfig",
]


@dataclass_json
@dataclass
class ComponentConfig(JsonSchemaMixin):
    """Component definition."""

    name: str
    source: str
    version: str
    short_name: Optional[str] = fields.none()
    integrations: List[str] = fields.list_()
    endpoints: Dict[str, str] = fields.dict_()
    health_check_path: Optional[str] = fields.none()

    # Azure-specific options
    azure: Optional[ComponentAzureConfig] = fields.none()

    # Development options
    branch: Optional[str] = fields.none()
    package_script: Optional[str] = fields.none()
    package_filename: Optional[str] = fields.none()

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
