package api

import (
	"github.com/docker/docker/client"
)

// dockerClient - docker client for platform_cc
type dockerClient struct {
	cli *client.Client
}

// newDockerClient - create docker client
func newDockerClient() (dockerClient, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return dockerClient{}, err
	}
	return dockerClient{
		cli: cli,
	}, nil
}

/*
// StartApp - create and start docker container for given app
func (d *DockerClient) StartApp(p *Project, a *AppDef) error {
	return d.createAppContainer(p, a)
}

// StartService - create and start docker container for given service
func (d *DockerClient) StartService(p *Project, s *ServiceDef) error {
	return nil
}

// StartProject - create and start all docker objects for given project
func (d *DockerClient) StartProject(p *Project) error {
	log.Printf("Start project '%s.'", p.ID)
	// create project network
	if err := d.createProjectNetwork(p); err != nil {
		return err
	}
	// create and start project services
	for _, service := range p.Services {
		if err := d.StartService(p, service); err != nil {
			return err
		}
	}
	// create and start project apps
	for _, app := range p.Apps {
		if err := d.StartApp(p, app); err != nil {
			return err
		}
	}
	return nil
}

// StopProject - stop all docker containers for given project
func (d *DockerClient) StopProject(p *Project) error {
	log.Printf("Stop project '%s.'", p.ID)
	// delete containers
	if err := d.deleteProjectContainers(p); err != nil {
		return err
	}
	// delete project network
	err := d.deleteProjectNetwork(p)
	log.Println(err)
	return err
}

// PurgeProject - purge all docker objects for given project
func (d *DockerClient) PurgeProject(p *Project) error {

	// stop project
	if err := d.StopProject(p); err != nil {
		return err
	}
	// purge
	log.Printf("Purge project '%s.'", p.ID)
	return d.deleteProjectVolumes(p)
}
*/
