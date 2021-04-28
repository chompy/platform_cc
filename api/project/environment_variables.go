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

package project

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/contextualcode/platform_cc/api/def"
)

const entropySalt = "Dyt+&&*^dKfD9,$rZRA$|I^DLKr%<By"

// GetPlatformRoutes returns the PLATFORM_ROUTES environment variable.
func (p *Project) GetPlatformRoutes() string {
	routes := make(map[string]def.Route)
	for _, route := range p.RoutesReplaceDefault(p.Routes) {
		if strings.HasPrefix(route.Path, ".") {
			continue
		}
		routes[route.Path] = route
	}
	routesJSON, _ := json.Marshal(routes)
	return base64.StdEncoding.EncodeToString(routesJSON)
}

// GetPlatformEntropy returns the PLATFORM_ENTROPY environment variable.
func (p *Project) GetPlatformEntropy() string {
	entH := md5.New()
	entH.Write([]byte(entropySalt))
	entH.Write([]byte(p.ID))
	entH.Write([]byte(entropySalt))
	return fmt.Sprintf("%x", entH.Sum(nil))
}

// GetPlatformApplication returns the PLATFORM_APPLICATION environment variable.
func (p *Project) GetPlatformApplication(d interface{}) string {
	disk := 0
	relationships := map[string]string{}
	mounts := map[string]*def.AppMount{}
	variables := def.Variables{}
	name := ""
	cType := ""
	runtime := def.AppRuntime{}
	web := def.AppWeb{}
	hooks := def.AppHooks{}
	crons := map[string]*def.AppCron{}
	dependencies := def.AppDependencies{}
	switch d := d.(type) {
	case def.App:
		{
			disk = d.Disk
			relationships = d.Relationships
			mounts = d.Mounts
			variables = d.Variables
			name = d.Name
			cType = d.Type
			runtime = d.Runtime
			web = d.Web
			hooks = d.Hooks
			crons = d.Crons
			dependencies = d.Dependencies
			break
		}
	case *def.AppWorker:
		{
			disk = d.Disk
			relationships = d.Relationships
			mounts = d.Mounts
			variables = d.Variables
			name = d.Name
			cType = d.Type
			break
		}
	}
	pfApp := map[string]interface{}{
		"resources":     nil,
		"size":          "AUTO",
		"disk":          disk,
		"access":        map[string]string{},
		"relationships": relationships,
		"mounts":        mounts,
		"timezone":      nil,
		"variables":     variables,
		"firewall":      nil,
		"name":          name,
		"type":          cType,
		"runtime":       runtime,
		"preflight": map[string]interface{}{
			"enabled":       true,
			"ignored_rules": []string{},
		},
		"tree_id":      "-",
		"slug_id":      "-",
		"app_dir":      def.AppDir,
		"web":          web,
		"hook":         hooks,
		"crons":        crons,
		"dependencies": dependencies,
	}
	pfAppJSON, _ := json.Marshal(pfApp)
	return base64.StdEncoding.EncodeToString(pfAppJSON)
}

// GetPlatformVariables returns the PLATFORM_VARIABLES environment variable.
func (p *Project) GetPlatformVariables(d interface{}) string {
	varJSON, _ := json.Marshal(p.GetDefinitionVariables(d))
	return base64.StdEncoding.EncodeToString(varJSON)
}

// GetPlatformRelationships returns the PLATFORM_RELATIONSHIPS environment variable.
func (p *Project) GetPlatformRelationships(d interface{}) string {
	rels := p.GetDefinitionRelationships(d)
	relsJSON, _ := json.Marshal(rels)
	return base64.StdEncoding.EncodeToString(relsJSON)
}

// GetPlatformEnvironmentVariables returns a map of all PLATFORM_ environment variables.
func (p *Project) GetPlatformEnvironmentVariables(d interface{}) map[string]string {
	name := ""
	switch d := d.(type) {
	case def.App:
		{
			name = d.Name
			break
		}
	case *def.AppWorker:
		{
			name = d.Name
			break
		}
	case def.Service:
		{
			name = d.Name
			break
		}
	}
	return map[string]string{
		"PLATFORM_DOCUMENT_ROOT":    def.AppDir + "/web",
		"PLATFORM_APPLICATION":      p.GetPlatformApplication(d),
		"PLATFORM_PROJECT":          p.ID,
		"PLATFORM_PROJECT_ENTROPY":  p.GetPlatformEntropy(),
		"PLATFORM_APPLICATION_NAME": name,
		"PLATFORM_BRANCH":           "pcc",
		"PLATFORM_DIR":              def.AppDir,
		"PLATFORM_APP_DIR":          def.AppDir,
		"PLATFORM_TREE_ID":          "-",
		"PLATFORM_ENVIRONMENT":      "pcc",
		"PLATFORM_VARIABLES":        p.GetPlatformVariables(d),
		"PLATFORM_RELATIONSHIPS":    p.GetPlatformRelationships(d),
		"PLATFORM_ROUTES":           p.GetPlatformRoutes(),
		"PLATFORM_CACHE_DIR":        "/tmp/cache",
	}
}
