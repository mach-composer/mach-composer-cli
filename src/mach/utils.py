import re
import unicodedata

PROTOCOL_RE = re.compile(r"^(http(s)?://)")


def strip_protocol(value: str) -> str:
    return PROTOCOL_RE.sub("", value)


def dns_zone_from_url(url: str) -> str:
    url = strip_protocol(url)
    parts = url.split("/")[0].split(".")
    if len(parts) == 2:
        raise ValueError("Given URL is already top-level domain")
    return ".".join(parts[1:])


def humanize_str(value: str) -> str:
    return re.sub(r"[-_]+", " ", value).title()


def slugify(value, allow_unicode=False, sep="_"):
    value = str(value)
    if allow_unicode:
        value = unicodedata.normalize("NFKC", value)
    else:
        value = (
            unicodedata.normalize("NFKD", value)
            .encode("ascii", "ignore")
            .decode("ascii")
        )
    value = re.sub(r"[^\w\s-]", "", value).strip().lower()
    return re.sub(r"[-\s]+", sep, value)
