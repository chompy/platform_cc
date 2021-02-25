package project

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/contextualcode/platform_cc/api/container"
	"gitlab.com/contextualcode/platform_cc/api/def"
)

const containerMntPath = "/mnt"
const symlinkMntPath = def.AppDir + "/.platform_cc_mnt"

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
		name: containerMntPath,
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
	var vars def.Variables
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
	for k, v := range vars.GetStringSubMap("env") {
		envVars[k] = v
	}
	// append environment variables from project (var:set command)
	for k, v := range p.Variables.GetStringSubMap("env") {
		envVars[k] = v
	}
	return envVars
}

// GetDefinitionVariables returns flattened variables for given definition.
func (p *Project) GetDefinitionVariables(d interface{}) map[string]interface{} {
	out := make(def.Variables)
	out.Merge(p.globalConfig.Variables)
	out.Merge(p.Variables)
	switch d.(type) {
	case def.App:
		{
			out.Merge(d.(def.App).Variables)
			break
		}
	case def.AppWorker:
		{
			out.Merge(d.(def.AppWorker).Variables)
			break
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
				"source_dir":  def.AppDir + "/",
				"output_dir":  def.AppDir + "/",
				"cache_dir":   "/tmp/cache",
				"uid":         uid,
				"gid":         gid,
			}
			buildJSON, _ := json.Marshal(buildData)
			buildB64 := base64.StdEncoding.EncodeToString(buildJSON)
			return fmt.Sprintf(
				appBuildCmd,
				buildB64,
			)
		}
	}
	return ""
}

// GetDefinitionPostDeployCommand returns post deploy command for given definition.
func (p *Project) GetDefinitionPostDeployCommand(d interface{}) string {
	switch d.(type) {
	case def.App:
		{
			return d.(def.App).Hooks.PostDeploy
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
		out := make([]string, 0)
		for dest, mount := range mounts {
			// build path to destination directory inside app root
			destPath := fmt.Sprintf(
				"%s/%s",
				def.AppDir,
				strings.Trim(dest, "/"),
			)
			destPath = strings.TrimRight(strings.ReplaceAll(
				destPath, ":", "_",
			), "/")
			// handle mount depending on user selected strategy
			switch p.GetOption(OptionMountStrategy) {
			case MountStrategyNone:
				{
					// create dest directory if it doesn't exist, fix persmission
					// the dest directory will not be mounted to anything, it'll just be part
					// of the main /app directory which will be mounted to the user's host file system
					out = append(out, fmt.Sprintf(
						"mkdir -p %s && chown -Rf web %s",
						destPath,
						destPath,
					))
					break
				}
			case MountStrategySymlink:
				{
					// build source path inside app root under .platform_cc_mnt
					srcPath := strings.TrimRight(fmt.Sprintf(
						"%s/%s",
						symlinkMntPath,
						strings.Trim(mount.SourcePath, "/"),
					), "/")
					srcPath = strings.ReplaceAll(
						srcPath, ":", "_",
					)
					// use symlink to link everything to main mount dir
					// this will allow the use of mount subdirectories while allowing
					// mount files to be accessable outside of the container
					out = append(out, fmt.Sprintf(
						"mkdir -p %s && chown -Rf web %s && rm -rf %s && ln -s -r %s %s",
						srcPath,
						srcPath,
						destPath,
						srcPath,
						destPath,
					))
					break
				}
			case MountStrategyVolume:
				{
					// build source path to mounted container volume
					srcPath := strings.TrimRight(fmt.Sprintf(
						"%s/%s",
						containerMntPath,
						strings.Trim(mount.SourcePath, "/"),
					), "/")
					srcPath = strings.ReplaceAll(
						srcPath, ":", "_",
					)
					// perform a bind mount between container volume and destination directory
					out = append(out, fmt.Sprintf(
						"mkdir -p %s && chown -Rf web %s && mkdir -p %s && chown -Rf web %s && mount -o user_xattr --bind %s %s",
						srcPath,
						srcPath,
						destPath,
						destPath,
						srcPath,
						destPath,
					))
					break
				}
			}
		}
		return strings.Join(out, " && ")
	}
	return ""
}

// GetDefinitionStartOrder given list of definitions reorder them for optimal start order for relationships.
func (p *Project) GetDefinitionStartOrder(defs []interface{}) ([]interface{}, error) {
	getRelNames := func(d interface{}) []string {
		rels := map[string]string{}
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
		out := make([]string, 0)
		for _, rel := range rels {
			relSplit := strings.Split(rel, ":")
			out = append(out, relSplit[0])
		}
		return out
	}
	nameList := make([]string, 0)
	out := make([]interface{}, 0)
	hasDef := func(def interface{}) bool {
		name := p.GetDefinitionName(def)
		for _, outDef := range out {
			if p.GetDefinitionName(outDef) == name {
				return true
			}
		}
		return false
	}
	hasAllRels := func(rels []string) bool {
		for _, rel := range rels {
			has := false
			for _, name := range nameList {
				if rel == name {
					has = true
					break
				}
			}
			if !has {
				return false
			}
		}
		return true
	}
	ittCount := 0
	for len(out) < len(defs) {
		ittCount++
		for _, def := range defs {
			if hasDef(def) {
				continue
			}
			name := p.GetDefinitionName(def)
			rels := getRelNames(def)
			if len(rels) == 0 || hasAllRels(rels) {
				out = append(out, def)
				nameList = append(nameList, name)
			}
		}
		// assume circular relationship if more than 100 itterations are needed
		if ittCount > 100 {
			invalidOut := make([]string, 0)
			for _, def := range defs {
				if hasDef(def) {
					continue
				}
				invalidOut = append(invalidOut, p.GetDefinitionName(def))
			}
			return nil, fmt.Errorf("one or more relationships are invalid: %s", strings.Join(invalidOut, ","))
		}

	}
	return out, nil
}
