package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/project"
)

// ReplaceDefault replaces {default} in given path with project id + domain suffix.
func ReplaceDefault(path string, p *project.Project) string {
	out := strings.ReplaceAll(
		path,
		"{default}",
		fmt.Sprintf(
			"%s.default",
			p.ID,
		),
	)
	out = strings.ReplaceAll(
		out,
		"__PID__",
		p.ID,
	)
	return out
}

// RouteReplaceDefault replaces {default} values in all paths.
func RouteReplaceDefault(r def.Route, p *project.Project) def.Route {
	r.OriginalURL = ReplaceDefault(r.OriginalURL, p)
	r.Path = ReplaceDefault(r.Path, p)
	r.To = ReplaceDefault(r.To, p)
	for i := range r.Redirects.Paths {
		r.Redirects.Paths[i].To = ReplaceDefault(r.Redirects.Paths[i].To, p)
	}
	for k, v := range r.Attributes {
		r.Attributes[k] = ReplaceDefault(v, p)
	}
	return r
}

// RoutesReplaceDefault replaces {default} values in all paths for multiple routes.
func RoutesReplaceDefault(routes []def.Route, p *project.Project) []def.Route {
	for i := range routes {
		routes[i] = RouteReplaceDefault(routes[i], p)
	}
	return routes
}

// ListActiveProjects returns list of project ids of projects current loaded in to the router.
func ListActiveProjects() ([]string, error) {
	ch, err := getContainerHandler()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	containerConf := GetContainerConfig()
	var buf bytes.Buffer
	if err := ch.ContainerCommand(
		containerConf.GetContainerName(),
		"root",
		[]string{"cat", "/www/projects.txt"},
		&buf,
	); err != nil {
		return nil, tracerr.Wrap(err)
	}
	projectIDs := make([]string, 0)
	for {
		line, err := buf.ReadString('\n')
		if line != "" {
			projectIDs = append(projectIDs, strings.TrimSpace(line))
		}
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				break
			}
			return nil, tracerr.Wrap(err)
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
		return []def.Route{}, tracerr.Wrap(err)
	}
	// itterate project ids and get route data from container
	ch, err := getContainerHandler()
	if err != nil {
		return []def.Route{}, tracerr.Wrap(err)
	}
	containerConf := GetContainerConfig()
	out := make([]def.Route, 0)
	for _, pid := range projectIDs {
		var buf bytes.Buffer
		if err := ch.ContainerCommand(
			containerConf.GetContainerName(),
			"root",
			[]string{"cat", fmt.Sprintf("/www/%s.json", pid)},
			&buf,
		); err != nil {
			return nil, tracerr.Wrap(err)
		}
		data := make([]activeRouterData, 0)
		if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
			return []def.Route{}, nil
		}
		for _, ar := range data {
			for _, route := range ar.Routes {
				route.Route.Path = fmt.Sprintf("http://%s%s", ar.Host, route.Path)
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
