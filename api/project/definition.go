package project

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gitlab.com/contextualcode/platform_cc/api/container"
	"gitlab.com/contextualcode/platform_cc/api/def"
)

const containerDataDirectory = "/mnt/data"

// GetDefinitionName returns the name of given definition.
func (p *Project) GetDefinitionName(d interface{}) string {
	switch d.(type) {
	case def.App:
		{
			return d.(def.App).Name
		}
	case def.AppWorker:
		{
			return d.(def.AppWorker).Name
		}
	case def.Service:
		{
			return d.(def.Service).Name
		}
	}
	return ""
}

// GetDefinitionType returns the service type for the given definition.
func (p *Project) GetDefinitionType(d interface{}) string {
	switch d.(type) {
	case def.App:
		{
			return d.(def.App).Type
		}
	case def.AppWorker:
		{
			return d.(def.AppWorker).Type
		}
	case def.Service:
		{
			return d.(def.Service).Type
		}
	}
	return ""
}

// GetDefinitionHostName returns the host name for the container of the given definition.
func (p *Project) GetDefinitionHostName(d interface{}) string {
	dummyConfig := container.Config{
		ProjectID:  p.ID,
		ObjectType: p.GetDefinitionContainerType(d),
		ObjectName: p.GetDefinitionName(d),
	}
	return dummyConfig.GetContainerName()
}

// GetDefinitionStartCommand returns the start command for the given definition.
func (p *Project) GetDefinitionStartCommand(d interface{}) []string {
	switch d.(type) {
	case def.App, def.AppWorker:
		{
			uid, gid := p.getUID()
			command := fmt.Sprintf(
				appContainerCmd,
				uid+1,
				gid+1,
				uid,
				gid,
			)
			return []string{"bash", "--login", "-c", command}
		}
	case def.Service:
		{
			return []string{"bash", "--login", "-c", serviceContainerCmd}
		}
	}
	return []string{}
}

// GetDefinitionContainerType returns the container type for given definition.
func (p *Project) GetDefinitionContainerType(d interface{}) container.ObjectContainerType {
	switch d.(type) {
	case def.App:
		{
			return container.ObjectContainerApp
		}
	case def.AppWorker:
		{
			return container.ObjectContainerWorker
		}
	case def.Service:
		{
			return container.ObjectContainerService
		}
	}
	return container.ObjectContainerNone
}

// GetDefinitionVolumes returns list of Docker volumes for given definition.
func (p *Project) GetDefinitionVolumes(d interface{}) map[string]string {
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
	out := map[string]string{
		name: containerDataDirectory,
	}
	return out
}

// GetDefinitionBinds returns list of Docker binds for given definition.
func (p *Project) GetDefinitionBinds(d interface{}) map[string]string {
	path := ""
	switch d.(type) {
	case def.App:
		{
			path = d.(def.App).Path
			break
		}
	case def.AppWorker:
		{
			path = d.(def.AppWorker).Path
			break
		}
	case def.Service:
		{
			return map[string]string{}
		}
	}
	return map[string]string{
		path: def.AppDir,
	}
}

// GetDefinitionEnvironmentVariables returns list of environment variables for given definition.
func (p *Project) GetDefinitionEnvironmentVariables(d interface{}) map[string]string {
	// get base environment vars
	envVars := p.GetPlatformEnvironmentVariables(d)
	// get variables
	var vars map[string]map[string]interface{}
	switch d.(type) {
	case def.App:
		{
			vars = d.(def.App).Variables
			break
		}
	case def.AppWorker:
		{
			vars = d.(def.AppWorker).Variables
			break
		}
	}
	// append environment variables from .platform.app.yaml
	for k, v := range vars["env"] {
		switch v.(type) {
		case int:
			{
				envVars[k] = fmt.Sprintf("%d", v.(int))
				break
			}
		case string:
			{
				envVars[k] = v.(string)
				break
			}
		}
	}
	// append environment variables from project (var:set command)
	for k, v := range p.Variables["env"] {
		envVars[k] = v
	}
	return envVars
}

