package api

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v3"
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
	Hooks         AppHooksDef                       `yaml:"hooks" json:"hooks"`
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

// BuildPlatformApplicationVar - build PLATFORM_APPLICATION env var
func (d *AppDef) BuildPlatformApplicationVar() string {
	jsonData, _ := json.Marshal(map[string]interface{}{
		"resources":     nil,
		"size":          "AUTO",
		"disk":          d.Disk,
		"access":        map[string]string{},
		"relationships": d.Relationships,
		"mounts":        d.Mounts,
		"timezone":      nil,
		"variables":     d.Variables,
		"firewall":      nil,
		"name":          d.Name,
		"type":          d.Type,
		"runtime":       d.Runtime,
		"preflight": map[string]interface{}{
			"enabled":       true,
			"ignored_rules": []string{},
		},
		"tree_id":      "-",
		"slug_id":      "-",
		"app_dir":      appDir,
		"web":          d.Web,
		"hook":         d.Hooks,
		"crons":        d.Crons,
		"dependencies": d.Dependencies,
	})
	return base64.StdEncoding.EncodeToString(jsonData)
}

// GetEmptyRelationship - get empty relationship
func (d AppDef) GetEmptyRelationship() map[string]interface{} {
	return map[string]interface{}{
		"host":        "",
		"hostname":    "",
		"ip":          "",
		"port":        80,
		"path":        "",
		"scheme":      d.Web.Upstream.Protocol,
		"fragment":    nil,
		"rel":         d.Web.Upstream.Protocol,
		"host_mapped": false,
		"public":      false,
		"type":        d.Type,
		"service":     d.Name,
	}
}

// ParseAppYaml - parse app yaml
func ParseAppYaml(d []byte) (*AppDef, error) {
	o := &AppDef{
		Crons:         make(map[string]*AppCronDef),
		Mounts:        make(map[string]*AppMountDef),
		Relationships: make(map[string]string),
		Variables:     make(map[string]map[string]interface{}),
	}
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
