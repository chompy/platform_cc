package router

import (
	"log"
	"net/url"

	"gitlab.com/contextualcode/platform_cc/api/def"
)

// HostRoute contains all routes mapped to a host.
type HostRoute struct {
	Host   string
	Routes []*def.Route
}

// AddRoute adds a given route to the host map if its host matches.
func (h *HostRoute) AddRoute(r *def.Route) bool {
	parsedURL, err := url.Parse(r.Path)
	if err != nil {
		log.Printf("URL parse error: %s", err)
		return true
	}
	if parsedURL.Hostname() == h.Host {
		h.Routes = append(h.Routes, r)
		return true
	}
	return false
}

// MapHostRoutes maps all routes to their hosts.
func MapHostRoutes(routes []*def.Route) []HostRoute {
	out := make([]HostRoute, 0)
	for _, route := range routes {
		hasHostRoute := false
		for i := range out {
			hasHostRoute = out[i].AddRoute(route)
			if hasHostRoute {
				break
			}
		}
		if !hasHostRoute {
			parsedURL, err := url.Parse(route.Path)
			if err != nil {
				log.Printf("URL parse error: %s", err)
				continue
			}
			hostRoute := HostRoute{
				Host:   parsedURL.Hostname(),
				Routes: make([]*def.Route, 0),
			}
			hostRoute.AddRoute(route)
			out = append(out, hostRoute)
		}
	}
	return out
}
