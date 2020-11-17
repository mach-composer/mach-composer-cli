import json
import os


def get_file(name):
    """Return the absolute path to a test file."""
    return os.path.join(os.path.dirname(__file__), "files", name)


def get_file_content(name):
    with open(get_file(name)) as f:
        return f.read()


def get_json(name):
    return json.loads(get_file_content(name))


def write_json(name, content):
    with open(get_file(name), "w") as f:
        f.write(content)
