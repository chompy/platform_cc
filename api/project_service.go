package api

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

const serviceContainerCmd = `
umount /etc/hosts
umount /etc/resolv.conf
mkdir -p /run/shared /run/rpc_pipefs/nfs
cat >/tmp/fake-rpc.py <<EOF
from gevent.monkey import patch_all;
patch_all();
from gevent_jsonrpc import RpcServer;
import json;
RpcServer(
	"/run/shared/agent.sock",
	"foo",
	root=None,
	root_factory=lambda c,a: c.send(json.dumps({"jsonrpc":"2.0","result":True,"id": json.loads(c.recv(1024))["id"]})))._accepter_greenlet.get();
EOF
python /tmp/fake-rpc.py &> /tmp/fake-rpc.log &
sleep 1
runsvdir -P /etc/service &> /tmp/runsvdir.log &
sleep 1
until [ -f /run/config.json ]; do sleep 1; done
/etc/platform/boot
exec init
`

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
	// get service config
	serviceConfig := def.GetServiceConfig()
	if serviceConfig == nil {
		return fmt.Errorf("service config not found for service of type %s", def.Type)
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

	// platform start cmd
	platformStartCmd := `
until [ -f /run/config.json ]; do sleep 1; done
sleep 1
/etc/platform/start &
`
	p.docker.RunContainerCommand(
		containerConfig.GetContainerName(),
		"root",
		[]string{"sh", "-c", platformStartCmd},
		os.Stdout,
	)
	// setup command
	setupCmd, err := serviceConfig.GetSetupCommand(def)
	if err != nil {
		return err
	}
	if len(setupCmd) > 0 {
		log.Print("Post start commands for service")
		if err := p.docker.RunContainerCommand(
			containerConfig.GetContainerName(),
			"root",
			setupCmd,
			nil,
		); err != nil {
			return err
		}
	}
	return nil
}

// getServiceRelationships - get service relationship value
func (p *Project) getServiceRelationships(def *ServiceDef) ([]map[string]interface{}, error) {
	// get service config
	serviceConfig := def.GetServiceConfig()
	if serviceConfig == nil {
		return nil, fmt.Errorf("service config not found for service of type %s", def.Type)
	}
	// get relationships
	out, err := serviceConfig.GetRelationship(def)
	if err != nil {
		return []map[string]interface{}{}, err
	}
	// get container ip address
	ipAddress, err := p.docker.GetContainerIP(p.getServiceContainerConfig(def).GetContainerName())
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
