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

package project

import (
	"path"
	"strings"
	"testing"

	"gitlab.com/contextualcode/platform_cc/v2/pkg/def"
)

func TestFromPath(t *testing.T) {
	projectPath := path.Join("_test_data", "sample2")
	p, e := LoadFromPath(projectPath, true)
	if e != nil {
		t.Errorf("failed to load project, %s", e)
	}
	def.AssertEqual(
		p.Apps[0].Name,
		"test_app2",
		"unspected app.type",
		t,
	)
}

func TestFromPathWithPCCAppYaml(t *testing.T) {
	projectPath := path.Join("_test_data", "sample4")
	p, e := LoadFromPath(projectPath, true)
	if e != nil {
		t.Errorf("failed to load project, %s", e)
	}
	def.AssertEqual(
		p.Apps[0].Variables.GetString("env:TEST_ENV"),
		"no",
		"unexpected variables.env.TEST_ENV",
		t,
	)
	def.AssertEqual(
		p.Apps[0].Variables.GetString("env:TEST_ENV_TWO"),
		"hello world",
		"unexpected variables.env.TEST_ENV_TWO",
		t,
	)
	def.AssertEqual(
		p.Apps[0].Variables.GetString("env:TEST_THREE"),
		"test123",
		"unexpected variables.env.TEST_THREE",
		t,
	)
	subMap := p.Apps[0].Variables.GetStringSubMap("env")
	def.AssertEqual(
		subMap["TEST_ENV"],
		"no",
		"unexpected variables.env.TEST_ENV (sub map)",
		t,
	)
	p.Variables.Set("php:memory_limit", "1024M")
	def.AssertEqual(
		p.Variables.Get("php:memory_limit"),
		"1024M",
		"unexpected variables.env.TEST_ENV (sub map)",
		t,
	)
	p.Variables.Delete("php:memory_limit")
	def.AssertEqual(
		p.Variables.Get("php:memory_limit"),
		nil,
		"unexpected variables.env.TEST_ENV (sub map)",
		t,
	)
	def.AssertEqual(
		p.Apps[0].Type,
		"php:7.4",
		"unexpected app.type",
		t,
	)
	def.AssertEqual(
		len(p.Apps[0].Runtime.Extensions),
		3,
		"unexpected app.runtime.extensions length",
		t,
	)
	// also test service override
	// three services defined but mysqldb disabled in services.pcc.yaml override
	def.AssertEqual(
		len(p.Services),
		2,
		"unexpected number of services",
		t,
	)
	def.AssertEqual(
		p.Services[0].Type,
		"redis:3.2",
		"unspected service type",
		t,
	)
	def.AssertEqual(
		p.Services[1].Type,
		"redis:3.2",
		"unspected service type",
		t,
	)
}

func TestConfigJSON(t *testing.T) {
	projectPath := path.Join("_test_data", "sample2")
	p, e := LoadFromPath(projectPath, true)
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
	p := Project{
		Variables: make(def.Variables),
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
	def.AssertEqual(v, env, "env:ENVIRONMENT unexpected value", t)
	v, err = p.VarGet("php:set_time_limit")
	if err != nil {
		t.Error(err)
	}
	def.AssertEqual(v, timeLimit, "php:set_time_limit unexpected value", t)
	v, err = p.VarGet("secret:api_secret")
	if err != nil {
		t.Error(err)
	}
	def.AssertEqual(v, apiSecret, "secret:api_secret unexpected value", t)
	if err := p.VarDelete("secret:api_secret"); err != nil {
		t.Error(err)
	}
	v, err = p.VarGet("secret:api_secret")
	if err != nil {
		t.Error(err)
	}
	def.AssertEqual(v, "", "secret:api_secret unexpected value", t)
}

func TestYAMLFunction(t *testing.T) {
	projectPath := path.Join("_test_data", "sample5")
	p, e := LoadFromPath(projectPath, true)
	if e != nil {
		t.Errorf("failed to load project, %s", e)
	}
	def.AssertEqual(
		p.Apps[0].Name,
		"test_app",
		"unspected app.type",
		t,
	)
	def.AssertEqual(
		strings.Contains(p.Apps[0].Hooks.Build, "TEST!"),
		true,
		"incorrect build hook",
		t,
	)
}

func TestFlags(t *testing.T) {
	p := Project{
		Flags: Flags{
			EnablePHPOpcache:     FlagOn,
			EnableOSXNFSMounts:   FlagOn,
			EnableCron:           FlagOff,
			EnableServiceRoutes:  FlagOff,
			EnableWorkers:        FlagUnset,
			DisableYamlOverrides: FlagUnset,
		},
	}
	p.SetGlobalConfig(
		def.GlobalConfig{
			Flags: []string{
				EnableOSXNFSMounts, EnableServiceRoutes, DisableYamlOverrides, DisableAutoCommit,
			},
		},
	)
	def.AssertEqual(
		p.HasFlag(EnablePHPOpcache),
		true,
		"expected flag 'enable_php_opcache' on",
		t,
	)
	def.AssertEqual(
		p.HasFlag(EnableOSXNFSMounts),
		true,
		"expected flag 'enable_osx_nfs_mounts' on",
		t,
	)
	def.AssertEqual(
		p.HasFlag(EnableCron),
		false,
		"expected flag 'enable_cron' off",
		t,
	)
	def.AssertEqual(
		p.HasFlag(EnableServiceRoutes),
		false,
		"expected flag 'enable_service_routes' off",
		t,
	)
	def.AssertEqual(
		p.HasFlag(EnableWorkers),
		false,
		"expected flag 'enable_workers' off",
		t,
	)
	def.AssertEqual(
		p.HasFlag(DisableYamlOverrides),
		true,
		"expected flag 'disable_yaml_overrides' on",
		t,
	)
	def.AssertEqual(
		p.HasFlag(DisableAutoCommit),
		true,
		"expected flag 'disable_auto_commit' on",
		t,
	)
}
