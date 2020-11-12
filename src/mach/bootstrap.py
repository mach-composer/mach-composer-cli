import os

import click
import yaml
from mach import types


def create_configuration(output_file: str):
    if os.path.exists(output_file):
        if not click.confirm(
            f"File {output_file} already exists. Do you want to overwrite?"
        ):
            return

    config = _create_config()
    data = config.to_dict()
    del data["output_path"]
    content = yaml.dump(data, explicit_start=True, sort_keys=False)
    with open(output_file, "w") as f:
        f.write(content)


def _create_config() -> types.MachConfig:
    environment = click.prompt("Environment", "test")
    cloud = click.prompt(
        "Cloud environment", type=click.Choice(["aws", "azure"]), default="aws"
    )
    site = click.prompt("Site identifier")

    return types.MachConfig(
        general_config=types.GeneralConfig(
            environment=environment,
            terraform_config=types.TerraformConfig(
                aws_remote_state=types.AWSTFState(
                    bucket="<your bucket>",
                    key_prefix="mach",
                )
            ),
            cloud=cloud,
            # sentry=types.SentryConfig(dsn="sentry-dsn"),
        ),
        sites=[
            types.Site(
                identifier=site,
                components=[
                    types.Component(
                        name="your-component",
                        variables={"FOO_VAR": "my-value"},
                        secrets={"MY_SECRET": "secretvalue"},
                    )
                ],
            ),
        ],
        components=[
            types.ComponentConfig(
                name="your-component",
                source="git::https://github.com/<username>/<your-component>.git//terraform",
                version="0.1.0",
            )
        ],
    )
