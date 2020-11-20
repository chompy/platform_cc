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

package router

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/docker"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

// GetContainerConfig gets container configuration for the router.
func GetContainerConfig() docker.ContainerConfig {
	return docker.ContainerConfig{
		ProjectID:  "_",
		ObjectName: "router",
		ObjectType: docker.ObjectContainerRouter,
		Command:    []string{"nginx", "-g", "daemon off;"},
		Image:      "docker.io/library/nginx:1.19-alpine",
		WorkingDir: "/routes",
		Ports: []string{
			"80:80/tcp",
		},
	}
}

// Start starts the router.
func Start() error {
	log.Println("Start main router.")
	d, err := docker.NewClient()
	if err != nil {
		return tracerr.Wrap(err)
	}
	// create network (if not already created)
	if err := d.CreateNetwork(); err != nil {
		return tracerr.Wrap(err)
	}
	// start container
	containerConf := GetContainerConfig()
	if err := d.StartContainer(containerConf); err != nil {
		return tracerr.Wrap(err)
	}
	// upload nginx conf
	nginxConfReader := bytes.NewReader([]byte(nginxBaseConf))
	if err := d.UploadDataToContainer(
		containerConf.GetContainerName(),
		"/etc/nginx/nginx.conf",
		nginxConfReader,
	); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// Stop stops the router.
func Stop() error {
	log.Println("Stop main router.")
	d, err := docker.NewClient()
	if err != nil {
		return tracerr.Wrap(err)
	}
	return tracerr.Wrap(d.DeleteProjectContainers("router"))
}

// Reload issues reload command to nginx in router container.
func Reload() error {
	d, err := docker.NewClient()
	if err != nil {
		return tracerr.Wrap(err)
	}
	containerConf := GetContainerConfig()
	return tracerr.Wrap(d.RunContainerCommand(
		containerConf.GetContainerName(),
		"root",
		[]string{"nginx", "-s", "reload"},
		nil,
	))
}

// AddProjectRoutes adds given project's routes to router.
func AddProjectRoutes(p *project.Project) error {
	log.Printf("Add routes for project '%s.'", p.ID)
	d, err := docker.NewClient()
	if err != nil {
		return tracerr.Wrap(err)
	}
	containerConf := GetContainerConfig()
	// generate route nginx
	routerNginxConf, err := GenerateNginxConfig(p)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// upload to container
	configReader := bytes.NewReader(routerNginxConf)
	if err := d.UploadDataToContainer(
		containerConf.GetContainerName(),
		fmt.Sprintf("/routes/%s.conf", p.ID),
		configReader,
	); err != nil {
		return tracerr.Wrap(err)
	}
	return tracerr.Wrap(Reload())
}

// DeleteProjectRoutes deletes routes for given project.
func DeleteProjectRoutes(p *project.Project) error {
	log.Printf("Delete routes for project '%s.'", p.ID)
	d, err := docker.NewClient()
	if err != nil {
		return tracerr.Wrap(err)
	}
	containerConf := GetContainerConfig()
	// delete config file
	if err := d.RunContainerCommand(
		containerConf.GetContainerName(),
		"root",
		[]string{"rm", "-rf", fmt.Sprintf("/routes/%s.conf", p.ID)},
		nil,
	); err != nil {
		if !strings.Contains(err.Error(), "No such container") {
			return tracerr.Wrap(err)
		}
		return nil
	}
	return tracerr.Wrap(Reload())
}
