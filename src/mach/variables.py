from contextlib import contextmanager

from mach.exceptions import MachError

_ignore_var_not_found = False


class VariableNotFound(MachError):
    def __init__(self, var_name: str):
        super().__init__(f"Variable {var_name} not found in variables")


def resolve_variable(var, variables):
    try:
        return _resolve_variable(var, variables)
    except VariableNotFound:
        if _ignore_var_not_found:
            return ""
        raise


def _resolve_variable(var, variables):
    lookup, *remain = var.split(".", maxsplit=1)

    if isinstance(variables, list):
        try:
            lookup = int(lookup)
        except ValueError:
            raise VariableNotFound("List indicies needs a number to index")
    elif not isinstance(variables, dict):
        # We've reached the end-node which is just
        raise VariableNotFound(var)

    try:
        value = variables[lookup]
    except KeyError:
        raise VariableNotFound(var)

    if remain:
        try:
            return _resolve_variable(remain[0], value)
        except VariableNotFound:
            raise VariableNotFound(var)
    return value


@contextmanager
def ignore_variable_not_found():
    global _ignore_var_not_found
    _ignore_var_not_found = True
    yield
    _ignore_var_not_found = False
