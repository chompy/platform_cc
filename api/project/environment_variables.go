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
	for _, route := range p.Routes {
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
	variables := map[string]map[string]interface{}{}
	name := ""
	cType := ""
	runtime := def.AppRuntime{}
	web := def.AppWeb{}
	hooks := def.AppHooks{}
	crons := map[string]*def.AppCron{}
	dependencies := def.AppDependencies{}
	switch d.(type) {
	case def.App:
		{
			defApp := d.(def.App)
			disk = defApp.Disk
			relationships = defApp.Relationships
			mounts = defApp.Mounts
			variables = defApp.Variables
			name = defApp.Name
			cType = defApp.Type
			runtime = defApp.Runtime
			web = defApp.Web
			hooks = defApp.Hooks
			crons = defApp.Crons
			dependencies = defApp.Dependencies
			break
		}
	case def.AppWorker:
		{
			defAppWorker := d.(def.AppWorker)
			disk = defAppWorker.Disk
			relationships = defAppWorker.Relationships
			mounts = defAppWorker.Mounts
			variables = defAppWorker.Variables
			name = defAppWorker.Name
			cType = defAppWorker.Type
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
	switch d.(type) {
	case def.App:
		{
			name = d.(def.App).Name
			break
		}
	case def.AppWorker:
		{
			name = d.(def.AppWorker).Name
			break
		}
	case def.Service:
		{
			name = d.(def.Service).Name
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
