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
    data = _clean_config_dump(data)
    content = yaml.dump(data, indent=2, explicit_start=True, sort_keys=False)
    with open(output_file, "w") as f:
        f.write(content)


def _create_config() -> types.MachConfig:
    environment = click.prompt("Environment", "test")
    cloud = click.prompt(
        "Cloud environment", type=click.Choice(["aws", "azure"]), default="aws"
    )
    site_id = click.prompt("Site identifier")
    use_commercetools = click.confirm("Use commercetools?", default=True)
    ct_project = ""
    if use_commercetools:
        ct_project = click.prompt("commercetools project name", default=site_id)
    use_sentry = click.confirm("Use Sentry?", default=False)
    use_contentful = click.confirm("Use Contentful?", default=False)

    integrations = []
    if use_commercetools:
        integrations.append("commercetools")
    if use_sentry:
        integrations.append("sentry")
    if use_contentful:
        integrations.append("contentful")
    # If we do have integrations, add the default (cloud) integration here as well
    if integrations:
        integrations = [cloud] + integrations

    if cloud == "aws":
        tf_config = types.TerraformConfig(
            aws_remote_state=types.AWSTFState(
                bucket="<your bucket>",
                key_prefix="mach",
            )
        )
    else:
        tf_config = types.TerraformConfig(
            azure_remote_state=types.AzureTFState(
                resource_group="<your-resource-group>",
                storage_account="<your-storage-account>",
                container_name="<your-container-name>",
                state_folder=environment,
            )
        )

    general_config_kwargs = dict(
        environment=environment,
        terraform_config=tf_config,
        cloud=cloud,
    )

    if use_sentry:
        general_config_kwargs["sentry"] = types.SentryConfig(
            auth_token="<your-auth-token>",
            project="<your-project>",
            organization="<your-organization>",
        )

    if use_contentful:
        general_config_kwargs["contentful"] = types.ContentfulConfig(
            cma_token="<your-cma-token>",
            organization_id="<your-organization-id>",
        )

    site = types.Site(
        identifier=site_id,
        components=[
            types.Component(
                name="your-component",
                variables={"FOO_VAR": "my-value"},
                secrets={"MY_SECRET": "secretvalue"},
            )
        ],
    )

    if use_commercetools:
        site.commercetools = types.CommercetoolsSettings(
            project_key=ct_project,
            client_id="<client-id>",
            client_secret="<client-secret>",
            scopes=f"manage_api_clients:{ct_project} manage_project:{ct_project} view_api_clients:{ct_project}",
        )

    return types.MachConfig(
        general_config=types.GeneralConfig(
            **general_config_kwargs,
        ),
        sites=[site],
        components=[
            types.ComponentConfig(
                name="your-component",
                source="git::https://github.com/<username>/<your-component>.git//terraform",
                version="0.1.0",
                integrations=integrations,
            )
        ],
    )


def _clean_config_dump(data: dict) -> dict:
    """Perform cleanup on the dump.

    TODO: These are actions that should be performed in the Marshmallow schema.
    """
    del data["output_path"]
    for component in data["components"]:
        if component["short_name"] == component["name"]:
            del component["short_name"]

    return data
