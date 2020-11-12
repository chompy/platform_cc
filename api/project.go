package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/martinlindhe/base36"
)

const appYamlFilename = ".platform.app.yaml"
const routesYamlPath = ".platform/routes.yaml"
const servicesYamlPath = ".platform/services.yaml"
const projectJSONFilename = ".platform_cc.json"

// Project - platform.sh/cc project
type Project struct {
	ID            string                       `json:"id"`
	Path          string                       `json:"-"`
	Apps          []*AppDef                    `json:"-"`
	Routes        []*RouteDef                  `json:"-"`
	Services      []*ServiceDef                `json:"-"`
	Variables     map[string]map[string]string `json:"vars"`
	relationships []map[string]interface{}
	docker        dockerClient
}

// LoadProjectFromPath - load a project from its path
func LoadProjectFromPath(path string, parseYaml bool) (*Project, error) {
	log.Printf("Load project at '%s.'", path)
	var err error
	// read app yaml
	apps := make([]*AppDef, 0)
	if parseYaml {
		appYamlFiles := scanPlatformAppYaml(path)
		if len(appYamlFiles) == 0 {
			return nil, missingAppYamlError{path: path}
		}
		for i := range appYamlFiles {
			app, err := ParseAppYamlFile(appYamlFiles[i])
			if err != nil {
				return nil, err
			}
			apps = append(apps, app)
		}
	}
	// read services yaml
	services := []*ServiceDef{}
	if parseYaml {
		fullServiceYamlPath := filepath.Join(path, servicesYamlPath)
		services, err = ParseServicesYamlFile(fullServiceYamlPath)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	// read routes yaml
	routes := make([]*RouteDef, 0)
	if parseYaml {
		fullRouteYamlPath := filepath.Join(path, routesYamlPath)
		routes, err = ParseRoutesYamlFile(fullRouteYamlPath)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}
	path, _ = filepath.Abs(path)
	dockerClient, err := newDockerClient()
	if err != nil {
		return nil, err
	}
	o := &Project{
		ID:            "",
		Path:          path,
		Apps:          apps,
		Services:      services,
		Routes:        routes,
		Variables:     make(map[string]map[string]string),
		docker:        dockerClient,
		relationships: make([]map[string]interface{}, 0),
	}
	o.Load()
	if o.ID == "" {
		o.ID = generateProjectID()
		o.Save()
	}
	log.Printf("Project '%s' loaded.", o.ID)
	return o, nil
}

func scanPlatformAppYaml(path string) []string {
	o := make([]string, 0)
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f.Name() == appYamlFilename {
			o = append(o, path)
		}
		return nil
	})
	return o
}

func generateProjectID() string {
	return strings.ToLower(base36.Encode(uint64(time.Now().Unix())))
}

// Load - load project info from file
func (p *Project) Load() error {
	data, err := ioutil.ReadFile(
		filepath.Join(p.Path, projectJSONFilename),
	)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	return json.Unmarshal(data, p)
}

// Save - save project info to file
func (p *Project) Save() error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(
		filepath.Join(p.Path, projectJSONFilename),
		data,
		0655,
	)
}

// Start - start project
func (p *Project) Start() error {
	log.Printf("Start project '%s.'", p.ID)
	// create network
	if err := p.docker.CreateProjectNetwork(p.ID); err != nil {
		return err
	}
	// services
	for _, service := range p.Services {
		if err := p.startService(service); err != nil {
			return err
		}
		if err := p.openService(service); err != nil {
			return err
		}
	}
	// app
	for _, app := range p.Apps {
		if err := p.startApp(app); err != nil {
			return err
		}
		if err := p.openApp(app); err != nil {
			return err
		}
	}
	log.Println("Project started.")
	return nil
}

// Stop - stop project
func (p *Project) Stop() error {
	log.Printf("Stop project '%s.'", p.ID)
	p.docker.DeleteProjectContainers(p.ID)
	p.docker.DeleteProjectNetwork(p.ID)
	log.Println("Project stopped.")
	return nil
}

// Build - build applications in project
func (p *Project) Build() error {
	log.Printf("Build project '%s.'", p.ID)
	// app
	for _, app := range p.Apps {
		if err := p.BuildApp(app); err != nil {
			return err
		}
	}
	return nil
}

// Purge - purge project
func (p *Project) Purge() error {
	if err := p.Stop(); err != nil {
		return err
	}
	log.Printf("Purge project '%s.'", p.ID)
	if err := p.docker.DeleteProjectVolumes(p.ID); err != nil {
		return err
	}
	os.Remove(
		filepath.Join(p.Path, projectJSONFilename),
	)
	log.Println("Project purged.")
	return nil
}

// getUID - get os user's uid and gid
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
