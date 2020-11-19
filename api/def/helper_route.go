package def

import (
	"fmt"
	"net/url"
	"strings"
)

// replaceHTTPS replaces https:// with http:// in given string.
func replaceHTTPS(path string) string {
	return strings.ReplaceAll(path, "https://", "http://")
}

// convertRoutePathToInternal takes a route path and converts it to an internal url.
func convertRoutePathToInternal(path string, domainSufix string) (string, error) {
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	domainSufix = strings.Trim(strings.Trim(domainSufix, "."), "/")
	return fmt.Sprintf(
		"%s://%s.%s/%s",
		parsedURL.Scheme,
		strings.ReplaceAll(
			parsedURL.Hostname(), ".", "-",
		),
		domainSufix,
		strings.TrimLeft(parsedURL.Path, "/"),
	), nil
}
