from mach.exceptions import MachError


class VariableNotFound(MachError):
    def __init__(self, var_name: str):
        super().__init__(f"Variable {var_name} not found in variables")


def resolve_variable(var, variables):
    lookup, *remain = var.split(".", maxsplit=1)

    try:
        value = variables[lookup]
    except KeyError:
        raise VariableNotFound(var)

    if remain:
        try:
            return resolve_variable(remain[0], value)
        except VariableNotFound:
            raise VariableNotFound(var)
    return value
