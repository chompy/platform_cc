package tests

import (
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/contextualcode/platform_cc/api/project"
	"gitlab.com/contextualcode/platform_cc/api/router"
)

// TestGenerateNginx tests the generation of nginx config.
func TestGenerateNginx(t *testing.T) {

	proj, err := project.LoadFromPath(
		filepath.Join("data", "sample1"),
		true,
	)
	if err != nil {
		t.Error(err)
	}

	out, err := router.GenerateNginxConfig(proj)
	if err != nil {
		t.Error(err)
	}
	//t.Error(string(out))
	stringConf := string(out)
	assertEqual(
		strings.Contains(stringConf, "return 301 /test"),
		true,
		"expected test.contextualcode.com /test/ redirect",
		t,
	)
	assertEqual(
		strings.Contains(stringConf, "server_name contextualcode-com.pcc."),
		true,
		"expected contextualcode-com route",
		t,
	)
}
