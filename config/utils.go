package config

import (
	"log"
	"net/url"
	"strings"
)

func StripProtocol(value string) string {
	if strings.HasPrefix(value, "http://") {
		return strings.TrimPrefix(value, "http://")
	}
	if strings.HasPrefix(value, "https://") {
		return strings.TrimPrefix(value, "https://")
	}
	return value
}

func ZoneFromURL(value string) string {
	u, err := url.Parse(value)
	if err != nil {
		log.Fatal(err)
	}

	var domains []string
	if !strings.Contains(value, "://") {
		parts := strings.SplitN(value, "/", 2)
		domains = strings.Split(parts[0], ".")
	} else {
		domains = strings.Split(u.Hostname(), ".")
	}

	if len(domains) < 3 {
		return strings.Join(domains, ".")
	} else {
		return strings.Join(domains[1:], ".")
	}
}
