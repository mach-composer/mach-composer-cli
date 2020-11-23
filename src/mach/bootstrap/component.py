import os
import re
import unicodedata

import click
from cookiecutter.main import cookiecutter


def create_component(output_dir: str, cookiecutter_location: str):

    COOKIECUTTER_TEMPLATE = cookiecutter_location

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

    use_public_api = click.confirm("Generate public API?", default=True)
    use_commercetools_api_extension = click.confirm(
        "Generate commercetools API extension?", default=True
    )
    use_commercetools_subscription = click.confirm(
        "Generate commercetools Subcription?", default=True
    )

    context["use_public_api"] = 1 if use_public_api else 0
    context["use_commercetools_api_extension"] = 1 if use_commercetools_api_extension else 0
    context["use_commercetools_subscription"] = 1 if use_commercetools_subscription else 0

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
