from typing import TYPE_CHECKING

from cdktf import App, TerraformStack

from .synthers import SiteSynther

if TYPE_CHECKING:
    from .types import MachConfig, Site


__all__ = ["generate"]


class SiteStack(TerraformStack):
    config: "MachConfig"
    site: "Site"

    def __init__(self, app, *, config: "MachConfig", site: "Site"):
        super().__init__(app, site.identifier)
        self.config = config
        self.site = site

        synther = SiteSynther(self)
        synther.synth(site)


def generate(config: "MachConfig", site: "Site", outdir: str):
    app = App(outdir=outdir)
    SiteStack(app, config=config, site=site)
    app.synth()
