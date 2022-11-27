package azure

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

// azureServicePlanResourceName Retrieve the resource name for a Azure app
// service plan.  The reason to make this conditional is because of backwards
// compatability; existing environments already have a `functionapp` resource.
// We want to keep that intact.
func azureServicePlanResourceName(value string) string {
	name := "funtionapps"
	if value != "default" {
		name = fmt.Sprintf("functionapps_%s", value)
	}
	return name
}

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

func SubdomainFromURL(value string) string {
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

	return domains[0]
}
