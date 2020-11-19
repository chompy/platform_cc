package tests

import (
	"log"
	"path"
	"testing"

	"gitlab.com/contextualcode/platform_cc/api/def"
)

func TestParseFile(t *testing.T) {
	p := path.Join("data", "sample1", ".platform.app.yaml")
	d, e := def.ParseAppYamlFile(p, nil)
	if e != nil {
		t.Errorf("failed to parse app yaml, %s", e)
	}
	assertEqual(d.Name, "test_app", "unexpected name", t)
	assertEqual(d.Type, "php:7.4", "unexpected type", t)
	assertEqual(d.Build.Flavor, "none", "unexpected build.flavor", t)
	// locations
	assertEqual(d.Web.Locations["/"].Passthru, "/index.php", "unexpected web.locations./.passthru", t)
	// mounts
	assertEqual(d.Mounts["/test"].SourcePath, "test", "unexpected web.mounts./test.source_path", t)
	assertEqual(d.Mounts["/test2"].Source, "service", "unexpected web.mounts./test2.source", t)
	assertEqual(d.Mounts["/test2"].SourcePath, "test2", "unexpected web.mounts./test2.source_path", t)
	assertEqual(d.Mounts["/test2"].Service, "files", "unexpected web.mounts./test2.service", t)
	// php extensions
	assertEqual(d.Runtime.Extensions[0].Name, "imagick", "unexpected runtime.extensions[].name", t)
	assertEqual(d.Runtime.Extensions[1].Name, "xsl", "unexpected runtime.extensions[].name", t)
	assertEqual(d.Runtime.Extensions[2].Name, "blackfire", "unexpected runtime.extensions[].name", t)
	assertEqual(d.Runtime.Extensions[2].Configuration["server_id"], "test123", "unexpected runtime.extensions[].configuration.server_id", t)
}

func TestInvalidCron(t *testing.T) {
	d, e := def.ParseAppYaml([]byte(`
name: test_app_cron
type: php:7.4
crons:
    test:
        spec: "*/5 * * *"
        cmd: "sleep 5"
`), nil)
	if e != nil {
		t.Errorf("failed to parse app yaml, %s", e)
	}
	if e := d.Validate(); len(e) == 0 {
		t.Error("expected cron parse error")
	}
}

func TestInvalidMount(t *testing.T) {
	d, e := def.ParseAppYaml([]byte(`
name: test_app_cron
type: php:7.4
mounts:
    /test:
        source_path: test2
        source: this_does_not_exist
        service: files
`), nil)
	if e != nil {
		t.Errorf("failed to parse app yaml, %s", e)
	}
	ve := d.Validate()
	if len(ve) != 1 {
		t.Error("expected mount parse error")
	}
	switch ve[0].(type) {
	case *def.ValidateError:
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

	d, e := def.ParseAppYaml([]byte(`
name: test_app_cron
type: php:7.4
dependencies:
    php:
        "platformsh/client": "dev-master"
        "something/something": "~1.4"
`), nil)
	if e != nil {
		t.Errorf("failed to parse app yaml, %s", e)
	}
	log.Println(d.Dependencies.PHP)
	assertEqual(
		d.Dependencies.PHP.Require["platformsh/client"],
		"dev-master",
		"unexpected dependencies.php.require[]",
		t,
	)
	assertEqual(
		d.Dependencies.PHP.Require["something/something"],
		"~1.4",
		"unexpected dependencies.php.require[]",
		t,
	)
	assertEqual(len(d.Dependencies.PHP.Repositories), 0, "unexpected length for dependencies.php.repositories", t)
}

func TestPHPDependenciesExpanded(t *testing.T) {
	d, e := def.ParseAppYaml([]byte(`
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
`), nil)
	if e != nil {
		t.Errorf("failed to parse app yaml, %s", e)
	}
	log.Println(d.Dependencies.PHP)
	assertEqual(
		d.Dependencies.PHP.Require["platformsh/client"],
		"dev-master",
		"unexpected dependencies.php.require[]",
		t,
	)
	assertEqual(
		d.Dependencies.PHP.Require["something/something"],
		"~1.4",
		"unexpected dependencies.php.require[]",
		t,
	)
	assertEqual(
		d.Dependencies.PHP.Repositories[0].Type,
		"vcs",
		"unexpected dependencies.php.repositories[].type",
		t,
	)
	assertEqual(
		d.Dependencies.PHP.Repositories[0].URL,
		"git@github.com:platformsh/platformsh-client-php.git",
		"unexpected dependencies.php.repositories[].url",
		t,
	)
}
