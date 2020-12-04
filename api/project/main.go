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

// Package project provides tools to run a Platform.sh local development environment.
package project

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/martinlindhe/base36"
	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/docker"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

var appYamlFilenames = []string{".platform.app.yaml", ".platform.app.pcc.yaml"}

const routesYamlPath = ".platform/routes.yaml"
const servicesYamlPath = ".platform/services.yaml"
const projectJSONFilename = ".platform_cc.json"
const platformShDockerImagePrefix = "docker.registry.platform.sh/"

// Project defines a platform.sh/cc project.
type Project struct {
	ID            string                       `json:"id"`
	Path          string                       `json:"-"`
	Apps          []def.App                    `json:"-"`
	Routes        []def.Route                  `json:"-"`
	Services      []def.Service                `json:"-"`
	Variables     map[string]map[string]string `json:"vars"`
	Flags         Flag                         `json:"flags"`
	Options       map[Option]string            `json:"options"`
	relationships []map[string]interface{}
	docker        docker.Client
}

// LoadFromPath loads a project from its path.
func LoadFromPath(path string, parseYaml bool) (*Project, error) {
	done := output.Duration(
		fmt.Sprintf("Load project at '%s.'", path),
	)
	var err error
	// build project
	path, _ = filepath.Abs(path)
	DockerClient, err := docker.NewClient()
	if err != nil {
		return nil, err
	}
	o := &Project{
		ID:            "",
		Path:          path,
		Variables:     make(map[string]map[string]string),
		Options:       make(map[Option]string),
		docker:        DockerClient,
		relationships: make([]map[string]interface{}, 0),
	}
	o.Load()
	if o.ID == "" {
		o.ID = generateProjectID()
		o.Save()
	}
	// read app yaml
	apps := make([]def.App, 0)
	if parseYaml {
		appYamlFiles := scanPlatformAppYaml(path)
		if len(appYamlFiles) == 0 {
			return nil, tracerr.Wrap(fmt.Errorf("could not location app yaml file"))
		}
		for _, appYamlFileList := range appYamlFiles {
			app, err := def.ParseAppYamlFiles(appYamlFileList)
			if err != nil {
				return nil, tracerr.Wrap(err)
			}
			apps = append(apps, *app)
		}
		o.Apps = apps
	}
	// read services yaml
	services := []def.Service{}
	if parseYaml {
		fullServiceYamlPath := filepath.Join(path, servicesYamlPath)
		services, err = def.ParseServicesYamlFile(fullServiceYamlPath)
		if err != nil && !os.IsNotExist(err) {
			return nil, tracerr.Wrap(err)
		}
		o.Services = services
	}
	// read routes yaml
	routes := make([]def.Route, 0)
	if parseYaml {
		fullRouteYamlPath := filepath.Join(path, routesYamlPath)
		routes, err = def.ParseRoutesYamlFile(fullRouteYamlPath)
		if err != nil && !os.IsNotExist(err) {
			return nil, tracerr.Wrap(err)
		}
		routes, err = def.ExpandRoutes(
			routes,
			OptionDomainSuffix.Value(o.Options),
		)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		o.Routes = routes
	}
	if !parseYaml {
		output.Info("Skipped (parseYaml=false).")
	}
	o.setAppFlags()
	done()
	output.Info(fmt.Sprintf("Loaded project '%s.'", o.ID))
	return o, nil
}

func scanPlatformAppYaml(path string) [][]string {
	o := make([][]string, 0)
	appYamlPaths := make([]string, 0)
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		for _, appYamlFilename := range appYamlFilenames {
			if f.Name() == appYamlFilename {
				appYamlPaths = append(appYamlPaths, path)
			}
		}
		return nil
	})
	for _, appYamlFilename := range appYamlFilenames {
		for _, appYamlPath := range appYamlPaths {
			if strings.HasSuffix(appYamlPath, appYamlFilename) {
				hasOut := false
				for i := range o {
					if filepath.Dir(o[i][0]) == filepath.Dir(appYamlPath) {
						o[i] = append(o[i], appYamlPath)
						hasOut = true
					}
				}
				if !hasOut {
					oo := make([]string, 1)
					oo[0] = appYamlPath
					o = append(o, oo)
				}
			}
		}
	}
	return o
}

func generateProjectID() string {
	return strings.ToLower(base36.Encode(uint64(time.Now().Unix())))
}

// Load loads the project info from file.
func (p *Project) Load() error {
	output.LogDebug("Load project.", p.ID)
	data, err := ioutil.ReadFile(
		filepath.Join(p.Path, projectJSONFilename),
	)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return tracerr.Wrap(err)
	}
	return json.Unmarshal(data, p)
}

// Save saves the project info to file.
func (p *Project) Save() error {
	output.LogDebug("Save project.", p.ID)
	data, err := json.Marshal(p)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return ioutil.WriteFile(
		filepath.Join(p.Path, projectJSONFilename),
		data,
		0655,
	)
}

