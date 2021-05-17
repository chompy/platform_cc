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

package project

import (
	"fmt"
	"strings"

	"gitlab.com/contextualcode/platform_cc/pkg/def"
)

// ReplaceDefault replaces {default} in given path with project id + domain suffix.
func (p *Project) ReplaceDefault(path string) string {
	out := strings.ReplaceAll(
		path,
		"{default}",
		fmt.Sprintf(
			"%s.default",
			p.ID,
		),
	)
	out = strings.ReplaceAll(
		out,
		"__PID__",
		p.ID,
	)
	return out
}

// RouteReplaceDefault replaces {default} values in all paths.
func (p *Project) RouteReplaceDefault(r def.Route) def.Route {
	r.OriginalURL = p.ReplaceDefault(r.OriginalURL)
	r.Path = p.ReplaceDefault(r.Path)
	r.To = p.ReplaceDefault(r.To)
	for i := range r.Redirects.Paths {
		r.Redirects.Paths[i].To = p.ReplaceDefault(r.Redirects.Paths[i].To)
	}
	for k, v := range r.Attributes {
		r.Attributes[k] = p.ReplaceDefault(v)
	}
	return r
}

// RoutesReplaceDefault replaces {default} values in all paths for multiple routes.
func (p *Project) RoutesReplaceDefault(routes []def.Route) []def.Route {
	for i := range routes {
		routes[i] = p.RouteReplaceDefault(routes[i])
	}
	return routes
}
