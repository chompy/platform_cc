package project

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/pkg/stdcopy"

	"gitlab.com/contextualcode/platform_cc/api/container"

	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/def"
)

// Container contains information needed to run a container.
type Container struct {
	Config            container.Config
	Name              string
	Definition        interface{}
	Relationships     map[string][]map[string]interface{}
	containerHandler  container.Interface
	configJSON        []byte
	buildCommand      string
	mountCommand      string
	patchCommand      string
	mountStrategy     string
	postDeployCommand string
}

// NewContainer creates a new container.
func (p *Project) NewContainer(d interface{}) Container {
	configJSON, _ := p.BuildConfigJSON(d)
	o := Container{
		Name:          p.GetDefinitionName(d),
		Definition:    d,
		Relationships: p.GetDefinitionRelationships(d),
		Config: container.Config{
			ProjectID:    p.ID,
			Slot:         p.slot,
			ObjectType:   p.GetDefinitionContainerType(d),
			ObjectName:   p.GetDefinitionName(d),
			Command:      p.GetDefinitionStartCommand(d),
			Image:        p.GetDefinitionImage(d),
			Volumes:      p.GetDefinitionVolumes(d),
			Binds:        p.GetDefinitionBinds(d),
			Env:          p.GetDefinitionEnvironmentVariables(d),
			WorkingDir:   def.AppDir,
			EnableOSXNFS: p.Flags.IsOn(EnableOSXNFSMounts),
		},
		containerHandler:  p.containerHandler,
		configJSON:        configJSON,
		buildCommand:      p.GetDefinitionBuildCommand(d),
		mountCommand:      p.GetDefinitionMountCommand(d),
		patchCommand:      p.GetDefinitionPatch(d),
		mountStrategy:     p.GetOption(OptionMountStrategy),
		postDeployCommand: p.GetDefinitionPostDeployCommand(d),
	}
	return o
}

// Start starts the container.
func (c Container) Start() error {
	done := output.Duration(
		fmt.Sprintf("Start %s '%s.'", c.Config.ObjectType.TypeName(), c.Name),
	)
	// start container
	if err := c.containerHandler.ContainerStart(c.Config); err != nil {
		return tracerr.Wrap(err)
	}
	// upload config.json
	d2 := output.Duration("Upload config.json.")
	if err := c.Upload(
		"/run/config.json",
		bytes.NewReader(c.configJSON),
	); err != nil {
		return tracerr.Wrap(err)
	}
	d2()
	// patch
	if c.patchCommand != "" {
		d2 = output.Duration("Patch container.")
		if err := c.containerHandler.ContainerCommand(
			c.Config.GetContainerName(),
			"root",
			[]string{"bash", "--login", "-c", c.patchCommand},
			nil,
		); err != nil {
			return tracerr.Wrap(err)
		}
		d2()
	}
	done()
	return c.Log()
}

// Open opens the container and returns the relationships.
func (c Container) Open() ([]map[string]interface{}, error) {
	indentLevel := output.IndentLevel
	done := output.Duration(
		fmt.Sprintf("Open %s '%s.'", c.Config.ObjectType.TypeName(), c.Name),
	)
	// start service
	d2 := output.Duration("Start service.")
	if err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", serviceStartCmd},
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
	// enable authentication it requested
	if err := c.openEnableAuthentication(); err != nil {
		return nil, tracerr.Wrap(err)
	}
	// open service and retrieve relationships
	d2 = output.Duration("Open service.")
	var openOutput bytes.Buffer
	cmd := fmt.Sprintf(serviceOpenCmd, relB64)
	if err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", cmd},
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
	containerStatus, err := c.containerHandler.ContainerStatus(c.Config.GetContainerName())
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	if !containerStatus.Running {
		return nil, tracerr.New("container not running")
	}
	if containerStatus.IPAddress == "" {
		return nil, tracerr.New("container has no ip address")
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
		rel["ip"] = containerStatus.IPAddress
		out = append(out, rel)
	}
	d2()
	done()
	output.IndentLevel = indentLevel
	return out, nil
}

// HasBuild returns true if container has been built.
func (c Container) HasBuild() bool {
	var buf bytes.Buffer
	if err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", "[ -f /config/built ] && echo 'YES'"},
		&buf,
	); err != nil {
		output.LogError(err)
		return false
	}
	return strings.TrimSpace(buf.String()) == "YES"
}

