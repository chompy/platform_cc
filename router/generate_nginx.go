package router

import (
	"bytes"
	"log"
	"net/url"
	"strings"
	"text/template"

	"gitlab.com/contextualcode/platform_cc/api"
)

// GetUpstreamHost retrieves upstream hostname from upstream value in route.
func GetUpstreamHost(proj *api.Project, upstream string) string {
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
			containerConfig := proj.GetServiceContainerConfig(serv)
			return containerConfig.GetContainerName()
		}
	}
	return ""
}

// GenerateTemplateVars generates variables to inject in nginx template.
func GenerateTemplateVars(proj *api.Project) []map[string]interface{} {
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
			if urlParse.Path == "/" {
				hasDefaultPath = true
			}
			outHm["routes"] = append(
				outHm["routes"].([]map[string]interface{}),
				map[string]interface{}{
					"path":     urlParse.Path,
					"type":     route.Type,
					"upstream": GetUpstreamHost(proj, route.Upstream),
					"to":       route.To,
					"route":    route,
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
func GenerateNginxConfig(proj *api.Project) ([]byte, error) {

	t, err := template.New("nginx.conf").Parse(nginxTemplate)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, GenerateTemplateVars(proj)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
