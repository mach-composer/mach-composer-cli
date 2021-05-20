import pytest
from mach import variables


def test_resolve_variables():
    vars = {"my-value": "foo", "secrets": {"site1": {"my-value": "bar"}}}

    variables.resolve_variable("my-value", vars) == "foo"
    with pytest.raises(variables.VariableNotFound):
        variables.resolve_variable("my-other-value", vars)

    variables.resolve_variable("secrets.site1.my-value", vars) == "bar"
    with pytest.raises(variables.VariableNotFound):
        variables.resolve_variable("secrets.site2.my-value", vars)
