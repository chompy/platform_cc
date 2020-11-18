package router

import (
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/contextualcode/platform_cc/api/project"
)

// TestGenerateNginx tests the generation of nginx config.
func TestGenerateNginx(t *testing.T) {

	proj, err := project.LoadFromPath(
		filepath.Join("..", "def", "test_data", "sample1"),
		true,
	)
	if err != nil {
		t.Error(err)
	}

	out, err := GenerateNginxConfig(proj)
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
		strings.Contains(stringConf, "server_name contextualcode-com.pcc.test"),
		true,
		"expected contextualcode-com route",
		t,
	)
}
