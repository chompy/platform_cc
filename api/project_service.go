package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// getServiceContainerConfig - get container configuration for service
func (p *Project) getServiceContainerConfig(def interface{}) dockerContainerConfig {
	var name string
	var typeName []string
	var objectType objectContainerType
	var command string
	var env map[string]string
	var binds map[string]string
	var workingDir string
	switch def.(type) {
	case *AppDef:
		{
			name = def.(*AppDef).Name
			typeName = strings.Split(def.(*AppDef).Type, ":")
			objectType = objectContainerApp
			uid, gid := p.getUID()
			command = fmt.Sprintf(
				appContainerCmd,
				uid+1,
				gid+1,
				uid,
				gid,
			)
			env = p.getAppEnvironmentVariables(def.(*AppDef))
			binds = map[string]string{
				def.(*AppDef).Path: "/app",
			}
			workingDir = "/app"
			break
		}
	case *ServiceDef:
		{
			name = def.(*ServiceDef).Name
			typeName = strings.Split(def.(*ServiceDef).Type, ":")
			objectType = objectContainerService
			command = serviceContainerCmd
			env = map[string]string{}
			binds = map[string]string{}
			workingDir = "/"
			break
		}
	default:
		{
			return dockerContainerConfig{}
		}
	}
	return dockerContainerConfig{
		projectID:  p.ID,
		objectName: name,
		objectType: objectType,
		command:    []string{"sh", "-c", command},
		Image:      fmt.Sprintf("%s%s-%s", platformShDockerImagePrefix, typeName[0], typeName[1]),
		Volumes: map[string]string{
			fmt.Sprintf(dockerNamingPrefix+"%s-v", p.ID, name): "/mnt/data",
		},
		Binds:      binds,
		Env:        env,
		WorkingDir: workingDir,
	}
}

// startService - start an service
func (p *Project) startService(def interface{}) error {
	var name string
	switch def.(type) {
	case *AppDef:
		{
			name = def.(*AppDef).Name
			log.Printf("Start application '%s.'", name)
			break
		}
	case *ServiceDef:
		{
			name = def.(*ServiceDef).Name
			log.Printf("Start service '%s.'", name)
			break
		}
	default:
		{
			return fmt.Errorf("passed definition is not an application or service")
		}
	}
	// get container config
	containerConfig := p.getServiceContainerConfig(def)
	// build config.json
	configJSON, err := p.BuildConfigJSON(def)
	if err != nil {
		return err
	}
	configJSONReader := bytes.NewReader(configJSON)
	// start container
	if err := p.docker.StartContainer(containerConfig); err != nil {
		return err
	}
	// upload config.json
	if err := p.docker.UploadDataToContainer(
		containerConfig.GetContainerName(),
		"/run/config.json",
		configJSONReader,
	); err != nil {
		return err
	}
	return nil
}

// openService - open service and get relationships
func (p *Project) openService(def interface{}) error {
	var name string
	switch def.(type) {
	case *AppDef:
		{
			name = def.(*AppDef).Name
			log.Printf("Open application '%s.'", name)
			break
		}
	case *ServiceDef:
		{
			name = def.(*ServiceDef).Name
			log.Printf("Open service '%s.'", name)
			break
		}
	default:
		{
			return fmt.Errorf("passed definition is not an application or service")
		}
	}
	// get container config
	containerConfig := p.getServiceContainerConfig(def)
	// start service
	if err := p.docker.RunContainerCommand(
		containerConfig.GetContainerName(),
		"root",
		[]string{"sh", "-c", serviceStartCmd},
		os.Stdout,
	); err != nil {
		return err
	}
	// prepare input relationships
	relationshipsVar, err := p.getServiceRelationships(def)
	if err != nil {
		return err
	}
	relationships := map[string]interface{}{
		"relationships": relationshipsVar,
	}
	relationshipsJSON, err := json.Marshal(relationships)
	if err != nil {
		return err
	}
	relationshipsB64 := base64.StdEncoding.EncodeToString(relationshipsJSON)
	// open service and retrieve relationships
	var openOutput bytes.Buffer
	cmd := fmt.Sprintf(serviceOpenCmd, relationshipsB64)
	if err := p.docker.RunContainerCommand(
		containerConfig.GetContainerName(),
		"root",
		[]string{"sh", "-c", cmd},
		&openOutput,
	); err != nil {
		return err
	}
	// process output relationships
	openOutlineLines := bytes.Split(openOutput.Bytes(), []byte{'\n'})
	rlRaw := openOutlineLines[len(openOutlineLines)-1]
	data := make(map[string]interface{})
	if err := json.Unmarshal(rlRaw, &data); err != nil {
		return err
	}
	ipAddress, err := p.docker.GetContainerIP(p.getServiceContainerConfig(def).GetContainerName())
	if err != nil {
		return err
	}
	for k, v := range data {
		var rel map[string]interface{}
		switch def.(type) {
		case *AppDef:
			{
				rel = def.(*AppDef).GetEmptyRelationship()
				break
			}
		case *ServiceDef:
			{
				rel = def.(*ServiceDef).GetEmptyRelationship()
				break
			}
		}
		for kk, vv := range v.(map[string]interface{}) {
			rel[kk] = vv
		}
		rel["rel"] = k
		rel["host"] = ipAddress
		rel["hostname"] = ipAddress
		rel["ip"] = ipAddress
		p.relationships = append(p.relationships, rel)
	}
	return nil
}

// getServiceRelationships - generate relationships variable for service
func (p *Project) getServiceRelationships(def interface{}) (map[string][]map[string]interface{}, error) {
	var relmap map[string]string
	switch def.(type) {
	case *AppDef:
		{
			relmap = def.(*AppDef).Relationships
			break
		}
	case *ServiceDef:
		{
			relmap = def.(*ServiceDef).Relationships
			break
		}
	default:
		{
			return nil, fmt.Errorf("passed definition is not an application or service")
		}
	}
	out := make(map[string][]map[string]interface{})
	for name, rel := range relmap {
		out[name] = make([]map[string]interface{}, 0)
		relSplit := strings.Split(rel, ":")
		for _, v := range p.relationships {
			if v["service"].(string) == relSplit[0] && v["rel"] == relSplit[1] {
				out[name] = append(out[name], v)
			}
		}
	}
	return out, nil
}

// ShellService - shell in to service
func (p *Project) ShellService(def interface{}, command []string) error {
	var name string
	var user string
	switch def.(type) {
	case *AppDef:
		{
			name = def.(*AppDef).Name
			user = "web"
			log.Printf("Shell in to application '%s.'", name)
			break
		}
	case *ServiceDef:
		{
			name = def.(*ServiceDef).Name
			user = "root"
			log.Printf("Shell in to service '%s.'", name)
			break
		}
	default:
		{
			return fmt.Errorf("passed definition is not for an application or service")
		}
	}
	// get container config
	containerConfig := p.getServiceContainerConfig(def)
	return p.docker.ShellContainer(
		containerConfig.GetContainerName(),
		user,
		command,
	)
}
