from dataclasses import field

from dataclasses_json import config
from marshmallow import ValidationError, fields

# Define a none value as a custom dataclasses field so that
# null values get excluded in a dataclass dump
none = lambda: field(default=None, metadata=config(exclude=lambda x: x is None))  # type: ignore
default = lambda value: field(
    default_factory=lambda: value, metadata=config(exclude=lambda x: x == value)  # type: ignore
)
list_ = lambda: field(default_factory=list, metadata=config(exclude=lambda x: not x))  # type: ignore # noqa
dict_ = lambda: field(default_factory=dict, metadata=config(exclude=lambda x: not x))  # type: ignore # noqa


class EndpointsField(fields.Dict):
    def _serialize(self, value, attr, obj, **kwargs):
        result = {}
        for endpoint in value:
            if endpoint.contains_defaults:
                result[endpoint.key] = endpoint.url
            else:
                result[endpoint.key] = endpoint.to_dict()

        return super()._deserialize(result, attr, obj, **kwargs)

    def _deserialize(self, value, attr, data, **kwargs):
        from .sites import Endpoint

        value = super()._deserialize(value, attr, data, **kwargs)
        result = []
        for k, v in value.items():
            if isinstance(v, str):
                result.append(Endpoint(key=k, url=v))
            elif isinstance(v, dict):
                v["key"] = k
                result.append(Endpoint.schema(infer_missing=True).load(v))  # type: ignore
            else:
                raise ValidationError(f"Unexpected value found for endpoint {k}")

        return result
