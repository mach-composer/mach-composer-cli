package aws

import "strings"

func StripProtocol(value string) string {
	if strings.HasPrefix(value, "http://") {
		return strings.TrimPrefix(value, "http://")
	}
	if strings.HasPrefix(value, "https://") {
		return strings.TrimPrefix(value, "https://")
	}
	return value
}
