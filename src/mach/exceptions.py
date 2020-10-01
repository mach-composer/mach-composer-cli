class MachError(Exception):
    pass


class UpdateError(MachError):
    pass


class ParseError(MachError):
    def __init__(self, msg, details=None):
        super().__init__(msg)
        self.details = details

    def __str__(self):
        result = super().__str__()
        details = self.details_pretty
        if details:
            result += f"{details}"

        return result

    @property
    def details_pretty(self):
        if not self.details or not isinstance(self.details, dict):
            return ""

        return self._pretty(self.details)

    def _pretty(self, data, *, indent=0):
        if not isinstance(data, dict):
            return data

        return "\n" + "\n".join(
            [
                f"{' ' * indent}{k}: {self._pretty(v, indent=indent+2)}"
                for k, v in data.items()
            ]
        )


class ValidationError(MachError):
    pass
