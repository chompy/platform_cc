package def

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// AppDir defines the app directory.
const AppDir = "/app"

// App defines an application.
type App struct {
	Path          string
	Name          string                            `yaml:"name"`
	Type          string                            `yaml:"type"`
	Size          string                            `yaml:"size"`
	Disk          int                               `yaml:"disk"`
	Build         AppBuild                          `yaml:"build"`
	Variables     map[string]map[string]interface{} `yaml:"variables"`
	Relationships map[string]string                 `yaml:"relationships"`
	Web           AppWeb                            `yaml:"web"`
	Mounts        map[string]*AppMount              `yaml:"mounts" json:"mounts"`
	Hooks         AppHooks                          `yaml:"hooks" json:"hooks"`
	Crons         map[string]*AppCron               `yaml:"crons" json:"crons"`
	Dependencies  AppDependencies                   `yaml:"dependencies"`
	Runtime       AppRuntime                        `yaml:"runtime"`
}

// SetDefaults sets the default values.
func (d *App) SetDefaults() {
	if d.Name == "" {
		d.Name = "app"
	}
	if d.Type == "" {
		d.Type = "php:7.4"
	}
	if d.Size == "" {
		d.Size = "M"
	}
	if d.Disk == 0 {
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

// Validate checks for errors.
func (d App) Validate() []error {
	o := make([]error, 0)
	if err := validateMustContainOne(
		[]string{"php", "golang", "dotnet", "elixir", "java", "lisp", "nodejs", "python", "ruby"},
		d.GetTypeName(),
		"app.type",
	); err != nil {
		o = append(o, err)
	}
	if d.Disk < 256 {
		o = append(o, NewDefValidateError(
			"should be 256 or higher",
			"app.disk",
		))
	}
	if e := d.Build.Validate(&d); len(e) > 0 {
		o = append(o, e...)
	}
	if e := d.Web.Validate(&d); len(e) > 0 {
		o = append(o, e...)
	}
	for _, m := range d.Mounts {
		if e := m.Validate(&d); len(e) > 0 {
			o = append(o, e...)
		}
	}
	if e := d.Hooks.Validate(&d); len(e) > 0 {
		o = append(o, e...)
	}
	for _, c := range d.Crons {
		if e := c.Validate(&d); len(e) > 0 {
			o = append(o, e...)
		}
	}
	if e := d.Dependencies.Validate(&d); len(e) > 0 {
		o = append(o, e...)
	}
	if e := d.Runtime.Validate(&d); len(e) > 0 {
		o = append(o, e...)
	}
	return o
}

// GetTypeName gets the type of app.
func (d App) GetTypeName() string {
	return strings.Split(d.Type, ":")[0]
}

// BuildPlatformApplicationVar builds the PLATFORM_APPLICATION env var.
func (d *App) BuildPlatformApplicationVar() string {
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
		"app_dir":      AppDir,
		"web":          d.Web,
		"hook":         d.Hooks,
		"crons":        d.Crons,
		"dependencies": d.Dependencies,
	})
	return base64.StdEncoding.EncodeToString(jsonData)
}

// GetEmptyRelationship returns an empty relationship.
func (d App) GetEmptyRelationship() map[string]interface{} {
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

// ParseAppYaml parses the contents of a .platform.app.yaml file.
func ParseAppYaml(d []byte) (*App, error) {
	o := &App{
		Crons:         make(map[string]*AppCron),
		Mounts:        make(map[string]*AppMount),
		Relationships: make(map[string]string),
		Variables:     make(map[string]map[string]interface{}),
	}
	e := yaml.Unmarshal(d, &o)
	o.SetDefaults()
	return o, e
}

// ParseAppYamlFile opens the .platform.app.yaml file and parses it.
func ParseAppYamlFile(f string) (*App, error) {
	log.Printf("Parse app at '%s.'", f)
	d, e := ioutil.ReadFile(f)
	if e != nil {
		return nil, e
	}
	out, e := ParseAppYaml(d)
	out.Path, _ = filepath.Abs(filepath.Dir(f))
	return out, e
}
