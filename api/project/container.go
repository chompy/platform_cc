package project

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/docker"
)

// Container contains information needed to run a container.
type Container struct {
	Config        docker.ContainerConfig
	Name          string
	Definition    interface{}
	Relationships map[string][]map[string]interface{}
	docker        docker.Client
	configJSON    []byte
	buildCommand  string
	mountCommand  string
}

// NewContainer creates a new container.
func (p *Project) NewContainer(d interface{}) Container {
	typeName := strings.Split(p.GetDefinitionType(d), ":")
	configJSON, _ := p.BuildConfigJSON(d)
	return Container{
		Name:          p.GetDefinitionName(d),
		Definition:    d,
		Relationships: p.GetDefinitionRelationships(d),
		Config: docker.ContainerConfig{
			ProjectID:  p.ID,
			ObjectType: p.GetDefinitionContainerType(d),
			ObjectName: p.GetDefinitionName(d),
			Command:    p.GetDefinitionStartCommand(d),
			Image:      fmt.Sprintf("%s%s-%s", platformShDockerImagePrefix, typeName[0], typeName[1]),
			Volumes:    p.GetDefinitionVolumes(d),
			Binds:      p.GetDefinitionBinds(d),
			Env:        p.GetDefinitionEnvironmentVariables(d),
			WorkingDir: def.AppDir,
		},
		docker:       p.docker,
		configJSON:   configJSON,
		buildCommand: p.GetDefinitionBuildCommand(d),
		mountCommand: p.GetDefinitionMountCommand(d),
	}
}

// Start starts the container.
func (c Container) Start() error {
	done := output.Duration(
		fmt.Sprintf("Start %s '%s.'", c.Config.ObjectType.TypeName(), c.Name),
	)
	// start container
	if err := c.docker.StartContainer(c.Config); err != nil {
		return tracerr.Wrap(err)
	}
	// upload config.json
	d2 := output.Duration("Upload config.json.")
	if err := c.docker.UploadDataToContainer(
		c.Config.GetContainerName(),
		"/run/config.json",
		bytes.NewReader(c.configJSON),
	); err != nil {
		return tracerr.Wrap(err)
	}
	d2()
	done()
	return nil
}

// Open opens the container and returns the relationships.
func (c Container) Open() ([]map[string]interface{}, error) {
	indentLevel := output.IndentLevel
	done := output.Duration(
		fmt.Sprintf("Open %s '%s.'", c.Config.ObjectType.TypeName(), c.Name),
	)
	// start service
	d2 := output.Duration("Start service.")
	if err := c.docker.RunContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"sh", "-c", serviceStartCmd},
		os.Stdout,
	); err != nil {
		return nil, tracerr.Wrap(err)
	}
	d2()
	// prepare relationships json
	d2 = output.Duration("Parse relationships.")
	relJSONData := map[string]interface{}{
		"relationships": c.Relationships,
	}
	relJSON, err := json.Marshal(relJSONData)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	relB64 := base64.StdEncoding.EncodeToString(relJSON)
	d2()
	// open service and retrieve relationships
	d2 = output.Duration("Open service.")
	var openOutput bytes.Buffer
	cmd := fmt.Sprintf(serviceOpenCmd, relB64)
	if err := c.docker.RunContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"sh", "-c", cmd},
		&openOutput,
	); err != nil {
		return nil, tracerr.Wrap(err)
	}
	d2()
	// process output relationships
	d2 = output.Duration("Build relationship.")
	openOutlineLines := bytes.Split(openOutput.Bytes(), []byte{'\n'})
	rlRaw := openOutlineLines[len(openOutlineLines)-1]
	data := make(map[string]interface{})
	json.Unmarshal(rlRaw, &data)
	// get ip address
	ipAddress, err := c.docker.GetContainerIP(c.Config.GetContainerName())
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	out := make([]map[string]interface{}, 0)
	for k, v := range data {
		rel := GetDefinitionEmptyRelationship(c.Definition)
		for kk, vv := range v.(map[string]interface{}) {
			rel[kk] = vv
		}
		rel["rel"] = k
		rel["host"] = c.Config.GetContainerName()
		rel["hostname"] = c.Config.GetContainerName()
		rel["ip"] = ipAddress
		out = append(out, rel)
	}
	d2()
	done()
	output.IndentLevel = indentLevel
	return out, nil
}

// Build runs the build hooks.
func (c Container) Build() error {
	if c.buildCommand == "" {
		return nil
	}
	done := output.Duration(
		fmt.Sprintf("Building %s '%s.'", c.Config.ObjectType.TypeName(), c.Name),
	)
	// run command
	if err := c.docker.RunContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"sh", "-c", c.buildCommand},
		os.Stdout,
	); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// SetupMounts sets up mounts in container.
func (c Container) SetupMounts() error {
	if c.mountCommand == "" {
		return nil
	}
	done := output.Duration(
		fmt.Sprintf("Set up mounts for %s '%s.'", c.Config.ObjectType.TypeName(), c.Name),
	)
	// run command
	if err := c.docker.RunContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"sh", "-c", c.mountCommand},
		os.Stdout,
	); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// Deploy runs the deploy hooks.
func (c Container) Deploy() error {
	done := output.Duration(
		fmt.Sprintf("Running deploy hook for %s '%s.'", c.Config.ObjectType.TypeName(), c.Name),
	)
	if err := c.docker.RunContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"sh", "-c", appDeployCmd},
		os.Stdout,
	); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// Shell accesses the container shell.
func (c Container) Shell(user string, cmd []string) error {
	output.Info(
		fmt.Sprintf(
			"Access shell for %s '%s.'",
			c.Config.ObjectType.TypeName(),
			c.Name,
		),
	)
	if len(cmd) == 0 {
		cmd = []string{"bash"}
	}
	return tracerr.Wrap(c.docker.ShellContainer(
		c.Config.GetContainerName(),
		user,
		cmd,
		nil,
	))
}
