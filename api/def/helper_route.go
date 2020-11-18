package def

import (
	"fmt"
	"net/url"
	"strings"
)

// RouterRootDomain defines domain suffix for internal routes.
const RouterRootDomain = ".pcc.localtest.me/"

// replaceHTTPS replaces https:// with http:// in given string.
func replaceHTTPS(path string) string {
	return strings.ReplaceAll(path, "https://", "http://")
}

// convertRoutePathToInternal takes a route path and converts it to an internal url.
func convertRoutePathToInternal(path string) (string, error) {
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s://%s%s%s",
		parsedURL.Scheme,
		strings.ReplaceAll(
			parsedURL.Hostname(), ".", "-",
		),
		RouterRootDomain,
		strings.TrimLeft(parsedURL.Path, "/"),
	), nil
}
