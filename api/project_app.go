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

const appContainerCmd = `
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
chown -R web:web /run
python /tmp/fake-rpc.py &> /tmp/fake-rpc.log &
until [ -f /run/config.json ]; do sleep 1; done
runsvdir -P /etc/service &> /tmp/runsvdir.log &
/etc/platform/boot
exec init
/etc/platform/start
`

// getAppContainerConfig - get container configuration for app
func (p *Project) getAppContainerConfig(app *AppDef) dockerContainerConfig {
	return dockerContainerConfig{
		projectID:  p.ID,
		objectName: app.Name,
		objectType: objectContainerApp,
		command:    []string{"sh", "-c", appContainerCmd},
		Image:      app.GetContainerImage(),
		Binds: map[string]string{
			app.Path: "/app",
		},
		Volumes: map[string]string{
			fmt.Sprintf(containerVolumeNameFormat, p.ID, app.Name): "/mnt/storage",
		},
	}
}

// startApp - start an app
func (p *Project) startApp(app *AppDef) error {
	log.Printf("Start app '%s.'", app.Name)
	// get container config
	containerConfig := p.getAppContainerConfig(app)
	// start container
	if err := p.docker.StartContainer(containerConfig); err != nil {
		return err
	}
	// build config.json
	configJSON, err := p.BuildConfigJSON()
	if err != nil {
		return err
	}
	configJSONReader := bytes.NewReader(configJSON)
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

// openApp - make the application available
func (p *Project) openApp(app *AppDef) error {
	log.Printf("Opening app '%s.'", app.Name)
	// get container config
	containerConfig := p.getAppContainerConfig(app)
	// create relationships payload
	relationshipsVar, err := p.getAppRelationships(app)
	if err != nil {
		return err
	}
	relationships := map[string]interface{}{
		"relationships": relationshipsVar,
	}
	relationshipsJSON, err := json.Marshal(relationships)
	log.Println(string(relationshipsJSON))
	if err != nil {
		return err
	}
	relationshipsB64 := base64.StdEncoding.EncodeToString(relationshipsJSON)
	cmd := fmt.Sprintf(`
sleep 3
echo '%s' | base64 -d | /etc/platform/commands/open
	`, relationshipsB64)
	// run open command
	return p.docker.RunContainerCommand(
		containerConfig.GetContainerName(),
		"root",
		[]string{"sh", "-c", cmd},
		os.Stdout,
	)
}

// getAppRelationships - generate relationships variable for app
func (p *Project) getAppRelationships(app *AppDef) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	for name, rel := range app.Relationships {
		out[name] = nil
		relSplit := strings.Split(rel, ":")
		for _, service := range p.Services {
			if service.Name == relSplit[0] {
				relationships, err := p.getServiceRelationships(service)
				if err != nil {
					return nil, err
				}
				for _, relationship := range relationships {
					if relationship["rel"] == relSplit[1] {
						out[name] = []map[string]interface{}{
							relationship,
						}
					}
				}

			}

		}
	}
	return out, nil
}

// getAppVariables - get variables to inject in to container
func (p *Project) getAppVariables(app *AppDef) map[string]string {
	out := make(map[string]string)
	for varType, varVal := range app.Variables {
		for k, v := range varVal {
			switch v.(type) {
			case string:
				{
					out[fmt.Sprintf("%s:%s", strings.ToLower(varType), k)] = v.(string)
					break
				}
			}
		}
	}
	return out
}
