package api

import (
	"log"
	"path"
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
	d, e := p.BuildConfigJSON()
	if e != nil {
		t.Errorf("failed to build config.json, %s", e)
	}
	t.Errorf(string(d))
}
