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
	"path"
	"strings"
	"testing"
)

func TestServices(t *testing.T) {
	p := path.Join("_test_data", "sample1", ".platform", "services.yaml")
	services, e := ParseServiceYamlFiles([]string{p})
	if e != nil {
		t.Errorf("failed to parse services yaml, %s", e)
	}
	AssertEqual(
		len(services),
		4,
		"unexpected number of services",
		t,
	)
	for _, service := range services {
		switch service.Name {
		case "mysqldb":
			{
				AssertEqual(
					service.Type,
					"mysql:10.0",
					"unexpected services[].type",
					t,
				)
				AssertEqual(
					service.Disk,
					512,
					"unexpected services[].disk",
					t,
				)
				break
			}
		case "redis-cache":
			{
				AssertEqual(
					service.Type,
					"redis:3.2",
					"unexpected services[].type",
					t,
				)
				AssertEqual(
					service.Disk,
					0,
					"unexpected services[].disk",
					t,
				)
				break
			}
		case "redis-session":
			{
				AssertEqual(
					service.Type,
					"redis-persistent:3.2",
					"unexpected services[].type",
					t,
				)
				AssertEqual(
					service.Disk,
					1024,
					"unexpected services[].disk",
					t,
				)
				break
			}
		case "solr-search":
			{
				AssertEqual(
					service.Type,
					"solr:6.6",
					"unexpected services[].type",
					t,
				)
				conf := service.Configuration["cores"].(map[string]interface{})["test"].(map[string]interface{})["conf_dir"].(string)
				AssertEqual(
					strings.HasPrefix(conf, "H4sIAAAAAAAA"),
					true,
					"unexpected services[].configuration.cores[].conf_dir",
					t,
				)
				break
			}
		}

	}
}
