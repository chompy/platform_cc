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
	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/config"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/container"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/def"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/output"
	"gitlab.com/contextualcode/platform_cc/v2/pkg/platformsh"
)

var appYamlFilenames = []string{".platform.app.yaml", ".platform.app.pcc.yaml"}
var serviceYamlFilenames = []string{".platform/services.yaml", ".platform/services.pcc.yaml"}
var routesYamlFilenames = []string{".platform/routes.yaml", ".platform/routes.pcc.yaml"}

const projectJSONFilename = ".platform_cc.json"
const platformShDockerImagePrefix = "docker.registry.platform.sh/"

// Project defines a platform.sh/cc project.
type Project struct {
	ID               string            `json:"id"`
	Path             string            `json:"-"`
	Apps             []def.App         `json:"-"`
	Routes           []def.Route       `json:"-"`
	Services         []def.Service     `json:"-"`
	Variables        def.Variables     `json:"vars"`
	Flags            Flags             `json:"flags"`   // local project flags
	Options          map[Option]string `json:"options"` // local project options
	relationships    []map[string]interface{}
	containerHandler container.Interface
	globalConfig     def.GlobalConfig
	PlatformSH       *platformsh.Project `json:"-"`
	slot             int                 // set volume slot
	noCommit         bool                // flag that signifies apps should not be committed
	noBuild          bool                // flag that signifies apps should not be built on start up
}

