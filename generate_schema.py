import json

from mach import types

schema = types.MachConfig.json_schema()
with open("schema.json", 'w') as f:
    f.write(json.dumps(schema, indent=2))