// Build runs the build hooks.
func (c Container) Build() error {
	if c.buildCommand == "" {
		output.LogDebug(
			fmt.Sprintf("Skip build for %s, no build command defined.", c.Config.GetContainerName()),
			nil,
		)
		return nil
	}
	done := output.Duration(
		fmt.Sprintf("Building %s '%s.'", c.Config.ObjectType.TypeName(), c.Name),
	)
	// run command
	if err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", c.buildCommand},
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
		fmt.Sprintf("Set up mounts for %s '%s' using '%s' strategy.", c.Config.ObjectType.TypeName(), c.Name, c.mountStrategy),
	)
	// run command
	if err := c.containerHandler.ContainerCommand(
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
	if err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", appDeployCmd},
		os.Stdout,
	); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// PostDeploy runs the post deploy hooks.
func (c Container) PostDeploy() error {
	if c.postDeployCommand == "" {
		return nil
	}
	done := output.Duration(
		fmt.Sprintf("Running post-deploy hook for %s '%s.'", c.Config.ObjectType.TypeName(), c.Name),
	)
	if err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", c.postDeployCommand},
		os.Stdout,
	); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// openEnableAuthentication enables authentication in the service.
func (c Container) openEnableAuthentication() error {
	switch c.Definition.(type) {
	case def.Service:
		{
			// check that authentication is enabled
			serviceConfig := c.Definition.(def.Service).Configuration
			if !serviceConfig.IsAuthenticationEnabled() {
				return nil
			}
			done := output.Duration("Enable authentication.")
			// build state json
			currentState := getDefaultServiceState()
			currentState.Image = c.Config.Image
			desiredState := getDefaultServiceState()
			desiredState.Image = c.Config.Image
			desiredState.Configuration = serviceConfig
			containerStatus, err := c.containerHandler.ContainerStatus(c.Config.GetContainerName())
			if err != nil {
				return tracerr.Wrap(err)
			}
			stateJSON, err := buildStateJSON(containerStatus.ID[0:12], currentState, desiredState)
			if err != nil {
				return tracerr.Wrap(err)
			}
			r := bytes.NewReader(stateJSON)
			// issue service state update
			if err := c.containerHandler.ContainerShell(
				c.Config.GetContainerName(),
				"root",
				[]string{"bash", "--login", "-c", "/etc/platform/commands/on_service_state_update"},
				r,
			); err != nil {
				return tracerr.Wrap(err)
			}
			done()
			break
		}
	}
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
		cmd = []string{"bash", "--login"}
	}
	return tracerr.Wrap(c.containerHandler.ContainerShell(
		c.Config.GetContainerName(),
		user,
		cmd,
		nil,
	))
}

// Log outputs container logs to log file.
func (c Container) Log() error {
	output.LogInfo(fmt.Sprintf("Read logs for container '%s.'", c.Config.GetContainerName()))
	go func() {
		out, err := c.containerHandler.ContainerLog(c.Config.GetContainerName(), true)
		if err != nil {
			output.LogError(err)
			return
		}
		scanner := bufio.NewScanner(out)
		defer out.Close()
		for {
			for scanner.Scan() {
				output.LogDebug(fmt.Sprintf("[%s] %s", c.Config.GetContainerName(), scanner.Text()), nil)
			}
			if err := scanner.Err(); err != nil {
				output.LogError(err)
			}
		}
	}()
	return nil
}

// LogStdout dumps container log to stdout.
func (c Container) LogStdout(follow bool) error {
	output.LogInfo(fmt.Sprintf("Read logs for container '%s.'", c.Config.GetContainerName()))
	// get log stream
	out, err := c.containerHandler.ContainerLog(c.Config.GetContainerName(), follow)
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer out.Close()
	var buf bytes.Buffer
	// inline func, scan lines in buffer and output to stdout
	scanLines := func() error {
		scanner := bufio.NewScanner(&buf)
		for scanner.Scan() {
			output.ContainerLog(c.Config.GetContainerName(), scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return tracerr.Wrap(err)
		}
		return nil
	}
	// follow logs
	if follow {
		go func() {
			for {
				err := scanLines()
				if err != nil {
					output.LogError(err)
					return
				}
				buf.Reset()
			}
		}()
	}
	// copy to buf
	_, err = stdcopy.StdCopy(&buf, &buf, out)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return tracerr.Wrap(scanLines())
}

// Commit commits the container.
func (c Container) Commit() error {
	return tracerr.Wrap(c.containerHandler.ContainerCommit(c.Config.GetContainerName()))
}

// DeleteCommit deletes the commit image.
func (c Container) DeleteCommit() error {
	return tracerr.Wrap(c.containerHandler.ContainerDeleteCommit(c.Config.GetContainerName()))
}

// Upload uploads given reader to container as a file at given path.
func (c Container) Upload(path string, reader io.Reader) error {
	return tracerr.Wrap(c.containerHandler.ContainerUpload(
		c.Config.GetContainerName(),
		path,
		reader,
	))
}

// Download downloads given container path to given writer.
func (c Container) Download(path string, writer io.Writer) error {
	return tracerr.Wrap(c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"cat", path},
		writer,
	))
}
