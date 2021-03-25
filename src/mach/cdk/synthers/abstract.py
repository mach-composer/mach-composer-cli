import logging
from abc import ABC
from typing import TYPE_CHECKING, Any

if TYPE_CHECKING:
    from mach.cdk.stack import SiteStack
    from mach.types import MachConfig, Site


logger = logging.getLogger(__name__)


class Synther(ABC):
    stack: "SiteStack"

    def __init__(self, stack: "SiteStack"):
        self.stack = stack

    def synth(self, obj):
        pass

    @property
    def config(self) -> "MachConfig":
        return self.stack.config

    @property
    def site(self) -> "Site":
        return self.stack.site

    def subsynth(self, cls: "Synther", obj: Any):
        if not obj:
            logger.debug("No object to synth. Will skip")
            return

        cls(self.stack).synth(obj)
