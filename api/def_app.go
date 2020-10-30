package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

const appDir = "/app"

// AppDef - defines an application
type AppDef struct {
	Path          string
	Name          string                            `yaml:"name"`
	Type          string                            `yaml:"type"`
	Size          string                            `yaml:"size"`
	Disk          int                               `yaml:"disk"`
	Build         AppBuildDef                       `yaml:"build"`
	Variables     map[string]map[string]interface{} `yaml:"variables"`
	Relationships map[string]string                 `yaml:"relationships"`
	Web           AppWebDef                         `yaml:"web"`
	Mounts        map[string]*AppMountDef           `yaml:"mounts" json:"mounts"`
	Hooks         AppHooksDef                       `yaml:"hooks"`
	Crons         map[string]*AppCronDef            `yaml:"crons" json:"crons"`
	Dependencies  AppDependenciesDef                `yaml:"dependencies"`
	Runtime       AppRuntimeDef                     `yaml:"runtime"`
}

// SetDefaults - set default values
func (d *AppDef) SetDefaults() {
	if d.Name == "" {
		d.Name = "app"
	}
	if d.Type == "" {
		d.Type = "php:7.4"
	}
	if d.Size == "" {
		d.Size = "M"
	}
	if d.Disk < 256 {
		d.Disk = 256
	}
	d.Build.SetDefaults()
	d.Web.SetDefaults()
	for i := range d.Mounts {
		d.Mounts[i].SetDefaults()
	}
	d.Hooks = AppHooksDef{}
	d.Hooks.SetDefaults()
	for i := range d.Crons {
		d.Crons[i].SetDefaults()
	}
	d.Dependencies.SetDefaults()
	d.Runtime.SetDefaults()
}

// Validate - validate AppDef
func (d AppDef) Validate() []error {
	o := make([]error, 0)
	if e := d.Build.Validate(); len(e) > 0 {
		o = append(o, e...)
	}
	if e := d.Web.Validate(); len(e) > 0 {
		o = append(o, e...)
	}
	for _, m := range d.Mounts {
		if e := m.Validate(); len(e) > 0 {
			o = append(o, e...)
		}
	}
	if e := d.Hooks.Validate(); len(e) > 0 {
		o = append(o, e...)
	}
	for _, c := range d.Crons {
		if e := c.Validate(); len(e) > 0 {
			o = append(o, e...)
		}
	}
	if e := d.Dependencies.Validate(); len(e) > 0 {
		o = append(o, e...)
	}
	if e := d.Runtime.Validate(); len(e) > 0 {
		o = append(o, e...)
	}
	return o
}

// GetContainerImage - get container image for app
func (d AppDef) GetContainerImage() string {
	typeName := strings.Split(d.Type, ":")
	return fmt.Sprintf("%s%s-%s", platformShDockerImagePrefix, typeName[0], typeName[1])
}

// MarshalJSON - implement json marshaler interface
func (d *AppDef) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"crons":                 d.Crons,
		"enable_smtp":           false,
		"mounts":                d.Mounts,
		"cron_minimum_interval": "1",
		"configuration": map[string]interface{}{
			"app_dir":       appDir,
			"hooks":         d.Hooks,
			"variables":     map[string]string{},
			"timezone":      nil,
			"disk":          d.Disk,
			"slug_id":       "--",
			"size":          "AUTO",
			"relationships": d.Relationships,
			"web":           d.Web,
			"is_production": false,
			"name":          d.Name,
			"access":        map[string]string{},
			"preflight": map[string]interface{}{
				"enabled":       true,
				"ignored_rules": []string{},
			},
			"tree_id":   "-",
			"mounts":    d.Mounts,
			"runtime":   d.Runtime,
			"type":      d.Type,
			"crons":     d.Crons,
			"slug":      "-",
			"resources": nil,
		},
	})
}

// ParseAppYaml - parse app yaml
func ParseAppYaml(d []byte) (*AppDef, error) {
	o := &AppDef{}
	e := yaml.Unmarshal(d, &o)
	o.SetDefaults()
	return o, e
}

// ParseAppYamlFile - open app yaml file and parse it
func ParseAppYamlFile(f string) (*AppDef, error) {
	log.Printf("Parse app at '%s.'", f)
	d, e := ioutil.ReadFile(f)
	if e != nil {
		return nil, e
	}
	out, e := ParseAppYaml(d)
	out.Path, _ = filepath.Abs(filepath.Dir(f))
	return out, e
}
