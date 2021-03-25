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
	"testing"
)

func TestParseFile(t *testing.T) {
	p := []string{path.Join("_test_data", "sample1", ".platform.app.yaml")}
	d, e := ParseAppYamlFiles(p, nil)
	if e != nil {
		t.Errorf("failed to parse app yaml, %s", e)
	}
	AssertEqual(d.Name, "test_app", "unexpected name", t)
	AssertEqual(d.Type, "php:7.4", "unexpected type", t)
	AssertEqual(d.Build.Flavor, "none", "unexpected build.flavor", t)
	// locations
	AssertEqual(d.Web.Locations["/"].Passthru.GetString(), "/index.php", "unexpected web.locations./.passthru", t)
	// mounts
	AssertEqual(d.Mounts["/test"].SourcePath, "test", "unexpected web.mounts./test.source_path", t)
	AssertEqual(d.Mounts["/test2"].Source, "service", "unexpected web.mounts./test2.source", t)
	AssertEqual(d.Mounts["/test2"].SourcePath, "test2", "unexpected web.mounts./test2.source_path", t)
	AssertEqual(d.Mounts["/test2"].Service, "files", "unexpected web.mounts./test2.service", t)
	// php extensions
	AssertEqual(d.Runtime.Extensions[0].Name, "imagick", "unexpected runtime.extensions[].name", t)
	AssertEqual(d.Runtime.Extensions[1].Name, "xsl", "unexpected runtime.extensions[].name", t)
	AssertEqual(d.Runtime.Extensions[2].Name, "blackfire", "unexpected runtime.extensions[].name", t)
	AssertEqual(d.Runtime.Extensions[2].Configuration["server_id"], "test123", "unexpected runtime.extensions[].configuration.server_id", t)
}

func TestInvalidCron(t *testing.T) {
	_, e := ParseAppYamls([][]byte{[]byte(`
name: test_app_cron
type: php:7.4
crons:
    test:
        spec: "*/5 * * *"
        cmd: "sleep 5"
`)}, nil)
	if e != nil {
		t.Errorf("failed to parse app yaml, %s", e)
	}
	/*if e := d.Validate(); len(e) == 0 {
		// TODO fix cron parser
		//t.Error("expected cron parse error")
	}*/
}

func TestInvalidMount(t *testing.T) {
	d, e := ParseAppYamls([][]byte{[]byte(`
name: test_app_cron
type: php:7.4
mounts:
    /test:
        source_path: test2
        source: this_does_not_exist
        service: files
`)}, nil)
	if e != nil {
		t.Errorf("failed to parse app yaml, %s", e)
	}
	ve := d.Validate()
	if len(ve) != 1 {
		t.Error("expected mount parse error")
	}
	switch ve[0].(type) {
	case *ValidateError:
		{
			break
		}
	default:
		{
			t.Errorf("expected def.ValidateError")
			break
		}
	}
}

func TestPHPDependencies(t *testing.T) {
	d, e := ParseAppYamls([][]byte{[]byte(`
name: test_app_cron
type: php:7.4
dependencies:
    php:
        "platformsh/client": "dev-master"
        "something/something": "~1.4"
`)}, nil)
	if e != nil {
		t.Errorf("failed to parse app yaml, %s", e)
	}
	AssertEqual(
		d.Dependencies.PHP.Require["platformsh/client"],
		"dev-master",
		"unexpected dependencies.php.require[]",
		t,
	)
	AssertEqual(
		d.Dependencies.PHP.Require["something/something"],
		"~1.4",
		"unexpected dependencies.php.require[]",
		t,
	)
	AssertEqual(len(d.Dependencies.PHP.Repositories), 0, "unexpected length for dependencies.php.repositories", t)
}

func TestPHPDependenciesExpanded(t *testing.T) {
	d, e := ParseAppYamls([][]byte{[]byte(`
name: test_app_cron
type: php:7.4
dependencies:
    php:
        require:
            "platformsh/client": "dev-master"
            "something/something": "~1.4"
        repositories:
            - type: vcs
              url: "git@github.com:platformsh/platformsh-client-php.git"
`)}, nil)
	if e != nil {
		t.Errorf("failed to parse app yaml, %s", e)
	}
	AssertEqual(
		d.Dependencies.PHP.Require["platformsh/client"],
		"dev-master",
		"unexpected dependencies.php.require[]",
		t,
	)
	AssertEqual(
		d.Dependencies.PHP.Require["something/something"],
		"~1.4",
		"unexpected dependencies.php.require[]",
		t,
	)
	AssertEqual(
		d.Dependencies.PHP.Repositories[0].Type,
		"vcs",
		"unexpected dependencies.php.repositories[].type",
		t,
	)
	AssertEqual(
		d.Dependencies.PHP.Repositories[0].URL,
		"git@github.com:platformsh/platformsh-client-php.git",
		"unexpected dependencies.php.repositories[].url",
		t,
	)
}
