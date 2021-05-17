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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/api/output"

	"gitlab.com/contextualcode/platform_cc/api/def"
)

// ListActiveProjects returns list of project ids of projects current loaded in to the router.
func ListActiveProjects() ([]string, error) {
	ch, err := getContainerHandler()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	containerConf := GetContainerConfig()
	var buf bytes.Buffer
	if _, err := ch.ContainerCommand(
		containerConf.GetContainerName(),
		"root",
		[]string{"cat", "/www/projects.txt"},
		&buf,
	); err != nil {
		return nil, errors.WithStack(err)
	}
	projectIDs := make([]string, 0)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			hasProjectID := false
			pidToInsert := strings.TrimSpace(line)
			for _, pid := range projectIDs {
				if pid == pidToInsert {
					hasProjectID = true
					break
				}
			}
			if !hasProjectID {
				projectIDs = append(projectIDs, pidToInsert)
			}
		}
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				break
			}
			return nil, errors.WithStack(err)
		}
	}
	return projectIDs, nil
}

type activeRouterData struct {
	Host   string `json:"host"`
	Routes []struct {
		Path     string    `json:"path"`
		Type     string    `json:"type"`
		Upstream string    `json:"upstream"`
		To       string    `json:"to"`
		Route    def.Route `json:"route"`
	} `json:"routes"`
}

// ListActiveRoutes returns list of routes currently active in the router.
func ListActiveRoutes() ([]def.Route, error) {
	// get project ids
	projectIDs, err := ListActiveProjects()
	if err != nil {
		return []def.Route{}, errors.WithStack(err)
	}
	// itterate project ids and get route data from container
	ch, err := getContainerHandler()
	if err != nil {
		return []def.Route{}, errors.WithStack(err)
	}
	containerConf := GetContainerConfig()
	out := make([]def.Route, 0)
	for _, pid := range projectIDs {
		var buf bytes.Buffer
		if _, err := ch.ContainerCommand(
			containerConf.GetContainerName(),
			"root",
			[]string{"cat", fmt.Sprintf("/www/%s.json", pid)},
			&buf,
		); err != nil {
			return nil, errors.WithStack(err)
		}
		data := make([]activeRouterData, 0)
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			output.LogError(err)
			return []def.Route{}, nil
		}
		for _, ar := range data {
			for _, route := range ar.Routes {
				route.Route.Path = fmt.Sprintf("https://%s%s", ar.Host, route.Path)
				route.Route.Attributes = map[string]string{
					"host":          ar.Host,
					"upstream_host": route.Upstream,
					"project_id":    pid,
				}
				out = append(out, route.Route)
			}
		}
	}
	return out, nil
}
