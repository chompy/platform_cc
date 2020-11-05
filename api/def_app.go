package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

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
	envVars := map[string]string{
		"PLATFORM_DOCUMENT_ROOT":    "/app/web",
		"PLATFORM_APPLICATION":      d.BuildPlatformApplicationVar(),
		"PLATFORM_PROJECT":          "-",
		"PLATFORM_PROJECT_ENTROPY":  "1234abc",
		"PLATFORM_APPLICATION_NAME": d.Name,
		"PLATFORM_BRANCH":           "pcc",
		"PLATFORM_DIR":              appDir,
		"PLATFORM_TREE_ID":          "-",
		"PLATFORM_ENVIRONMENT":      "pcc",
		"PLATFORM_VARIABLES":        d.BuildPlatformVariablesVar(),
		"PLATFORM_ROUTES":           "e30=",
	}
	for k, v := range d.Variables["env"] {
		envVars[k] = v.(string)
	}
	return json.Marshal(map[string]interface{}{
		"crons":                 d.Crons,
		"enable_smtp":           "false",
		"mounts":                d.Mounts,
		"cron_minimum_interval": "1",
		"configuration": map[string]interface{}{
			"app_dir":       appDir,
			"hooks":         d.Hooks,
			"variables":     envVars,
			"timezone":      nil,
			"disk":          d.Disk,
			"slug_id":       "-",
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
		"tree_id": "-",
		"slug_id": "-",
		"app_dir": appDir,
		"web":     d.Web,
		"hook":    d.Hooks,
		"crons":   d.Crons,
	})
	return base64.StdEncoding.EncodeToString(jsonData)
}

// BuildPlatformVariablesVar - build PLATFORM_VARIABLES env var
func (d *AppDef) BuildPlatformVariablesVar() string {
	data := make(map[string]string)
	for varType, varVal := range d.Variables {
		for k, v := range varVal {
			switch v.(type) {
			case string:
				{
					data[fmt.Sprintf("%s:%s", strings.ToLower(varType), k)] = v.(string)
					break
				}
			}
		}
	}
	jsonData, _ := json.Marshal(data)
	return base64.StdEncoding.EncodeToString(jsonData)
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
