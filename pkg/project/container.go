/*
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
*/

package project

import (
	"archive/tar"
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitlab.com/contextualcode/platform_cc/pkg/container"

	"gitlab.com/contextualcode/platform_cc/pkg/output"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/pkg/def"
)

// Container contains information needed to run a container.
type Container struct {
	Config            container.Config
	Name              string
	Definition        interface{}
	Relationships     map[string][]map[string]interface{}
	containerHandler  container.Interface
	configJSON        []byte
	initCommand       string
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
		initCommand:       p.GetDefinitionInitCommand(d),
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
		return errors.WithStack(err)
	}
	// upload config.json
	d2 := output.Duration("Upload config.json.")
	if err := c.Upload(
		"/run/config.json",
		bytes.NewReader(c.configJSON),
	); err != nil {
		return errors.WithStack(err)
	}
	d2()
	// patch
	if c.patchCommand != "" {
		d2 = output.Duration("Patch container.")
		if _, err := c.containerHandler.ContainerCommand(
			c.Config.GetContainerName(),
			"root",
			[]string{"bash", "--login", "-c", c.patchCommand},
			nil,
		); err != nil {
			return errors.WithStack(err)
		}
		d2()
	}
	// run init command
	d2 = output.Duration("Init container.")
	if _, err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", c.initCommand},
		os.Stdout,
	); err != nil {
		return errors.WithStack(err)
	}
	d2()
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
	if _, err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", serviceStartCmd},
		os.Stdout,
	); err != nil {
		return nil, errors.WithStack(err)
	}
	d2()
	// prepare relationships json
	d2 = output.Duration("Parse relationships.")
	relJSONData := map[string]interface{}{
		"relationships": c.Relationships,
	}
	relJSON, err := json.Marshal(relJSONData)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	relB64 := base64.StdEncoding.EncodeToString(relJSON)
	d2()
	// enable authentication it requested
	if err := c.openEnableAuthentication(); err != nil {
		return nil, errors.WithStack(err)
	}
	// open service and retrieve relationships
	d2 = output.Duration("Open service.")
	var openOutput bytes.Buffer
	cmd := fmt.Sprintf(serviceOpenCmd, relB64)
	if _, err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", cmd},
		&openOutput,
	); err != nil {
		return nil, errors.WithStack(err)
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
		return nil, errors.WithStack(err)
	}
	if !containerStatus.Running {
		return nil, errors.WithStack(container.ErrContainerNotRunning)
	}
	if containerStatus.IPAddress == "" {
		return nil, errors.WithStack(ErrContainerNoIP)
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
	if _, err := c.containerHandler.ContainerCommand(
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
	if _, err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", c.buildCommand},
		os.Stdout,
	); err != nil {
		if !errors.Is(err, container.ErrCommandExited) {
			return errors.WithStack(err)
		}
		output.Warn("Build exited with non zero code.")
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
	if _, err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"sh", "-c", c.mountCommand},
		os.Stdout,
	); err != nil {
		return errors.WithStack(err)
	}
	done()
	return nil
}

// Deploy runs the deploy hooks.
func (c Container) Deploy() error {
	done := output.Duration(
		fmt.Sprintf("Running deploy hook for %s '%s.'", c.Config.ObjectType.TypeName(), c.Name),
	)
	if _, err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", appDeployCmd},
		os.Stdout,
	); err != nil {
		if !errors.Is(err, container.ErrCommandExited) {
			return errors.WithStack(err)
		}
		output.Warn("Deploy exited with non zero code.")
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
	if _, err := c.containerHandler.ContainerCommand(
		c.Config.GetContainerName(),
		"root",
		[]string{"bash", "--login", "-c", c.postDeployCommand},
		os.Stdout,
	); err != nil {
		if !errors.Is(err, container.ErrCommandExited) {
			return errors.WithStack(err)
		}
		output.Warn("Post deploy exited with non zero code.")
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
				return errors.WithStack(err)
			}
			stateJSON, err := buildStateJSON(containerStatus.ID[0:12], currentState, desiredState)
			if err != nil {
				return errors.WithStack(err)
			}
			r := bytes.NewReader(stateJSON)
			// issue service state update
			if _, err := c.containerHandler.ContainerShell(
				c.Config.GetContainerName(),
				"root",
				[]string{"bash", "--login", "-c", "/etc/platform/commands/on_service_state_update"},
				r,
			); err != nil {
				return errors.WithStack(err)
			}
			done()
			break
		}
	}
	return nil
}

