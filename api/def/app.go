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

package def

import (
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
	Workers       map[string]*AppWorker             `yaml:"workers" json:"workers"`
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
		o = append(o, NewValidateError(
			"app.disk",
			"should be 256 or higher",
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
func ParseAppYaml(d []byte, append *App) (*App, error) {
	o := &App{
		Crons:         make(map[string]*AppCron),
		Mounts:        make(map[string]*AppMount),
		Workers:       make(map[string]*AppWorker),
		Relationships: make(map[string]string),
		Variables:     make(map[string]map[string]interface{}),
	}
	if append != nil {
		o = append
	}
	e := yaml.Unmarshal(d, o)
	o.SetDefaults()
	for name, w := range o.Workers {
		w.Name = name
		w.Type = o.Type
		w.Runtime = o.Runtime
		w.Dependencies = o.Dependencies
	}
	return o, e
}

// ParseAppYamlFile opens the .platform.app.yaml file and parses it.
func ParseAppYamlFile(f string, append *App) (*App, error) {
	log.Printf("Parse app at '%s.'", f)
	d, e := ioutil.ReadFile(f)
	if e != nil {
		return nil, e
	}
	o, e := ParseAppYaml(d, append)
	o.Path, _ = filepath.Abs(filepath.Dir(f))
	for _, w := range o.Workers {
		w.Path = o.Path
	}
	return o, e
}
