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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"
	"gopkg.in/yaml.v2"
)

const defaultPath = "__PID__.default"

// Route defines a route.
type Route struct {
	Path        string            `json:"-"`
	Type        string            `yaml:"type" json:"type"`
	Upstream    string            `yaml:"upstream" json:"upstream"`
	To          string            `yaml:"to" json:"to"`
	ID          string            `yaml:"id" json:"id"`
	Attributes  map[string]string `yaml:"attributes" json:"attributes"`
	Cache       RouteCache        `yaml:"cache" json:"-"`
	Redirects   RouteRedirects    `yaml:"redirects" json:"-"`
	SSI         RoutesSsi         `yaml:"ssi" json:"-"`
	Primary     Bool              `json:"primary"`
	OriginalURL string            `json:"original_url"`
	Disable     bool              `json:"_disable"`
}

// SetDefaults sets the default values.
func (d *Route) SetDefaults() {
	if d.Type == "" {
		d.Type = "upstream"
	}
	d.Cache.SetDefaults()
	d.Redirects.SetDefaults()
	d.SSI.SetDefaults()
	d.Primary.DefaultValue = false
	d.Primary.SetDefaults()
}

// Validate checks for errors.
func (d Route) Validate() []error {
	o := make([]error, 0)
	if e := d.Cache.Validate(&d); e != nil {
		o = append(o, e...)
	}
	if e := d.Redirects.Validate(&d); e != nil {
		o = append(o, e...)
	}
	if e := d.SSI.Validate(&d); e != nil {
		o = append(o, e...)
	}
	return o
}

// ParseRoutesYaml parses contents of routes.yaml file.
func ParseRoutesYaml(d []byte) ([]Route, error) {
	out := make([]Route, 0)
	routes := make(map[string]*Route)
	if err := yaml.Unmarshal(d, &routes); err != nil {
		return nil, err
	}
	for path, route := range routes {
		route.Path = strings.ReplaceAll(path, "{default}", defaultPath)
		route.SetDefaults()
		out = append(out, *route)
	}
	return out, nil
}

// ParseRoutesYamlFile opens the routes.yaml file and parses it.
func ParseRoutesYamlFile(f string) ([]Route, error) {
	d, err := ioutil.ReadFile(f)
	if err != nil {
		if os.IsNotExist(err) {
			return []Route{}, nil
		}
		return []Route{}, errors.WithStack(err)
	}
	r, err := ParseRoutesYaml(d)
	if err != nil {
		return r, errors.WithStack(err)
	}
	return r, nil
}

// ParseRoutesYamlFiles reads multiple routes yaml files and merges them.
func ParseRoutesYamlFiles(fileList []string) ([]Route, error) {
	done := output.Duration(
		fmt.Sprintf("Parse routes at '%s.'", strings.Join(fileList, ", ")),
	)
	out := make([]Route, 0)
	for _, f := range fileList {
		routes, err := ParseRoutesYamlFile(f)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		for _, route := range routes {
			hasOut := false
			for i := range out {
				if out[i].Path == route.Path {
					// attempt to merge routes using json marshaling and unmarshaling
					// because there is no other way to merge structs
					routeJSON, err := json.Marshal(route)
					if err != nil {
						return nil, errors.WithStack(err)
					}
					if err := json.Unmarshal(routeJSON, &out[i]); err != nil {
						return nil, errors.WithStack(err)
					}
					hasOut = true
					break
				}
			}
			if !hasOut {
				out = append(out, route)
			}
		}

	}
	done()
	return out, nil
}

// ExpandRoutes expands routes to include internal verisons and makes modifications for use with PCC.
func ExpandRoutes(inRoutes []Route, internalDomainSuffix string) ([]Route, error) {
	out := make([]Route, 0)
	// make copy of routes
	copyRoutes := make([]Route, 0)
	copyRoutes = append(copyRoutes, inRoutes...)
	redirectMatchesRoute := func(routes []Route, to string) bool {
		for _, route := range routes {
			url, err := url.Parse(route.Path)
			if err != nil {
				continue
			}
			toURL, err := url.Parse(to)
			if err != nil {
				continue
			}
			if url.Hostname() == toURL.Hostname() {
				return true
			}
		}
		return false
	}
	// internal routes will replace periods(.) with hyphens (-)
	// to create urls structures as follows...
	// example-com.platform.cc
	for _, route := range copyRoutes {
		// ignore if route doesn't start with http
		if !strings.HasPrefix(route.Path, "http") {
			continue
		}
		// original route
		route.Primary.DefaultValue = true
		route.SetDefaults()
		route.OriginalURL = route.Path
		r := route
		out = append(out, r)
		// internal route
		route.Primary.DefaultValue = false
		route.Primary.SetDefaults()
		var err error
		route.Path, err = convertRoutePathToInternal(route.Path, internalDomainSuffix)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		// - if redirect routes to internal then fix it up
		if route.To != "" &&
			redirectMatchesRoute(copyRoutes, route.To) {
			route.To, err = convertRoutePathToInternal(route.To, internalDomainSuffix)
			if err != nil {
				return nil, errors.WithStack(err)
			}
		}
		ri := route
		out = append(out, ri)
	}
	return out, nil
}

// RoutesToMap converts route array in to map.
func RoutesToMap(routes []Route) map[string]Route {
	data := make(map[string]Route)
	if routes == nil {
		return data
	}
	for _, route := range routes {
		data[route.Path] = route
	}
	return data
}
