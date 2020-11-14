package api

import "fmt"

const dockerNamingPrefix = "pcc-%s-"

// containerNamePrefix - container name prefix
const containerNameFormat = dockerNamingPrefix + "%s-%s"

// containerNetworkNameFormat - container network name format
const containerNetworkNameFormat = dockerNamingPrefix + "n"

// objectContainerType - object type for container
type objectContainerType byte

const (
	objectContainerApp     objectContainerType = 'a'
	objectContainerService objectContainerType = 's'
)

// TypeName - get object container type name
func (o objectContainerType) TypeName() string {
	switch o {
	case objectContainerApp:
		{
			return "app"
		}
	case objectContainerService:
		{
			return "service"
		}
	}
	return "unknown"
}

// DockerContainerConfig contains configuration for a Docker container.
type DockerContainerConfig struct {
	projectID  string
	objectType objectContainerType
	objectName string
	command    []string
	Image      string
	Volumes    map[string]string
	Binds      map[string]string
	Env        map[string]string
	WorkingDir string
}

// GetContainerName - get name of docker container
func (d DockerContainerConfig) GetContainerName() string {
	return fmt.Sprintf(containerNameFormat, d.projectID, string(d.objectType), d.objectName)
}

// GetNetworkName - get name for docker network
func (d DockerContainerConfig) GetNetworkName() string {
	return fmt.Sprintf(containerNetworkNameFormat, d.projectID)
}

// GetCommand - get container command
func (d DockerContainerConfig) GetCommand() []string {
	if len(d.command) > 0 {
		return d.command
	}
	return []string{"tail", "-f", "/dev/null"}
}

// GetEnv - convert environment vars to format needed to start docker container
func (d DockerContainerConfig) GetEnv() []string {
	out := make([]string, 0)
	for k, v := range d.Env {
		out = append(out, fmt.Sprintf("%s=%v", k, v))
	}
	return out
}
