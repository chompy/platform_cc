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

package router

import (
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/contextualcode/platform_cc/v2/pkg/def"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/project"
)

// TestGenerateNginx tests the generation of nginx config.
func TestGenerateNginx(t *testing.T) {
	proj, err := project.LoadFromPath(
		filepath.Join("_test_data", "sample1"),
		true,
	)
	if err != nil {
		t.Error(err)
	}
	out, err := GenerateNginxConfig(proj)
	if err != nil {
		t.Error(err)
	}
	stringConf := string(out)
	def.AssertEqual(
		strings.Contains(stringConf, "return 301 /test"),
		true,
		"expected test.contextualcode.com /test/ redirect",
		t,
	)
	def.AssertEqual(
		strings.Contains(stringConf, "server_name contextualcode-com.platform.cc"),
		true,
		"expected contextualcode-com route",
		t,
	)
}
