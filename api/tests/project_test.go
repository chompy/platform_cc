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

package tests

import (
	"path"
	"strings"
	"testing"

	"gitlab.com/contextualcode/platform_cc/api/def"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

func TestFromPath(t *testing.T) {
	projectPath := path.Join("data", "sample2")
	p, e := project.LoadFromPath(projectPath, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	assertEqual(
		p.Apps[0].Name,
		"test_app2",
		"unspected app.type",
		t,
	)
}

func TestFromPathWithPCCAppYaml(t *testing.T) {
	projectPath := path.Join("data", "sample4")
	p, e := project.LoadFromPath(projectPath, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	assertEqual(
		p.Apps[0].Variables["env"]["TEST_ENV"],
		"no",
		"unexpected variables.env.TEST_ENV",
		t,
	)
	assertEqual(
		p.Apps[0].Variables["env"]["TEST_ENV_TWO"],
		"hello world",
		"unexpected variables.env.TEST_ENV_TWO",
		t,
	)
	assertEqual(
		p.Apps[0].Variables["env"]["TEST_THREE"],
		"test123",
		"unexpected variables.env.TEST_THREE",
		t,
	)
	assertEqual(
		p.Apps[0].Type,
		"php:7.4",
		"unexpected app.type",
		t,
	)
	assertEqual(
		len(p.Apps[0].Runtime.Extensions),
		3,
		"unexpected app.runtime.extensions length",
		t,
	)
	// also test service override
	// three services defined but mysqldb disabled in services.pcc.yaml override
	assertEqual(
		len(p.Services),
		2,
		"unexpected number of services",
		t,
	)
	assertEqual(
		p.Services[0].Type,
		"redis:3.2",
		"unspected service type",
		t,
	)
	assertEqual(
		p.Services[1].Type,
		"redis:3.2",
		"unspected service type",
		t,
	)
}

func TestConfigJSON(t *testing.T) {
	projectPath := path.Join("data", "sample2")
	p, e := project.LoadFromPath(projectPath, true)
	if e != nil {
		t.Errorf("failed to load project, %s", e)
	}
	d, e := p.BuildConfigJSON(p.Apps[0])
	if e != nil {
		t.Errorf("failed to build config.json, %s", e)
	}
	out := string(d)
	if !strings.Contains(out, "applications") {
		t.Error("config.json does not contain key applications")
	}
	if !strings.Contains(out, "PLATFORM_PROJECT_ENTROPY") {
		t.Error("config.json does not contain key PLATFORM_PROJECT_ENTROPY")
	}
}

func TestVariables(t *testing.T) {
	p := project.Project{
		Variables: make(map[string]map[string]string),
	}
	env := "dev"
	timeLimit := "30"
	apiSecret := "secret123"
	if err := p.VarSet("env:ENVIRONMENT", env); err != nil {
		t.Error(err)
	}
	if err := p.VarSet("php:set_time_limit", timeLimit); err != nil {
		t.Error(err)
	}
	if err := p.VarSet("secret:api_secret", apiSecret); err != nil {
		t.Error(err)
	}
	v, err := p.VarGet("env:ENVIRONMENT")
	if err != nil {
		t.Error(err)
	}
	assertEqual(v, env, "env:ENVIRONMENT unexpected value", t)
	v, err = p.VarGet("php:set_time_limit")
	if err != nil {
		t.Error(err)
	}
	assertEqual(v, timeLimit, "php:set_time_limit unexpected value", t)
	v, err = p.VarGet("secret:api_secret")
	if err != nil {
		t.Error(err)
	}
	assertEqual(v, apiSecret, "secret:api_secret unexpected value", t)
	if err := p.VarDelete("secret:api_secret"); err != nil {
		t.Error(err)
	}
	v, err = p.VarGet("secret:api_secret")
	if err != nil {
		t.Error(err)
	}
	assertEqual(v, "", "secret:api_secret unexpected value", t)
}

func TestYAMLFunction(t *testing.T) {
	projectPath := path.Join("data", "sample5")
	p, e := project.LoadFromPath(projectPath, true)
	if e != nil {
		tracerr.PrintSourceColor(e)
		t.Errorf("failed to load project, %s", e)
	}
	assertEqual(
		p.Apps[0].Name,
		"test_app",
		"unspected app.type",
		t,
	)
	assertEqual(
		strings.Contains(p.Apps[0].Hooks.Build, "TEST!"),
		true,
		"incorrect build hook",
		t,
	)
}

func TestFlags(t *testing.T) {
	gc := def.GlobalConfig{
		Flags: []string{
			project.EnableOSXNFSMounts, project.EnableServiceRoutes, project.DisableYamlOverrides, project.DisableAutoCommit,
		},
	}
	p := project.Project{
		Flags: project.Flags{
			project.EnablePHPOpcache:     project.FlagOn,
			project.EnableOSXNFSMounts:   project.FlagOn,
			project.EnableCron:           project.FlagOff,
			project.EnableServiceRoutes:  project.FlagOff,
			project.EnableWorkers:        project.FlagUnset,
			project.DisableYamlOverrides: project.FlagUnset,
		},
	}
	p.SetGlobalConfig(&gc)
	assertEqual(
		p.HasFlag(project.EnablePHPOpcache),
		true,
		"expected flag 'enable_php_opcache' on",
		t,
	)
	assertEqual(
		p.HasFlag(project.EnableOSXNFSMounts),
		true,
		"expected flag 'enable_osx_nfs_mounts' on",
		t,
	)
	assertEqual(
		p.HasFlag(project.EnableCron),
		false,
		"expected flag 'enable_cron' off",
		t,
	)
	assertEqual(
		p.HasFlag(project.EnableServiceRoutes),
		false,
		"expected flag 'enable_service_routes' off",
		t,
	)
	assertEqual(
		p.HasFlag(project.EnableWorkers),
		false,
		"expected flag 'enable_workers' off",
		t,
	)
	assertEqual(
		p.HasFlag(project.DisableYamlOverrides),
		true,
		"expected flag 'disable_yaml_overrides' on",
		t,
	)
	assertEqual(
		p.HasFlag(project.DisableAutoCommit),
		true,
		"expected flag 'disable_auto_commit' on",
		t,
	)
}
