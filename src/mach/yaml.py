import re
from typing import Union

import yaml
import yamlinclude

EXTERNAL_RE = re.compile(r"^(git::)?(http|https)://")


class YamlIncludeConstructor(yamlinclude.YamlIncludeConstructor):
    def load(
        self,
        loader: Union[yaml.Loader, yaml.FullLoader],
        pathname: str,
        recursive: bool = False,
        encoding: str = "",
        reader: str = "",
    ):
        if EXTERNAL_RE.match(pathname):
            # TODO: Download file to .mach directory and include from there
            raise Exception("External include paths not supported yet")

        return super().load(loader, pathname, recursive, encoding, reader)
