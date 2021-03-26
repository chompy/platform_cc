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

// Package router is the main router and provides the main HTTP entry point for Platform.CC.
package router

import (
	"bytes"
	"fmt"
	"strings"

	"gitlab.com/contextualcode/platform_cc/api/def"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/container"
	"gitlab.com/contextualcode/platform_cc/api/output"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

// HTTPPort is the port to accept HTTP requests on.
var HTTPPort = uint16(80)

// HTTPSPort is the port to accept HTTPS requests on.
var HTTPSPort = uint16(443)

// sslDomains is a list of domains to certain SSL certificates for (using minica https://github.com/jsha/minica).
var sslDomains = []string{"localhost", "*." + strings.TrimLeft(project.OptionDomainSuffix.DefaultValue(), ".")}

// GetContainerConfig gets container configuration for the router.
func GetContainerConfig() container.Config {
	// add global domain suffix to list of ssl certifs
	globalConfig, _ := def.ParseGlobalYamlFile()
	if globalConfig != nil && globalConfig.Options[string(project.OptionDomainSuffix)] != "" {
		sslDomains = append(sslDomains, "*."+strings.TrimLeft(globalConfig.Options[string(project.OptionDomainSuffix)], "."))
	}
	routerCmd := fmt.Sprintf(`
mkdir /www
if [ ! -f /var/ssl/minica.pem ]; then
	cd /var/ssl
	~/go/bin/minica -domains "%s"
fi
nginx -g "daemon off;"
`, strings.Join(sslDomains, ","))
	return container.Config{
		ProjectID:  "_",
		ObjectName: "router",
		ObjectType: container.ObjectContainerRouter,
		Command:    []string{"sh", "-c", routerCmd},
		Volumes: map[string]string{
			"pcc_router": "/var/ssl",
		},
		Image:      "docker.io/contextualcode/platform_cc_router:latest",
		WorkingDir: "/routes",
		Ports: []string{
			fmt.Sprintf("%d:80/tcp", HTTPPort),
			fmt.Sprintf("%d:443/tcp", HTTPSPort),
		},
	}
}

func getContainerHandler() (container.Interface, error) {
	// TODO make container handler configurable
	return container.NewDocker()
}

// Start starts the router.
func Start() error {
	done := output.Duration("Start main router.")
	ch, err := getContainerHandler()
	if err != nil {
		return tracerr.Wrap(err)
	}
	// load global config
	gc, err := def.ParseGlobalYamlFile()
	if err != nil {
		return tracerr.Wrap(err)
	}
	HTTPPort = gc.Router.HTTP
	HTTPSPort = gc.Router.HTTPS

	// get container config and pull image
	containerConf := GetContainerConfig()
	if err := ch.ImagePull([]container.Config{containerConf}); err != nil {
		return tracerr.Wrap(err)
	}
	// start container
	if err := ch.ContainerStart(containerConf); err != nil {
		return tracerr.Wrap(err)
	}
	// upload index html
	indexHTMLReader := bytes.NewReader([]byte(routeListHTML))
	if err := ch.ContainerUpload(
		containerConf.GetContainerName(),
		"/www/index.html",
		indexHTMLReader,
	); err != nil {
		return tracerr.Wrap(err)
	}
	// upload nginx conf
	nginxConfReader := bytes.NewReader([]byte(nginxBaseConf))
	if err := ch.ContainerUpload(
		containerConf.GetContainerName(),
		"/etc/nginx/nginx.conf",
		nginxConfReader,
	); err != nil {
		return tracerr.Wrap(err)
	}
	if err := ch.ContainerCommand(
		containerConf.GetContainerName(),
		"root",
		[]string{"sh", "-c", "touch /www/projects.txt && nginx -s reload"},
		nil,
	); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// Stop stops the router.
func Stop() error {
	done := output.Duration("Stop main router.")
	ch, err := getContainerHandler()
	if err != nil {
		return tracerr.Wrap(err)
	}
	if err := ch.ProjectStop("router"); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// Reload issues reload command to nginx in router container.
func Reload() error {
	ch, err := getContainerHandler()
	if err != nil {
		return tracerr.Wrap(err)
	}
	containerConf := GetContainerConfig()
	return tracerr.Wrap(ch.ContainerCommand(
		containerConf.GetContainerName(),
		"root",
		[]string{"nginx", "-s", "reload"},
		nil,
	))
}

// ClearCertificates deletes all certificates files generates by minica.
func ClearCertificates() error {
	ch, err := getContainerHandler()
	if err != nil {
		return tracerr.Wrap(err)
	}
	containerConf := GetContainerConfig()
	return tracerr.Wrap(ch.ContainerCommand(
		containerConf.GetContainerName(),
		"root",
		[]string{"sh", "-c", "rm -rf /var/ssl/*"},
		nil,
	))
}

// DumpCertificateCA returns the CA certificate.
func DumpCertificateCA() ([]byte, error) {
	ch, err := getContainerHandler()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	containerConf := GetContainerConfig()
	var buf bytes.Buffer
	if err := ch.ContainerCommand(
		containerConf.GetContainerName(),
		"root",
		[]string{"sh", "-c", "cat /var/ssl/minica.pem && cat /var/ssl/minica-key.pem"},
		&buf,
	); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return buf.Bytes(), nil
}

// AddProjectRoutes adds given project's routes to router.
func AddProjectRoutes(p *project.Project) error {
	done := output.Duration(
		fmt.Sprintf("Add routes for project '%s.'", p.ID),
	)
	ch, err := getContainerHandler()
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
	if err := ch.ContainerUpload(
		containerConf.GetContainerName(),
		fmt.Sprintf("/routes/%s.conf", p.ID),
		configReader,
	); err != nil {
		return tracerr.Wrap(err)
	}
	if err := Reload(); err != nil {
		return tracerr.Wrap(err)
	}
	// generate and upload route list json
	routeJSON, err := GenerateRouteListJSON(p)
	if err != nil {
		return tracerr.Wrap(err)
	}
	routeJSONReader := bytes.NewReader(routeJSON)
	if err := ch.ContainerUpload(
		containerConf.GetContainerName(),
		fmt.Sprintf("/www/%s.json", p.ID),
		routeJSONReader,
	); err != nil {
		return tracerr.Wrap(err)
	}
	// add to project list + make ssl certifs
	ch.ContainerCommand(
		containerConf.GetContainerName(),
		"root",
		[]string{"sh", "-c", fmt.Sprintf("echo '%s' >> /www/projects.txt", p.ID)},
		nil,
	)
	done()
	return nil
}

// DeleteProjectRoutes deletes routes for given project.
func DeleteProjectRoutes(p *project.Project) error {
	done := output.Duration(
		fmt.Sprintf("Delete routes for project '%s.'", p.ID),
	)
	ch, err := getContainerHandler()
	if err != nil {
		return tracerr.Wrap(err)
	}
	containerConf := GetContainerConfig()
	// delete config file
	if err := ch.ContainerCommand(
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
	if err := Reload(); err != nil {
		return tracerr.Wrap(err)
	}
	// delete json file and remove project id from projects.txt
	if err := ch.ContainerCommand(
		containerConf.GetContainerName(),
		"root",
		[]string{"sh", "-c", fmt.Sprintf("rm -rf /www/%s.json && sed -i \"/%s/d\" /www/projects.txt", p.ID, p.ID)},
		nil,
	); err != nil {
		if !strings.Contains(err.Error(), "No such container") {
			return tracerr.Wrap(err)
		}
		return nil
	}
	done()
	return nil
}
