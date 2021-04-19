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
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

// AppDir defines the app directory.
const AppDir = "/app"

// App defines an application.
type App struct {
	Path          string
	Name          string                `yaml:"name"`
	Type          string                `yaml:"type"`
	Size          string                `yaml:"size"`
	Disk          int                   `yaml:"disk"`
	Build         AppBuild              `yaml:"build" json:"build"`
	Variables     Variables             `yaml:"variables"`
	Relationships map[string]string     `yaml:"relationships"`
	Web           AppWeb                `yaml:"web"`
	Mounts        map[string]*AppMount  `yaml:"mounts" json:"mounts"`
	Hooks         AppHooks              `yaml:"hooks" json:"hooks"`
	Crons         map[string]*AppCron   `yaml:"crons" json:"crons"`
	Dependencies  AppDependencies       `yaml:"dependencies"`
	Runtime       AppRuntime            `yaml:"runtime"`
	Workers       map[string]*AppWorker `yaml:"workers" json:"workers"`
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
		fmt.Sprintf("app.%s.type", d.Name),
	); err != nil {
		o = append(o, err)
	}
	if d.Disk < 128 {
		o = append(o, NewValidateError(
			fmt.Sprintf("app.%s.disk", d.Name),
			"should be 128 or higher",
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
	for _, key := range []string{"env", "php"} {
		val := d.Variables.GetSubMap(key)
		if val == nil {
			continue
		}
		for sk, sv := range val {
			switch sv.(type) {
			case string, int, float32, float64, bool:
				{
					break
				}
			default:
				{
					o = append(o, NewValidateError(
						fmt.Sprintf("app.%s.variables.%s.%s", d.Name, key, sk),
						"should be a scalar",
					))
					break
				}
			}
		}
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

// ParseAppYamls parses multiple .platform.app.yaml contents and merges them in to one.
func ParseAppYamls(d [][]byte, global *GlobalConfig) (*App, error) {
	o := &App{
		Crons:         make(map[string]*AppCron),
		Mounts:        make(map[string]*AppMount),
		Workers:       make(map[string]*AppWorker),
		Relationships: make(map[string]string),
		Variables:     make(Variables),
	}
	if err := mergeYamls(d, o); err != nil {
		return &App{}, tracerr.Wrap(err)
	}
	// set defaults
	o.SetDefaults()
	// transfer app attributes to its workers
	for name, w := range o.Workers {
		w.Name = name
		w.Type = o.Type
		w.Runtime = o.Runtime
		w.Dependencies = o.Dependencies
		w.Variables = make(Variables)
		mergeMaps(w.Variables, o.Variables)
	}
	// merge global variables
	if global != nil {
		output.LogDebug(fmt.Sprintf("Merge app '%s' variables with global variables.", o.Name), nil)
		o.Variables.Merge(global.Variables)
	}
	return o, nil
}

// ParseAppYamlFiles parses multiple .platform.app.yaml files and merges them in to one.
func ParseAppYamlFiles(fileList []string, global *GlobalConfig) (*App, error) {
	done := output.Duration(
		fmt.Sprintf("Parse app at '%s.'", strings.Join(fileList, ", ")),
	)
	byteList := make([][]byte, 0)
	for _, f := range fileList {
		projectPlatformDir = filepath.Dir(f)
		d, err := ioutil.ReadFile(f)
		if err != nil {
			return &App{}, tracerr.Wrap(err)
		}
		byteList = append(byteList, d)
	}
	a, err := ParseAppYamls(byteList, global)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	// append path to app def
	a.Path, _ = filepath.Abs(filepath.Dir(fileList[0]))
	for _, w := range a.Workers {
		w.Path = a.Path
	}
	done()
	return a, nil
}
