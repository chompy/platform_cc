package api

import (
	"log"
	"path"
	"strings"
	"testing"
)

func TestProjectFromPath(t *testing.T) {
	projectPath := path.Join("test_data", "sample2")
	p, e := LoadProjectFromPath(projectPath, true)
	if e != nil {
		t.Errorf("failed to load project, %s", e)
	}
	log.Println(p.ID)
}

func TestProjectConfigJSON(t *testing.T) {
	projectPath := path.Join("test_data", "sample2")
	p, e := LoadProjectFromPath(projectPath, true)
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

func TestProjectVariables(t *testing.T) {
	p := Project{
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
