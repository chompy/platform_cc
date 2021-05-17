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
	"path"
	"testing"
)

const internalDomainSuffix = "pcc.test"

func TestRoutes(t *testing.T) {
	p := path.Join("_test_data", "sample1", ".platform", "routes.yaml")
	routes, err := ParseRoutesYamlFile(p)
	if err != nil {
		t.Errorf("failed to parse routes yaml, %s", err)
	}
	routes, err = ExpandRoutes(routes, internalDomainSuffix)
	if err != nil {
		t.Errorf("failed to expand routes, %s", err)
	}
	containsRoutes := []string{
		"https://www.contextualcode.com/",
		"https://www-contextualcode-com." + internalDomainSuffix + "/",
		"https://cdn-backend.contextualcode.ccplatform.net/",
		"https://cdn-backend-contextualcode-ccplatform-net." + internalDomainSuffix + "/",
		"http://health.contextualcode.ccplatform.net/",
		"http://health-contextualcode-ccplatform-net." + internalDomainSuffix + "/",
		"https://contextualcode.com/",
		"https://contextualcode-com." + internalDomainSuffix + "/",
	}
	routeMap := RoutesToMap(routes)
	for _, path := range containsRoutes {
		AssertEqual(
			routeMap[path].Type != "",
			true,
			fmt.Sprintf("missing expected route '%s'", path),
			t,
		)
	}
	AssertEqual(
		routeMap["https://contextualcode.com/"].To,
		"https://www.contextualcode.com/",
		"unexpected redirect",
		t,
	)
	AssertEqual(
		routeMap["https://contextualcode-com."+internalDomainSuffix+"/"].To,
		"https://www-contextualcode-com."+internalDomainSuffix+"/",
		"unexpected redirect",
		t,
	)
	AssertEqual(
		routeMap["https://www.contextualcode.com/"].Upstream,
		"test_app:http",
		"unexpected upstream",
		t,
	)
}

func TestRoutesRedirect(t *testing.T) {
	d, e := ParseRoutesYaml([]byte(`
http://example.com:
    redirects:
        expires: 1d
        paths:
            "/from":
                to: "http://example.com/to"
            "^/foo/(.*)/bar":
                to: "http://example.com/$1"
                regexp: true
`))
	if e != nil {
		t.Errorf("failed to parse app yaml, %s", e)
	}
	AssertEqual(
		d[0].Redirects.Paths["/from"].Prefix.Get(),
		true,
		"unexpected routes[].redirects.paths[].prefix",
		t,
	)
	AssertEqual(
		d[0].Redirects.Paths["/from"].AppendSuffix.Get(),
		true,
		"unexpected routes[].redirects.paths[].append_suffix",
		t,
	)
	AssertEqual(
		d[0].Redirects.Paths["/from"].Regexp.Get(),
		false,
		"unexpected routes[].redirects.paths[].regexp",
		t,
	)
	AssertEqual(
		d[0].Redirects.Paths["^/foo/(.*)/bar"].Regexp.Get(),
		true,
		"unexpected routes[].redirects.paths[].regexp",
		t,
	)
	AssertEqual(
		d[0].Redirects.Paths["^/foo/(.*)/bar"].To,
		"http://example.com/$1",
		"unexpected routes[].redirects.paths[].to",
		t,
	)
}

func TestLargeRoutes(t *testing.T) {
	p := path.Join("_test_data", "sample3", "routes.yaml")
	routes, e := ParseRoutesYamlFile(p)
	if e != nil {
		t.Errorf("failed to parse routes yaml, %s", e)
	}
	//log.Println(routes[0].Validate())
	AssertEqual(
		len(routes[0].Validate()),
		0,
		"expected route to validate",
		t,
	)
}

func TestRouteOverride(t *testing.T) {
	paths := []string{
		path.Join("_test_data", "route_override", ".platform", "routes.yaml"),
		path.Join("_test_data", "route_override", ".platform", "routes.pcc.yaml"),
	}
	routes, err := ParseRoutesYamlFiles(paths)
	if err != nil {
		t.Errorf("failed to parse routes yaml, %s", err)
	}
	routes, err = ExpandRoutes(routes, internalDomainSuffix)
	if err != nil {
		t.Errorf("failed to expand routes, %s", err)
	}
	routeMap := RoutesToMap(routes)
	AssertEqual(
		routeMap["https://www.contextualcode.com/"].Upstream,
		"test_app:http",
		"unexpected upstream for www.contextualcode.com",
		t,
	)
	AssertEqual(
		routeMap["https://dev.contextualcode.com/"].Upstream,
		"test_app:http",
		"unexpected upstream for dev.contextualcode.com",
		t,
	)
	AssertEqual(
		routeMap["https://contextualcode.com/"].To,
		"https://dev.contextualcode.com/",
		"unexpected redirect",
		t,
	)
}
