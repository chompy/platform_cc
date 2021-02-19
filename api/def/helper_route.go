/*
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
*/

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
func convertRoutePathToInternal(path string, domainSuffix string) (string, error) {
	parsedURL, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	domainSuffix = strings.Trim(strings.Trim(domainSuffix, "."), "/")
	subdomain := strings.TrimRight(strings.ReplaceAll(parsedURL.Hostname(), domainSuffix, ""), ".")
	return fmt.Sprintf(
		"%s://%s.%s/%s",
		parsedURL.Scheme,
		strings.ReplaceAll(subdomain, ".", "-"),
		domainSuffix,
		strings.TrimLeft(parsedURL.Path, "/"),
	), nil
}
