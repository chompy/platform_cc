package api

import (
	"fmt"
	"log"
)

// getServiceContainerConfig - get container configuration for service
func (p *Project) getServiceContainerConfig(service *ServiceDef) dockerContainerConfig {
	serviceContainerConfig := service.GetContainerConfig()
	if serviceContainerConfig == nil {
		return dockerContainerConfig{}
	}
	volumes := make(map[string]string)
	for i, dest := range serviceContainerConfig.GetVolumes() {
		volumes[fmt.Sprintf(dockerNamingPrefix+"%s-v-%d", p.ID, service.Name, i)] = dest
	}
	return dockerContainerConfig{
		projectID:  p.ID,
		objectName: service.Name,
		objectType: objectContainerService,
		command:    serviceContainerConfig.GetStartCommand(),
		Image:      serviceContainerConfig.GetImage(),
		Volumes:    volumes,
	}
}

// startService - start an service
func (p *Project) startService(service *ServiceDef) error {
	log.Printf("Start service '%s.'", service.Name)
	// get container config
	containerConfig := p.getServiceContainerConfig(service)
	serviceContainerConfig := service.GetContainerConfig()
	if serviceContainerConfig == nil {
		return fmt.Errorf("service container config not found for service of type %s", service.Type)
	}
	// start container
	if err := p.docker.StartContainer(containerConfig); err != nil {
		return err
	}
	// post start command
	log.Print("Post start commands for service")
	if err := p.docker.RunContainerCommand(
		containerConfig.GetContainerName(),
		"root",
		serviceContainerConfig.GetPostStartCommand(),
		nil,
	); err != nil {
		return err
	}
	return nil
}

// getServiceRelationships - get service relationship value
func (p *Project) getServiceRelationships(service *ServiceDef) ([]map[string]interface{}, error) {
	serviceContainerConfig := service.GetContainerConfig()
	if serviceContainerConfig == nil {
		return nil, fmt.Errorf("service container config not found for service of type %s", service.Type)
	}
	out := service.GetContainerConfig().GetRelationship()
	// get container ip address
	ipAddress, err := p.docker.GetContainerIP(p.getServiceContainerConfig(service).GetContainerName())
	if err != nil {
		return nil, err
	}
	// add host/ip
	for i := range out {
		out[i]["host"] = ipAddress
		out[i]["hostname"] = ipAddress
		out[i]["ip"] = ipAddress
	}
	return out, nil
}
