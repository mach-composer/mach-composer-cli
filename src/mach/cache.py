import os


def cache_dir_for(file: str, *, create=True):
    cache_dir = os.path.join(
        os.getcwd(), ".mach", os.path.splitext(os.path.basename(file))[0]
    )
    if create:
        os.makedirs(cache_dir, exist_ok=True)
    return cache_dir
