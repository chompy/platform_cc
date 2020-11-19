package tests

import (
	"path"
	"strings"
	"testing"

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
	/*assertEqual(
		p.Apps[0].Variables["env"]["TEST_THREE"],
		"test123",
		"unexpected variables.env.TEST_THREE",
		t,
	)*/
	assertEqual(
		p.Apps[0].Type,
		"php:7.4",
		"unexpected app.type",
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
