package tests

import (
	"fmt"
	"log"
	"path"
	"testing"

	"gitlab.com/contextualcode/platform_cc/api/def"
)

const internalDomainSuffix = "pcc.test"

func TestRoutes(t *testing.T) {
	p := path.Join("data", "sample1", ".platform", "routes.yaml")
	routes, err := def.ParseRoutesYamlFile(p)
	if err != nil {
		t.Errorf("failed to parse routes yaml, %s", err)
	}
	routes, err = def.ExpandRoutes(routes, internalDomainSuffix)
	if err != nil {
		t.Errorf("failed to expand routes, %s", err)
	}
	containsRoutes := []string{
		"http://www.contextualcode.com/",
		"http://www-contextualcode-com." + internalDomainSuffix + "/",
		"http://cdn-backend.contextualcode.ccplatform.net/",
		"http://cdn-backend-contextualcode-ccplatform-net." + internalDomainSuffix + "/",
		"http://health.contextualcode.ccplatform.net/",
		"http://health-contextualcode-ccplatform-net." + internalDomainSuffix + "/",
		"http://contextualcode.com/",
		"http://contextualcode-com." + internalDomainSuffix + "/",
	}
	for _, path := range containsRoutes {
		hasRoute := false
		for _, route := range routes {
			if route.Path == path {
				hasRoute = true
				if path == "http://contextualcode.com/" {
					assertEqual(
						route.To,
						"http://www.contextualcode.com/",
						"unexpected redirect",
						t,
					)
				} else if path == "http://contextualcode-com."+internalDomainSuffix+"/" {
					assertEqual(
						route.To,
						"http://www-contextualcode-com."+internalDomainSuffix+"/",
						"unexpected redirect",
						t,
					)
				}
				break
			}
		}
		assertEqual(
			hasRoute,
			true,
			fmt.Sprintf("missing expected route '%s'", path),
			t,
		)
	}
}

func TestRoutesRedirect(t *testing.T) {
	d, e := def.ParseRoutesYaml([]byte(`
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
	assertEqual(
		d[0].Redirects.Paths["/from"].Prefix.Get(),
		true,
		"unexpected routes[].redirects.paths[].prefix",
		t,
	)
	assertEqual(
		d[0].Redirects.Paths["/from"].AppendSuffix.Get(),
		true,
		"unexpected routes[].redirects.paths[].append_suffix",
		t,
	)
	assertEqual(
		d[0].Redirects.Paths["/from"].Regexp.Get(),
		false,
		"unexpected routes[].redirects.paths[].regexp",
		t,
	)
	assertEqual(
		d[0].Redirects.Paths["^/foo/(.*)/bar"].Regexp.Get(),
		true,
		"unexpected routes[].redirects.paths[].regexp",
		t,
	)
	assertEqual(
		d[0].Redirects.Paths["^/foo/(.*)/bar"].To,
		"http://example.com/$1",
		"unexpected routes[].redirects.paths[].to",
		t,
	)
}

func TestLargeRoutes(t *testing.T) {
	p := path.Join("data", "sample3", "routes.yaml")
	routes, e := def.ParseRoutesYamlFile(p)
	if e != nil {
		t.Errorf("failed to parse routes yaml, %s", e)
	}
	log.Println(routes[0].Validate())
	assertEqual(
		len(routes[0].Validate()),
		0,
		"expected route to validate",
		t,
	)
}
