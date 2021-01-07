import re
import unicodedata
from typing import Tuple

import tldextract

PROTOCOL_RE = re.compile(r"^(http(s)?://)")


def strip_protocol(value: str) -> str:
    return PROTOCOL_RE.sub("", value)


def dns_zone_from_url(url: str) -> str:
    return domain_parts_from_url(url)[1]


def subdomain_from_url(url: str) -> str:
    return domain_parts_from_url(url)[0]


def domain_parts_from_url(url: str) -> Tuple[str, str]:
    ext = tldextract.extract(url)
    if not ext.subdomain:
        raise ValueError("Given URL is already top-level domain")

    sd_parts = ext.subdomain.split(".")
    parts = sd_parts[1:] + [ext.domain, ext.suffix]
    return sd_parts[0], ".".join(parts)


def humanize_str(value: str) -> str:
    return re.sub(r"[-_]+", " ", value).title()


def slugify(value, sep="_"):
    value = str(value)
    value = (
        unicodedata.normalize("NFKD", value).encode("ascii", "ignore").decode("ascii")
    )
    value = re.sub(r"[^\w\s-]", "", value).strip().lower()
    return re.sub(r"[-\s]+", sep, value)
