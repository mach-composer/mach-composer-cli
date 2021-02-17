"""Collection of test utilities specifically to test out terraform cases."""
import os
from typing import List

from mach import terraform, types

from tests.utils import HclWrapper, load_hcl


def generate(config: types.MachConfig) -> HclWrapper:
    terraform.generate_terraform(config)
    file_path = os.path.join(config.output_path, config.sites[0].identifier, "site.tf")
    assert os.path.exists(file_path), "No site.tf file found"
    return load_hcl(file_path)


def get_resource_ids(data: HclWrapper) -> List[str]:
    """Get all resource ids in <resource-type>.<name> format."""
    result = []
    for type_, resources in data.resource.items():
        for key, resource in resources.items():
            result.append(f"{type_}.{key}")
    return sorted(result)


def get_data_ids(data: HclWrapper) -> List[str]:
    """Get all data ids in <resource-type>.<name> format."""
    result = []
    for type_, data_ in data.data.items():
        for key in data_.keys():
            result.append(f"{type_}.{key}")
    return sorted(result)