// GetDefinitionVariables returns variables for given definition.
func (p *Project) GetDefinitionVariables(d interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	var vars map[string]map[string]interface{}
	switch d.(type) {
	case def.App:
		{
			vars = d.(def.App).Variables
			break
		}
	case def.AppWorker:
		{
			vars = d.(def.AppWorker).Variables
			break
		}
	}
	for varType, varVal := range vars {
		for k, v := range varVal {
			switch v.(type) {
			case int:
				{
					out[fmt.Sprintf("%s:%s", strings.ToLower(varType), k)] = fmt.Sprintf("%d", v.(int))
					break
				}
			case string:
				{
					out[fmt.Sprintf("%s:%s", strings.ToLower(varType), k)] = v.(string)
					break
				}
			}
		}
	}
	for varType, varVal := range p.Variables {
		for k, v := range varVal {
			out[fmt.Sprintf("%s:%s", strings.ToLower(varType), k)] = v
		}
	}
	return out
}

// GetDefinitionRelationships returns relationships for given definition.
func (p *Project) GetDefinitionRelationships(d interface{}) map[string][]map[string]interface{} {
	var rels map[string]string
	switch d.(type) {
	case def.App:
		{
			rels = d.(def.App).Relationships
			break
		}
	case def.AppWorker:
		{
			rels = d.(def.AppWorker).Relationships
			break
		}
	case def.Service:
		{
			rels = d.(def.Service).Relationships
			break
		}
	}
	out := make(map[string][]map[string]interface{})
	for name, rel := range rels {
		out[name] = make([]map[string]interface{}, 0)
		relSplit := strings.Split(rel, ":")
		for _, v := range p.relationships {
			if v["service"].(string) == relSplit[0] && v["rel"] == relSplit[1] {
				out[name] = append(out[name], v)
			}
		}
	}
	return out
}

// GetDefinitionEmptyRelationship returns a relationship template/empty for given definition.
func GetDefinitionEmptyRelationship(d interface{}) map[string]interface{} {
	switch d.(type) {
	case def.App:
		{
			return d.(def.App).GetEmptyRelationship()
		}
	case def.Service:
		{
			return d.(def.Service).GetEmptyRelationship()
		}
	}
	return map[string]interface{}{}
}

// GetDefinitionBuildCommand returns build command for given definition.
func (p *Project) GetDefinitionBuildCommand(d interface{}) string {
	switch d.(type) {
	case def.App:
		{
			// prepare build payload
			uid, gid := p.getUID()
			buildData := map[string]interface{}{
				"application": p.buildConfigAppJSON(d),
				"source_dir":  def.AppDir,
				"output_dir":  def.AppDir,
				"cache_dir":   "/tmp/cache",
				"uid":         uid,
				"gid":         gid,
			}
			buildJSON, _ := json.Marshal(buildData)
			buildB64 := base64.StdEncoding.EncodeToString(buildJSON)
			// build flavor
			buildFlavorComposer := strings.ToLower(d.(def.App).Build.Flavor) == "composer"
			return fmt.Sprintf(
				appBuildCmd,
				strings.Title(strconv.FormatBool(buildFlavorComposer)),
				buildB64,
			)
			//return "sleep 5"
		}
	}
	return ""
}

// GetDefinitionMountCommand returns command to setup mounts for given definition.
func (p *Project) GetDefinitionMountCommand(d interface{}) string {
	var mounts map[string]*def.AppMount
	switch d.(type) {
	case def.App:
		{
			mounts = d.(def.App).Mounts
		}
	case def.AppWorker:
		{
			mounts = d.(def.AppWorker).Mounts
		}
	}
	if mounts != nil {
		out := ""
		for dest, mount := range mounts {
			destPath := fmt.Sprintf(
				"%s/%s",
				def.AppDir,
				strings.Trim(dest, "/"),
			)
			srcPath := strings.TrimRight(fmt.Sprintf(
				"%s/%s",
				containerDataDirectory,
				strings.Trim(mount.SourcePath, "/"),
			), "/")
			destPath = strings.TrimRight(strings.ReplaceAll(
				destPath, ":", "_",
			), "/")
			srcPath = strings.ReplaceAll(
				srcPath, ":", "_",
			)
			if out != "" {
				out += " && "
			}
			out += fmt.Sprintf(
				"mkdir -p %s && mkdir -p %s && chown -f web %s && chown -Rf web %s",
				destPath,
				srcPath,
				destPath,
				srcPath,
			)
			if p.Flags.Has(EnableMountVolume) {
				out += fmt.Sprintf(
					" && mount -o user_xattr --bind %s %s",
					srcPath,
					destPath,
				)
			}
		}
		return out
	}
	return ""
}
