import json
import os


def get_file(name: str) -> str:
    """Return the absolute path to a test file."""
    return os.path.join(os.path.dirname(__file__), "files", name)


def get_file_content(name: str) -> str:
    with open(get_file(name)) as f:
        return f.read()


def get_json(name: str) -> dict:
    return json.loads(get_file_content(name))


def write_json(name: str, data: dict):
    with open(get_file(name), "w") as f:
        f.write(json.dumps(data, indent=2))
