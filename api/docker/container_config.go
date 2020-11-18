package docker

import "fmt"

const dockerNamingPrefix = "pcc-%s-"

// containerNamePrefix is the container name prefix.
const containerNameFormat = dockerNamingPrefix + "%s-%s"

// containerNetworkNameFormat is the container network name format.
const containerNetworkNameFormat = dockerNamingPrefix + "n"

// ObjectContainerType defines the type of container.
type ObjectContainerType byte

const (
	// ObjectContainerApp is an application container.
	ObjectContainerApp ObjectContainerType = 'a'
	// ObjectContainerService is a service container.
	ObjectContainerService ObjectContainerType = 's'
	// ObjectContainerRouter is the router container.
	ObjectContainerRouter ObjectContainerType = 'r'
)

// TypeName gets the type of container as a string.
func (o ObjectContainerType) TypeName() string {
	switch o {
	case ObjectContainerApp:
		{
			return "app"
		}
	case ObjectContainerService:
		{
			return "service"
		}
	case ObjectContainerRouter:
		{
			return "router"
		}
	}
	return "unknown"
}

// ContainerConfig contains configuration for a Docker container.
type ContainerConfig struct {
	ProjectID  string
	ObjectType ObjectContainerType
	ObjectName string
	Command    []string
	Image      string
	Volumes    map[string]string
	Binds      map[string]string
	Env        map[string]string
	Ports      []string
	WorkingDir string
}

// GetContainerName gets the name of the Docker container.
func (d ContainerConfig) GetContainerName() string {
	if d.ObjectType == ObjectContainerRouter {
		return "pcc-router-1"
	}
	return fmt.Sprintf(containerNameFormat, d.ProjectID, string(d.ObjectType), d.ObjectName)
}

// GetNetworkName gets the name of the Docker network.
func (d ContainerConfig) GetNetworkName() string {
	return fmt.Sprintf(containerNetworkNameFormat, d.ProjectID)
}

// GetCommand gets the container command.
func (d ContainerConfig) GetCommand() []string {
	if len(d.Command) > 0 {
		return d.Command
	}
	return []string{"tail", "-f", "/dev/null"}
}

// GetEnv converts environment vars to format needed to start docker container.
func (d ContainerConfig) GetEnv() []string {
	out := make([]string, 0)
	for k, v := range d.Env {
		out = append(out, fmt.Sprintf("%s=%v", k, v))
	}
	return out
}
