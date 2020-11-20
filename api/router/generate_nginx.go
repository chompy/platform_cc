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
	"net/url"
	"strings"
	"text/template"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

// GetUpstreamHost retrieves upstream hostname from upstream value in route.
func GetUpstreamHost(proj *project.Project, upstream string, allowServices bool) string {
	upstreamSplit := strings.Split(upstream, ":")
	// itterate apps and services to find name match
	// TODO this should use relationships but those only get resolved when
	// services are opened...sooo??
	for _, app := range proj.Apps {
		if app.Name == upstreamSplit[0] {
			containerConfig := proj.GetAppContainerConfig(app)
			return containerConfig.GetContainerName()
		}
	}
	for _, serv := range proj.Services {
		if serv.Name == upstreamSplit[0] {
			// forward to app if allowServices is false
			if !allowServices {
				for _, relationship := range serv.Relationships {
					rlSplit := strings.Split(relationship, ":")
					return GetUpstreamHost(proj, fmt.Sprintf("%s:http", rlSplit[0]), allowServices)
				}
			}
			containerConfig := proj.GetServiceContainerConfig(serv)
			// TODO use relationship to determine port
			port := 80
			switch serv.GetTypeName() {
			case "varnish":
				{
					port = 8080
					break
				}
			}
			return fmt.Sprintf("%s:%d", containerConfig.GetContainerName(), port)
		}
	}
	return ""
}

// GenerateTemplateVars generates variables to inject in nginx template.
func GenerateTemplateVars(proj *project.Project) []map[string]interface{} {
	hostMaps := MapHostRoutes(proj.Routes)
	out := make([]map[string]interface{}, 0)
	for _, hostMap := range hostMaps {
		outHm := map[string]interface{}{
			"host":   hostMap.Host,
			"routes": make([]map[string]interface{}, 0),
		}
		hasDefaultPath := false
		for _, route := range hostMap.Routes {
			urlParse, _ := url.Parse(route.Path)
			path := strings.TrimSpace(urlParse.Path)
			if path == "" {
				path = "/"
			}
			if path == "/" {
				hasDefaultPath = true
			}
			outHm["routes"] = append(
				outHm["routes"].([]map[string]interface{}),
				map[string]interface{}{
					"path": path,
					"type": route.Type,
					"upstream": GetUpstreamHost(
						proj, route.Upstream, proj.Flags.Has(project.EnableServiceRoutes),
					),
					"to":    route.To,
					"route": route,
				},
			)
		}
		if !hasDefaultPath {
			log.Println(hostMap.Host)
			outHm["routes"] = append(
				outHm["routes"].([]map[string]interface{}),
				map[string]interface{}{
					"path":     "/",
					"type":     "redirect",
					"upstream": "",
					"to":       outHm["routes"].([]map[string]interface{})[0]["path"],
					"route":    nil,
				},
			)
		}
		out = append(out, outHm)
	}
	return out
}

// GenerateNginxConfig creates nginx configuration for given application.
func GenerateNginxConfig(proj *project.Project) ([]byte, error) {
	t, err := template.New("nginx.conf").Parse(nginxServerTemplate)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, GenerateTemplateVars(proj)); err != nil {
		return nil, tracerr.Wrap(err)
	}
	return buf.Bytes(), nil
}
