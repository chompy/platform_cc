package api

import (
	"path"
	"testing"
)

func TestRoutes(t *testing.T) {
	p := path.Join("test_data", "sample1", "routes.yaml")
	routes, e := ParseRoutesYamlFile(p)
	if e != nil {
		t.Errorf("failed to parse routes yaml, %s", e)
	}
	assertEqual(
		routes[1].Path,
		"https://www.contextualcode.com/",
		"unexpected routes[].path",
		t,
	)
	assertEqual(
		routes[1].Upstream,
		"app:http",
		"unexpected routes[].upstream",
		t,
	)
	assertEqual(
		routes[1].Cache.Enabled.Get(),
		false,
		"unexpected routes[].cache.enabled",
		t,
	)
	assertEqual(
		routes[1].SSI.Enabled.Get(),
		true,
		"unexpected routes[].ssi.enabled",
		t,
	)
	assertEqual(
		routes[3].Path,
		"http://health.contextualcode.ccplatform.net/",
		"unexpected routes[].path",
		t,
	)
	assertEqual(
		routes[3].Upstream,
		"app:http",
		"unexpected routes[].upstream",
		t,
	)
	assertEqual(
		routes[4].Type,
		"redirect",
		"unexpected routes[].type",
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
	p := path.Join("test_data", "sample3", "routes.yaml")
	routes, e := ParseRoutesYamlFile(p)
	if e != nil {
		t.Errorf("failed to parse routes yaml, %s", e)
	}
	assertEqual(
		len(routes[0].Validate()),
		0,
		"expected route to validate",
		t,
	)
}
