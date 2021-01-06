import os

import pytest
from mach.commands import bootstrap
from mach.parse import parse_config_from_file
from mach.validate import validate_config


# Add a time-out in case click expects more input
@pytest.mark.timeout(2)
@pytest.mark.parametrize("cloud", ["aws", "azure"])
def test_configuration(click_runner, click_dir, cloud):
    input_values = [
        "test",
        cloud,
        "unittest-nl",
        "y",  # Use commercetools?"
        "ct-unittest",
        "n",  # Use sentry?
        "n",
    ]
    result = click_runner.invoke(bootstrap, ["config"], input="\n".join(input_values))
    assert not result.exception

    file_path = os.path.join(click_dir, "main.yml")
    assert os.path.exists(file_path)

    config = parse_config_from_file(file_path)
    validate_config(config)

    config.general_config.environment == "test"
    assert not config.general_config.sentry
    config.sites[0].identifier == "unittest-nl"
    config.sites[0].commercetools.project_key == "ct-unittest"

    if cloud == "aws":
        assert not config.general_config.azure
        assert not config.general_config.terraform_config.azure_remote_state
        assert config.general_config.terraform_config.aws_remote_state
    elif cloud == "azure":
        assert config.general_config.azure
        assert config.general_config.terraform_config.azure_remote_state
        assert not config.general_config.terraform_config.aws_remote_state


# Add a time-out in case click expects more input
@pytest.mark.skipif(
    os.environ.get("CI"), reason="Disabled on CI for now (no Git access)"
)
@pytest.mark.timeout(5)
@pytest.mark.parametrize("language", ["python", "node"])
@pytest.mark.parametrize("cloud", ["aws", "azure"])
def test_component(cookiecutter, click_runner, click_dir, language, cloud):
    input_values = [
        cloud,
        language,
        "api-extensions",  # Name
        "API extensions component",  # Description
    ]

    if cloud == "aws":
        input_values += [
            "api-extensions-component",  # Directory name
            "y",  # use an HTTP endpoint?
            "y",  # support GraphQL?
            "y",  # use commercetools?
            "y",  # generate API extension?
            "y",  # generate API subscription?
            "y",  # token rotator?
            "n",  # Use Sentry?
            "mach-lambda-repository",  # s3 bucket
        ]

    elif cloud == "azure":
        input_values += [
            "apiext",  # Short name
            "api-extensions",  # Function name
            "api-extensions-component",  # Directory name
            "y",  # use an HTTP endpoint?
            "y",  # support GraphQL?
            "y",  # use commercetools?
            "y",  # generate API extension?
            "y",  # generate API subscription?
            "n",  # Use Sentry?
            "shared-rg",
            "sharedcomponentsrepo",
            "code",
        ]

    result = click_runner.invoke(
        bootstrap, ["component"], input="\n".join(input_values)
    )
    assert not result.exception

    component_path = os.path.join(click_dir, "api-extensions-component")
    assert os.path.exists(component_path)
