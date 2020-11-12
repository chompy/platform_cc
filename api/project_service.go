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
func (p *Project) getServiceContainerConfig(service *ServiceDef) dockerContainerConfig {
	typeName := strings.Split(service.Type, ":")
	return dockerContainerConfig{
		projectID:  p.ID,
		objectName: service.Name,
		objectType: objectContainerService,
		command:    []string{"sh", "-c", serviceContainerCmd},
		Image:      fmt.Sprintf("%s%s-%s", platformShDockerImagePrefix, typeName[0], typeName[1]),
		Volumes: map[string]string{
			fmt.Sprintf(dockerNamingPrefix+"%s-v", p.ID, service.Name): "/mnt/data",
		},
	}
}

// startService - start an service
func (p *Project) startService(def *ServiceDef) error {
	log.Printf("Start service '%s.'", def.Name)
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
func (p *Project) openService(def *ServiceDef) error {
	log.Printf("Open service '%s.'", def.Name)
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
	// TODO (assumed this is need by varnish)
	rel := map[string]interface{}{
		"relationships": map[string]interface{}{},
	}
	relJSON, err := json.Marshal(rel)
	if err != nil {
		return err
	}
	relB64 := base64.StdEncoding.EncodeToString(relJSON)
	// open service and retrieve relationships
	var openOutput bytes.Buffer
	cmd := fmt.Sprintf(serviceOpenCmd, relB64)
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
		rel := def.GetEmptyRelationship()
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

// ShellService - shell in to service
func (p *Project) ShellService(def *ServiceDef, command []string) error {
	log.Printf("Shell in to service '%s.'", def.Name)
	// get container config
	containerConfig := p.getServiceContainerConfig(def)
	return p.docker.ShellContainer(
		containerConfig.GetContainerName(),
		"root",
		command,
	)
}
