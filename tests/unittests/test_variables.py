import pytest
from mach import variables


def test_resolve_variables():
    vars = {
        "my-value": "foo",
        "secrets": {"site1": {"my-value": "bar"}},
        "list": ["one", "two", {"nested-key": "three"}],
    }

    variables.resolve_variable("my-value", vars) == "foo"
    with pytest.raises(variables.VariableNotFound):
        variables.resolve_variable("my-other-value", vars)

    variables.resolve_variable("secrets.site1.my-value", vars) == "bar"
    with pytest.raises(variables.VariableNotFound):
        variables.resolve_variable("secrets.site2.my-value", vars)

    variables.resolve_variable("list.0", vars) == "one"
    variables.resolve_variable("list.1", vars) == "two"
    variables.resolve_variable("list.2", vars) == {"nested-key": "three"}
    variables.resolve_variable("list.2.nested-key", vars) == "three"

    with pytest.raises(variables.VariableNotFound):
        variables.resolve_variable("my-value.string-attribute", vars)
