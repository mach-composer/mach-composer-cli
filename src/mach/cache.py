import os
from typing import TYPE_CHECKING

if TYPE_CHECKING:
    from mach.types import MachConfig


def cache_dir_for(config: "MachConfig", *, create=True):
    cache_dir = os.path.join(
        os.getcwd(), ".mach", os.path.splitext(os.path.basename(config.file))[0]
    )
    if create:
        os.makedirs(cache_dir, exist_ok=True)
    return cache_dir
