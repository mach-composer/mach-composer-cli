import os
from typing import Any, Dict

import click
from cookiecutter.main import cookiecutter
from mach import utils

DEFAULT_COOKIECUTTER = "git@github.com:labd/mach-component-cookiecutter.git"


def create_component(output_dir: str, cookiecutter_location: str):
    """Create a component using the given cookiecutter template or the default one."""
    cookiecutter_kwargs: Dict[str, Any] = {
        "output_dir": output_dir or ".",
    }

    if output_dir and os.path.exists(output_dir):
        if not click.confirm(
            f"Directory {output_dir} already exists. Do you want to overwrite?"
        ):
            return

    if not cookiecutter_location:
        cookiecutter_location = DEFAULT_COOKIECUTTER
        cloud = click.prompt(
            "Cloud environment", type=click.Choice(["aws", "azure"]), default="aws"
        )

        cookiecutter_kwargs.update(
            {
                "directory": cloud,
                "extra_context": _get_component_context(cloud),
                "no_input": True,
            }
        )

    result = cookiecutter(
        cookiecutter_location,
        **cookiecutter_kwargs,
    )

    dirname = os.path.basename(result)
    click.echo(f"New component {dirname} created 🎉")


def _get_component_context(cloud: str) -> dict:
    """Prompt user for input for the cookiecutter."""
    # We have a different default language here because the cookiecutter
    # doesnt contain examples for all languages yet.
    default_lang = "node" if cloud == "aws" else "python"

    language = click.prompt(
        "Language", type=click.Choice(["python", "node"]), default=default_lang
    )

    name = click.prompt("Name", default="example-name")
    description = click.prompt(
        "Description", default=f"{utils.humanize_str(name)} component"
    )

    dirname = click.prompt(
        "Directory name", default=f"{utils.slugify(name, sep='-')}-component"
    )

    context = {
        "language": language,
        "name": name,
        "description": description,
        "component_identifier": dirname,
        "function_name": utils.slugify(name),
        "use_public_api": 0,
        "include_graphql_server": 0,
        "use_commercetools": 0,
        "use_commercetools_api_extension": 0,
        "use_commercetools_subscription": 0,
        "use_commercetools_token_rotator": 0,
    }

    if click.confirm("Uses an HTTP endpoint?", default=True):
        context["use_public_api"] = 1
        if cloud == "azure":
            context["function_name"] = click.prompt(
                "Function name", default=utils.slugify(name)
            )
        context["include_graphql_server"] = int(
            click.confirm("Include GraphQL support?", default=True)
        )

    if click.confirm("Uses commercetools?", default=True):
        context["use_commercetools"] = 1
        context["use_commercetools_api_extension"] = int(
            click.confirm("Generate commercetools API extension?", default=True)
        )
        context["use_commercetools_subscription"] = int(
            click.confirm("Generate commercetools Subcription?", default=True)
        )
        if cloud == "aws":
            context["use_commercetools_token_rotator"] = int(
                click.confirm(
                    "Do you want use the commercetools token rotator component?",
                    default=False,
                )
            )

    if click.confirm("Use Sentry?", default=False):
        context["sentry_organization"] = click.prompt("Sentry Organization")
        context["sentry_project"] = click.prompt("Sentry Project")

    if cloud == "aws":
        context["lambda_s3_repository"] = click.prompt("Lambda repository S3 bucket")
    elif cloud == "azure":
        context.update(
            {
                "shared_resource_group": click.prompt("Shared resource group name"),
                "function_storage_account": click.prompt(
                    "Function storage account name"
                ),
                "function_container_name": click.prompt(
                    "Function container name", default="code"
                ),
            }
        )

    return context
