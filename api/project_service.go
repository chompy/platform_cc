package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"gitlab.com/contextualcode/platform_cc/def"
)

// GetServiceContainerConfig gets container configuration for a service.
func (p *Project) GetServiceContainerConfig(d interface{}) DockerContainerConfig {
	var name string
	var typeName []string
	var objectType objectContainerType
	var command string
	var env map[string]string
	var binds map[string]string
	var workingDir string
	switch d.(type) {
	case *def.App:
		{
			name = d.(*def.App).Name
			typeName = strings.Split(d.(*def.App).Type, ":")
			objectType = objectContainerApp
			uid, gid := p.getUID()
			command = fmt.Sprintf(
				appContainerCmd,
				uid+1,
				gid+1,
				uid,
				gid,
			)
			env = p.getAppEnvironmentVariables(d.(*def.App))
			binds = map[string]string{
				d.(*def.App).Path: "/app",
			}
			workingDir = "/app"
			break
		}
	case *def.Service:
		{
			name = d.(*def.Service).Name
			typeName = strings.Split(d.(*def.Service).Type, ":")
			objectType = objectContainerService
			command = serviceContainerCmd
			env = map[string]string{}
			binds = map[string]string{}
			workingDir = "/"
			break
		}
	default:
		{
			return DockerContainerConfig{}
		}
	}
	return DockerContainerConfig{
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
func (p *Project) startService(d interface{}) error {
	var name string
	switch d.(type) {
	case *def.App:
		{
			name = d.(*def.App).Name
			log.Printf("Start application '%s.'", name)
			break
		}
	case *def.Service:
		{
			name = d.(*def.Service).Name
			log.Printf("Start service '%s.'", name)
			break
		}
	default:
		{
			return fmt.Errorf("passed dinition is not an application or service")
		}
	}
	// get container config
	containerConfig := p.GetServiceContainerConfig(d)
	// build config.json
	configJSON, err := p.BuildConfigJSON(d)
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
func (p *Project) openService(d interface{}) error {
	var name string
	switch d.(type) {
	case *def.App:
		{
			name = d.(*def.App).Name
			log.Printf("Open application '%s.'", name)
			break
		}
	case *def.Service:
		{
			name = d.(*def.Service).Name
			log.Printf("Open service '%s.'", name)
			break
		}
	default:
		{
			return fmt.Errorf("passed dinition is not an application or service")
		}
	}
	// get container config
	containerConfig := p.GetServiceContainerConfig(d)
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
	relationshipsVar, err := p.getServiceRelationships(d)
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
	ipAddress, err := p.docker.GetContainerIP(p.GetServiceContainerConfig(d).GetContainerName())
	if err != nil {
		return err
	}
	for k, v := range data {
		var rel map[string]interface{}
		switch d.(type) {
		case *def.App:
			{
				rel = d.(*def.App).GetEmptyRelationship()
				break
			}
		case *def.Service:
			{
				rel = d.(*def.Service).GetEmptyRelationship()
				break
			}
		}
		for kk, vv := range v.(map[string]interface{}) {
			rel[kk] = vv
		}
		rel["rel"] = k
		rel["host"] = containerConfig.GetContainerName()
		rel["hostname"] = containerConfig.GetContainerName()
		rel["ip"] = ipAddress
		p.relationships = append(p.relationships, rel)
	}
	return nil
}

// getServiceRelationships - generate relationships variable for service
func (p *Project) getServiceRelationships(d interface{}) (map[string][]map[string]interface{}, error) {
	var relmap map[string]string
	switch d.(type) {
	case *def.App:
		{
			relmap = d.(*def.App).Relationships
			break
		}
	case *def.Service:
		{
			relmap = d.(*def.Service).Relationships
			break
		}
	default:
		{
			return nil, fmt.Errorf("passed dinition is not an application or service")
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
func (p *Project) ShellService(d interface{}, command []string) error {
	var name string
	var user string
	switch d.(type) {
	case *def.App:
		{
			name = d.(*def.App).Name
			user = "web"
			log.Printf("Shell in to application '%s.'", name)
			break
		}
	case *def.Service:
		{
			name = d.(*def.Service).Name
			user = "root"
			log.Printf("Shell in to service '%s.'", name)
			break
		}
	default:
		{
			return fmt.Errorf("passed dinition is not for an application or service")
		}
	}
	// get container config
	containerConfig := p.GetServiceContainerConfig(d)
	return p.docker.ShellContainer(
		containerConfig.GetContainerName(),
		user,
		command,
	)
}
