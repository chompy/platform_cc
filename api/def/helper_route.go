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