// Shell accesses the container shell.
func (c Container) Shell(user string, cmd []string) (int, error) {
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
	code, err := c.containerHandler.ContainerShell(
		c.Config.GetContainerName(),
		user,
		cmd,
		nil,
	)
	return code, errors.WithStack(err)
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
	// open logs
	followOption := ""
	if follow {
		followOption = "-f"
	}
	errChan := make(chan error)
	go func(err chan error) {
		_, e := c.containerHandler.ContainerCommand(
			c.Config.GetContainerName(), "root",
			[]string{"sh", "-c", fmt.Sprintf("tail %s /var/log/*.log %s /var/log/*/*.log %s /tmp/*.log", followOption, followOption, followOption)},
			os.Stdout,
		)
		err <- e
	}(errChan)
	// follow logs
	if follow {
		err := <-errChan
		return errors.WithStack(err)
	}
	// wait a second for buffer to fill
	select {
	case err := <-errChan:
		{
			return errors.WithStack(err)
		}
	case <-time.After(time.Second):
		{
			return nil
		}
	}
}

// Commit commits the container.
func (c Container) Commit() error {
	return errors.WithStack(c.containerHandler.ContainerCommit(c.Config.GetContainerName()))
}

// DeleteCommit deletes the commit image.
func (c Container) DeleteCommit() error {
	return errors.WithStack(c.containerHandler.ContainerDeleteCommit(c.Config.GetContainerName()))
}

// Upload uploads given reader to container as a single file at given path.
func (c Container) Upload(path string, reader io.ReadSeeker) error {
	// get size
	size, err := reader.Seek(0, io.SeekEnd)
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = reader.Seek(0, io.SeekStart)
	if err != nil {
		return errors.WithStack(err)
	}
	// build tar
	var buf bytes.Buffer
	tarball := tar.NewWriter(&buf)
	header := &tar.Header{
		Name:  filepath.Base(path),
		Mode:  0664,
		Uname: "root",
		Size:  size,
	}
	if err := tarball.WriteHeader(header); err != nil {
		return errors.WithStack(err)
	}
	if _, err := io.Copy(tarball, reader); err != nil {
		return errors.WithStack(err)
	}
	if err := tarball.Close(); err != nil {
		return errors.WithStack(err)
	}
	if err := tarball.Close(); err != nil {
		return errors.WithStack(err)
	}
	// upload
	return errors.WithStack(c.UploadMulti(
		filepath.Dir(path), &buf,
	))
}

// UploadMulti uploads given tarball reader to container.
func (c Container) UploadMulti(path string, reader io.Reader) error {
	return errors.WithStack(c.containerHandler.ContainerUpload(
		c.Config.GetContainerName(),
		path,
		reader,
	))
}

// Download downloads given container path to given writer.
func (c Container) Download(path string, writer io.Writer) error {
	// download
	var buf bytes.Buffer
	if err := c.DownloadMulti(path, &buf); err != nil {
		return errors.WithStack(err)
	}
	// untar file
	tarball := tar.NewReader(&buf)
	header, err := tarball.Next()
	if err != nil {
		return errors.WithStack(err)
	}
	// copy to writer
	if _, err := io.CopyN(writer, tarball, header.Size); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// DownloadMulti downloads container path and writes to given writer as tarball.
func (c Container) DownloadMulti(path string, writer io.Writer) error {
	return errors.WithStack(c.containerHandler.ContainerDownload(
		c.Config.GetContainerName(),
		path,
		writer,
	))
}
