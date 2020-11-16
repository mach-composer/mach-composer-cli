import os
import re
import unicodedata

import click
from cookiecutter.main import cookiecutter

COOKIECUTTER_TEMPLATE = "git@github.com:labd/mach-component-cookiecutter.git"


def create_component(output_dir: str):
    if output_dir and os.path.exists(output_dir):
        if not click.confirm(
            f"Directory {output_dir} already exists. Do you want to overwrite?"
        ):
            return

    cloud = click.prompt(
        "Cloud environment", type=click.Choice(["aws", "azure"]), default="aws"
    )
    language = click.prompt(
        "Language", type=click.Choice(["python", "node"]), default="python"
    )

    name = click.prompt("Name", default="example")
    description = click.prompt(
        "Description", default=f"{_humanize_str(name)} component"
    )
    short_name = click.prompt(
        "Short name", default=name.replace("_", "").replace("-", "")
    )
    component_identifier = click.prompt(
        "Component identifier", default=f"{_slugify(name, sep='-')}-component"
    )
    function_name = click.prompt("Function name", default=_slugify(name))

    context = {
        "language": language,
        "name": name,
        "description": description,
        "short_name": short_name,
        "component_identifier": component_identifier,
        "function_name": function_name,
    }

    if click.confirm("Use Sentry?", default=False):
        context["sentry_organization"] = click.prompt("Sentry Organization")
        context["sentry_project"] = click.prompt("Sentry Project")

    if cloud == "aws":
        context.update(
            {
                "lambda_s3_repository": click.prompt(
                    "Lambda repository S3 bucket", default="mach-lambda-repository"
                ),
            }
        )
    elif cloud == "azure":
        function_template = click.prompt(
            "Function template to use",
            type=click.Choice(["api-extension", "public-api"]),
            default="",
        )

        context.update(
            {
                "shared_resource_group": click.prompt("Shared resource group name"),
                "function_storage_account": click.prompt(
                    "Function storage account name"
                ),
                "function_container_name": click.prompt(
                    "Function container name", default="code"
                ),
                "function_template": function_template or "none",
            }
        )

    result = cookiecutter(
        COOKIECUTTER_TEMPLATE,
        directory=cloud,
        no_input=True,
        output_dir=output_dir or ".",
        extra_context=context,
    )

    dirname = os.path.basename(result)
    click.echo(f"New component {dirname} created 🎉")


def _humanize_str(value: str) -> str:
    return re.sub(r"[-_]+", " ", value).title()


def _slugify(value, allow_unicode=False, sep="_"):
    value = str(value)
    if allow_unicode:
        value = unicodedata.normalize("NFKC", value)
    else:
        value = (
            unicodedata.normalize("NFKD", value)
            .encode("ascii", "ignore")
            .decode("ascii")
        )
    value = re.sub(r"[^\w\s-]", "", value).strip().lower()
    return re.sub(r"[-\s]+", sep, value)