// LoadFromPath loads a project from its path.
func LoadFromPath(path string, parseYaml bool) (*Project, error) {
	done := output.Duration(
		fmt.Sprintf("Search for project at '%s.'", path),
	)
	// look for a psh project path (git repo with remote pointed at platform.sh)
	psh, err := platformsh.LoadProjectFromPath(path)
	if err == nil {
		path = psh.LocalPath
	} else {
		output.LogDebug("Platform.sh project not found.", err)
	}
	// global config
	gc, err := config.Load()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// build project
	path, _ = filepath.Abs(path)
	// TODO allow container handler to be configured
	containerHandler, err := container.NewDocker()
	if err != nil {
		return nil, err
	}
	o := &Project{
		ID:               "",
		Path:             path,
		Variables:        make(map[string]interface{}),
		Options:          make(map[Option]string),
		PlatformSH:       psh,
		containerHandler: containerHandler,
		relationships:    make([]map[string]interface{}, 0),
		slot:             1,
		globalConfig:     gc,
	}
	o.Load()
	if o.ID == "" {
		if psh != nil && psh.ID != "" {
			o.ID = psh.ID
		} else {
			o.ID = generateProjectID()
		}
		o.Save()
	}
	// read app yaml
	apps := make([]def.App, 0)
	if parseYaml {
		appYamlFiles := scanPlatformAppYaml(path, o.HasFlag(DisableYamlOverrides))
		if len(appYamlFiles) == 0 {
			return nil, errors.WithStack(fmt.Errorf("could not locate app yaml file"))
		}
		for _, appYamlFileList := range appYamlFiles {
			app, err := def.ParseAppYamlFiles(appYamlFileList, &gc)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			apps = append(apps, *app)
		}
		o.Apps = apps
	}
	// read services yaml
	if parseYaml {
		serviceYamlPaths := make([]string, 0)
		for _, fn := range serviceYamlFilenames {
			serviceYamlPaths = append(
				serviceYamlPaths,
				filepath.Join(path, fn),
			)
			if o.HasFlag(DisableYamlOverrides) {
				break
			}
		}
		o.Services, err = def.ParseServiceYamlFiles(serviceYamlPaths)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	// read routes yaml
	if parseYaml {
		routesYamlPaths := make([]string, 0)
		for _, fn := range routesYamlFilenames {
			routesYamlPaths = append(
				routesYamlPaths,
				filepath.Join(path, fn),
			)
			if o.HasFlag(DisableYamlOverrides) {
				break
			}
		}
		o.Routes, err = def.ParseRoutesYamlFiles(routesYamlPaths)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		o.Routes, err = def.ExpandRoutes(
			o.Routes,
			o.GetOption(OptionDomainSuffix),
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if !parseYaml {
		output.Info("Skipped (parseYaml=false).")
	}
	o.setAppFlags()
	done()
	output.Info(fmt.Sprintf("Loaded project '%s.'", o.ID))
	return o, nil
}

func scanPlatformAppYaml(topPath string, disableOverrides bool) [][]string {
	o := make([][]string, 0)
	appYamlPaths := make([]string, 0)
	filepath.Walk(topPath, func(path string, f os.FileInfo, err error) error {
		// check sub directory
		if f.IsDir() && f.Name() != "." && path != topPath {
			for _, appYamlFilename := range appYamlFilenames {
				possiblePath := filepath.Join(path, appYamlFilename)
				if _, err := os.Stat(possiblePath); !os.IsNotExist(err) {
					appYamlPaths = append(appYamlPaths, possiblePath)
				}
			}
			return filepath.SkipDir
		}
		// check root directory
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
						if !disableOverrides {
							o[i] = append(o[i], appYamlPath)
						}
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
		return errors.WithStack(err)
	}
	return json.Unmarshal(data, p)
}

// Save saves the project info to file.
func (p *Project) Save() error {
	output.LogDebug("Save project.", p.ID)
	data, err := json.Marshal(p)
	if err != nil {
		return errors.WithStack(err)
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
		return errors.WithStack(err)
	}
	// build list of services (apps, services, workers) that need to start
	serviceList := make([]interface{}, 0)
	for _, s := range p.Services {
		// ignore network storage service as it is not needed by pcc
		if s.GetTypeName() == "network-storage" {
			continue
		}
		serviceList = append(serviceList, s)
	}
	for _, a := range p.Apps {
		serviceList = append(serviceList, a)
		if p.HasFlag(EnableWorkers) {
			for _, w := range a.Workers {
				serviceList = append(serviceList, w)
			}
		}
	}
	// determine start order
	var err error
	serviceList, err = p.GetDefinitionStartOrder(serviceList)
	if err != nil {
		return errors.WithStack(err)
	}
	// itterate and start
	for _, service := range serviceList {
		// start
		c := p.NewContainer(service)
		if err := c.Start(); err != nil {
			if errors.Is(err, ErrContainerRunning) {
				output.Info(fmt.Sprintf("Container '%s' is already running.", c.Config.GetContainerName()))
				output.IndentLevel--
				continue
			}
			return errors.WithStack(err)
		}
		// container type specific operations
		switch c.Config.ObjectType {
		case container.ObjectContainerApp, container.ObjectContainerWorker:
			{
				// build
				if !p.noBuild && !c.HasBuild() {
					if err := c.Build(); err != nil {
						return errors.WithStack(err)
					}
					if !p.noCommit {
						if err := c.Commit(); err != nil {
							return errors.WithStack(err)
						}
					}
				}
				// setup mounts
				if err := c.SetupMounts(); err != nil {
					return errors.WithStack(err)
				}
			}
		}
		// open + process relationships
		rels, err := c.Open()
		if err != nil {
			return errors.WithStack(err)
		}
		p.relationships = append(p.relationships, rels...)
	}
	// post-deploy
	for _, service := range serviceList {
		c := p.NewContainer(service)
		if err := c.PostDeploy(); err != nil {
			return errors.WithStack(err)
		}
	}
	done()
	return nil
}

// Stop stops the project.
func (p *Project) Stop() error {
	done := output.Duration(
		fmt.Sprintf("Stop project '%s.'", p.ID),
	)
	if err := p.containerHandler.ProjectStop(p.ID); err != nil {
		return errors.WithStack(err)
	}
	done()
	return nil
}

// Build builds all applications in the project.
func (p *Project) Build(force bool) error {
	if p.noBuild {
		return nil
	}
	done := output.Duration(
		fmt.Sprintf("Build project '%s.'", p.ID),
	)
	// app
	for _, app := range p.Apps {
		c := p.NewContainer(app)
		if c.HasBuild() && !force {
			output.Info("Already built, skipped.")
			output.LogDebug(
				fmt.Sprintf("Skip build for %s, already committed, not forced.", c.Config.GetContainerName()),
				nil,
			)
			continue
		}
		if err := c.Build(); err != nil {
			return errors.WithStack(err)
		}
		if !p.noCommit {
			if err := c.Commit(); err != nil {
				return errors.WithStack(err)
			}
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
			return errors.WithStack(err)
		}
	}
	done()
	return nil
}

// PostDeploy runs post-deploy hooks for all applications in the project.
func (p *Project) PostDeploy() error {
	done := output.Duration(
		fmt.Sprintf("Run post-deploy hooks for project '%s.'", p.ID),
	)
	// app
	for _, app := range p.Apps {
		c := p.NewContainer(app)
		if err := c.PostDeploy(); err != nil {
			return errors.WithStack(err)
		}
	}
	done()
	return nil
}

// Purge purges data related to the project.
func (p *Project) Purge() error {
	done := output.Duration(fmt.Sprintf("Purge project '%s.'", p.ID))
	if err := p.containerHandler.ProjectPurge(p.ID); err != nil {
		return errors.WithStack(err)
	}
	// TODO not sure if purge should really delete platform_cc.json
	/*if err := os.Remove(
		filepath.Join(p.Path, projectJSONFilename),
	); err != nil {
		output.Warn(err.Error())
	}*/
	done()
	return nil
}

// PurgeSlot purges volumes for current slot.
func (p *Project) PurgeSlot() error {
	done := output.Duration(fmt.Sprintf("Purge project '%s' slot %d.", p.ID, p.slot))
	if err := p.containerHandler.ProjectPurgeSlot(p.ID, p.slot); err != nil {
		return errors.WithStack(err)
	}
	done()
	return nil
}

// Pull fetches all the Docker container images needed by the project.
func (p *Project) Pull() error {
	done := output.Duration("Pull images.")
	containerConfigs := make([]container.Config, 0)
	for _, d := range p.Services {
		c := p.NewContainer(d)
		containerConfigs = append(containerConfigs, c.Config)
	}
	for _, d := range p.Apps {
		c := p.NewContainer(d)
		if c.HasBuild() {
			output.LogDebug(
				fmt.Sprintf("Skip pulling image for %s, has committed image.", c.Config.GetContainerName()),
				nil,
			)
			continue
		}
		containerConfigs = append(containerConfigs, c.Config)
	}
	if err := p.containerHandler.ImagePull(containerConfigs); err != nil {
		return errors.WithStack(err)
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
		if !p.HasFlag(EnablePHPOpcache) {
			output.LogDebug("Disable opcache.", nil)
			p.Apps[i].Variables.Set("php:enable.opcache", "off")
			p.Apps[i].Variables.Set("php:opcache.enable", "off")
			for j := range p.Apps[i].Workers {
				p.Apps[i].Workers[j].Variables.Set("php:enable.opcache", "off")
				p.Apps[i].Workers[j].Variables.Set("php:opcache.enable", "off")
			}
		}
	}
}

// SetContainerHandler sets the container handler.
func (p *Project) SetContainerHandler(c container.Interface) {
	p.containerHandler = c
}

// SetGlobalConfig sets the global config, used for testing.
func (p *Project) SetGlobalConfig(gc def.GlobalConfig) {
	p.globalConfig = gc
}

// SetSlot sets the current slot.
func (p *Project) SetSlot(slot int) {
	if p.slot != slot {
		if slot <= 0 {
			slot = 1
		}
		output.Info(fmt.Sprintf("Set slot %d.", slot))
		p.slot = slot
	}
}

// CopySlot copies the current slot to a given destination slot.
func (p *Project) CopySlot(destSlot int) error {
	return errors.WithStack(p.containerHandler.ProjectCopySlot(
		p.ID, p.slot, destSlot,
	))
}

// SetNoCommit sets the no commit flag.
func (p *Project) SetNoCommit() {
	p.noCommit = true
}

// SetNoBuild sets the no build flag.
func (p *Project) SetNoBuild() {
	p.noBuild = true
}

// Validate returns list of validation errors for project.
func (p *Project) Validate() []error {
	done := output.Duration("Validate project.")
	out := make([]error, 0)
	for _, app := range p.Apps {
		out = append(out, app.Validate()...)
		for _, rel := range app.Relationships {
			hasRel := false
			relSplit := strings.Split(rel, ":")
			for _, serv := range p.Services {
				if relSplit[0] == serv.Name {
					hasRel = true
					break
				}
			}
			if !hasRel {
				out = append(out, def.NewValidateError(
					fmt.Sprintf("app.%s.relationships", app.Name),
					fmt.Sprintf("service %s is not defined", relSplit[0]),
				))
			}
		}
	}
	for _, serv := range p.Services {
		out = append(out, serv.Validate()...)
	}
	for _, route := range p.Routes {
		out = append(out, route.Validate()...)
	}
	done()
	return out
}