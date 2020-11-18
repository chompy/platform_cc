package def

import (
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

const defaultPath = "default.pcc.localtest.me"

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
func ParseRoutesYaml(d []byte) ([]*Route, error) {
	out := make([]*Route, 0)
	routes := make(map[string]*Route)
	if err := yaml.Unmarshal(d, &routes); err != nil {
		return nil, err
	}
	for path, route := range routes {
		route.Path = strings.ReplaceAll(path, "{default}", defaultPath)
		route.SetDefaults()
		out = append(out, route)
	}
	return out, nil
}

// ParseRoutesYamlFile opens the routes.yaml file and parses it.
func ParseRoutesYamlFile(f string) ([]*Route, error) {
	log.Printf("Parse routes at '%s.'", f)
	d, e := ioutil.ReadFile(f)
	if e != nil {
		return []*Route{}, e
	}
	return ParseRoutesYaml(d)
}

// ExpandRoutes expands routes to include internal verisons and makes modifications for use with PCC.
func ExpandRoutes(inRoutes []*Route) ([]*Route, error) {
	out := make([]*Route, 0)
	// make copy of routes
	copyRoutes := make([]Route, 0)
	for _, route := range inRoutes {
		copyRoutes = append(copyRoutes, *route)
	}
	// convienence functions
	isRouteHTTPS := func(path string) bool {
		return strings.HasPrefix(path, "https://")
	}
	hasDuplicateHTTPSRoute := func(routes []Route, path string) bool {
		for _, route := range routes {
			if path != route.Path &&
				replaceHTTPS(route.Path) == replaceHTTPS(path) {
				return true
			}
		}
		return false
	}
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
	// PCC won't support https route, rename to http
	// internal routes will replace periods(.) with hyphens (-)
	// to create urls structures as followed...
	// example-com.pcc.local
	for _, route := range copyRoutes {
		// ignore if route doesn't start with http
		if !strings.HasPrefix(route.Path, "http") {
			continue
		}
		// if route is not https and the same route with https
		// exists then ignore this one
		if !isRouteHTTPS(route.Path) && hasDuplicateHTTPSRoute(copyRoutes, route.Path) {
			continue
		}
		// original route
		route.Primary.DefaultValue = true
		route.SetDefaults()
		route.Path = replaceHTTPS(route.Path)
		route.OriginalURL = route.Path
		// - if redirect routes to internal then fix it up
		if route.To != "" &&
			redirectMatchesRoute(copyRoutes, route.To) {
			route.To = replaceHTTPS(route.To)
		}
		r := route
		out = append(out, &r)
		// internal route
		route.Primary.DefaultValue = false
		route.Primary.SetDefaults()
		var err error
		route.Path, err = convertRoutePathToInternal(route.Path)
		if err != nil {
			return nil, err
		}
		// - if redirect routes to internal then fix it up
		if route.To != "" &&
			redirectMatchesRoute(copyRoutes, route.To) {
			route.To, err = convertRoutePathToInternal(replaceHTTPS(route.To))
			if err != nil {
				return nil, err
			}
		}
		ri := route
		out = append(out, &ri)
	}
	return out, nil
}