// Start starts the project.
func (p *Project) Start() error {
	done := output.Duration("Start project.")
	// pull images
	if err := p.Pull(); err != nil {
		return tracerr.Wrap(err)
	}
	// create network (if not already created)
	if err := p.docker.CreateNetwork(); err != nil {
		return tracerr.Wrap(err)
	}
	// services
	for _, service := range p.Services {
		// skip if service requires relationships (start after apps start)
		// this is needed by varnish
		if len(service.Relationships) > 0 {
			continue
		}
		// start
		c := p.NewContainer(service)
		if err := c.Start(); err != nil {
			return tracerr.Wrap(err)
		}
		// open + process relationships
		rels, err := c.Open()
		if err != nil {
			return tracerr.Wrap(err)
		}
		p.relationships = append(p.relationships, rels...)
	}
	// app
	for _, app := range p.Apps {
		// prepare nfs volume (if needed)
		if err := prepareNfsVolume(p, app); err != nil {
			return tracerr.Wrap(err)
		}
		// start
		c := p.NewContainer(app)
		if err := c.Start(); err != nil {
			return tracerr.Wrap(err)
		}
		// setup mounts
		if err := c.SetupMounts(); err != nil {
			return tracerr.Wrap(err)
		}
		// open + process relationships
		rels, err := c.Open()
		if err != nil {
			return tracerr.Wrap(err)
		}
		p.relationships = append(p.relationships, rels...)

		// workers
		if p.Flags.Has(EnableWorkers) {
			for _, worker := range app.Workers {
				wc := p.NewContainer(*worker)
				// prepare nfs volume (if needed)
				if err := prepareNfsVolume(p, *worker); err != nil {
					return tracerr.Wrap(err)
				}
				// start
				if err := wc.Start(); err != nil {
					return tracerr.Wrap(err)
				}
				// setup mounts
				if err := wc.SetupMounts(); err != nil {
					return tracerr.Wrap(err)
				}
				// open
				_, err := c.Open()
				if err != nil {
					return tracerr.Wrap(err)
				}

			}
		}
	}
	// services with relationships (varnish, etc)
	for _, service := range p.Services {
		// skip if service has no relationships, it was already started
		if len(service.Relationships) == 0 {
			continue
		}
		// start
		c := p.NewContainer(service)
		if err := c.Start(); err != nil {
			return tracerr.Wrap(err)
		}
		// open + process relationships
		rels, err := c.Open()
		if err != nil {
			return tracerr.Wrap(err)
		}
		p.relationships = append(p.relationships, rels...)
	}
	done()
	return nil
}

// Stop stops the project.
func (p *Project) Stop() error {
	done := output.Duration(
		fmt.Sprintf("Stop project '%s.'", p.ID),
	)
	if err := p.docker.DeleteProjectContainers(p.ID); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// Build builds all applications in the project.
func (p *Project) Build() error {
	done := output.Duration(
		fmt.Sprintf("Build project '%s.'", p.ID),
	)
	// app
	for _, app := range p.Apps {
		c := p.NewContainer(app)
		if err := c.Build(); err != nil {
			return tracerr.Wrap(err)
		}
	}
	done()
	return nil
}

// Deploy runs deploy hooks for all applications in the project.
func (p *Project) Deploy() error {
	done := output.Duration(
		fmt.Sprintf("Run deploy hooks for project '%s.'", p.ID),
	)
	// app
	for _, app := range p.Apps {
		c := p.NewContainer(app)
		if err := c.Deploy(); err != nil {
			return tracerr.Wrap(err)
		}
	}
	done()
	return nil
}

// Purge purges data related to the project.
func (p *Project) Purge() error {
	if err := p.Stop(); err != nil {
		return tracerr.Wrap(err)
	}
	done := output.Duration(fmt.Sprintf("Purge project '%s.'", p.ID))
	if err := p.docker.DeleteProjectVolumes(p.ID); err != nil {
		return tracerr.Wrap(err)
	}
	if err := os.Remove(
		filepath.Join(p.Path, projectJSONFilename),
	); err != nil {
		output.Warn(err.Error())
	}
	done()
	return nil
}

// Pull fetches all the Docker container images needed by the project.
func (p *Project) Pull() error {
	done := output.Duration("Pull images.")
	containerConfigs := make([]docker.ContainerConfig, 0)
	for _, d := range p.Services {
		c := p.NewContainer(d)
		containerConfigs = append(containerConfigs, c.Config)
	}
	for _, d := range p.Apps {
		c := p.NewContainer(d)
		containerConfigs = append(containerConfigs, c.Config)
	}
	if err := p.docker.PullImages(containerConfigs); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// getUID gets the os user's uid and gid.
func (p *Project) getUID() (int, int) {
	uid := 0
	gid := 0
	currentUser, _ := user.Current()
	if currentUser != nil {
		uid, _ = strconv.Atoi(currentUser.Uid)
		gid, _ = strconv.Atoi(currentUser.Gid)
	}
	if uid == 0 {
		uid = 1000
	}
	if gid == 0 {
		gid = 1000
	}
	return uid, gid
}

func (p *Project) setAppFlags() {
	for i := range p.Apps {
		// disable opcache if flag not present
		if !p.Flags.Has(EnablePHPOpcache) {
			p.Apps[i].Variables["php"]["enable.opcache"] = "0"
			for j := range p.Apps[i].Workers {
				p.Apps[i].Workers[j].Variables["php"]["enable.opcache"] = "0"
			}
		}
	}
}
