import json

from mach import types

schema = types.MachConfig.json_schema()

del schema["properties"]["output_path"]
del schema["properties"]["file"]
del schema["properties"]["file_encrypted"]

with open("schema.json", 'w') as f:
    f.write(json.dumps(schema, indent=2))
