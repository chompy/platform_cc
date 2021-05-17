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
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/pkg/config"

	"gitlab.com/contextualcode/platform_cc/pkg/def"
)

// BuildConfigJSON makes config.json for container runtime.
func (p *Project) BuildConfigJSON(d interface{}) ([]byte, error) {
	// determine uid/gid to set
	uid, gid := p.getUID()
	// get name + build app json
	name := ""
	appJsons := make([]map[string]interface{}, 0)
	switch d := d.(type) {
	case def.App:
		{
			name = d.Name
			// ensure this app is the first item in application json list
			appJsons = append(appJsons, p.buildConfigAppJSON(d))
			// build rest of application json
			for _, app := range p.Apps {
				if app.Name == name {
					continue
				}
				appJsons = append(appJsons, p.buildConfigAppJSON(app))
			}
			break
		}
	case *def.AppWorker:
		{
			name = d.Name
			// put worker in application json list
			appJsons = append(appJsons, p.buildConfigAppJSON(d))
			break
		}
	case def.Service:
		{
			name = d.Name
			break
		}
	}
	// get private key
	privKey, err := config.PrivateKey()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	out := map[string]interface{}{
		"primary_ip":    "127.0.0.1",
		"features":      []string{},
		"domainname":    fmt.Sprintf("%s.pcc.local", p.ID),
		"host_ip":       "127.0.0.1",
		"all_node_ips":  []string{"127.0.0.1"},
		"all_hostnames": []string{"localhost"},
		"applications":  appJsons,
		"configuration": map[string]interface{}{
			"access": map[string]interface{}{
				"ssh": []string{},
			},
			"privileged_digest": "44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
			"environment_info": map[string]interface{}{
				"is_production": false,
				"machine_name":  "pcc-1",
				"name":          "pcc",
				"reference":     "refs/heads/pcc",
				"is_main":       false,
			},
			"project_info": map[string]interface{}{
				"name":    p.ID,
				"ssh_key": privKey,
				"settings": map[string]interface{}{
					"systemd":          false,
					"variables_prefix": "PLATFORM_",
					"crons_in_git":     false,
					"product_code":     "platformsh",
					"product_name":     "Platform.sh",
					"enforce_mfa":      false,
					"bot_email":        "bot@platform.sh",
				},
			},
			"privileged": map[string]interface{}{},
		},
		"info": map[string]interface{}{
			"mail_relay_host":    "",
			"mail_relay_host_v2": "127.0.0.1",
			"limits": map[string]interface{}{
				"disk":   p.Apps[0].Disk,
				"cpu":    1.0,
				"memory": 1024,
			},
			"external ip": "127.0.0.1",
		},
		"name":       p.ID,
		"service":    name,
		"cluster":    "-",
		"region":     "pcc.local",
		"hostname":   "pcc.local",
		"instance":   p.ID,
		"nameserver": "127.0.0.11",
		"web_uid":    uid,
		"web_gid":    gid,
		"log_uid":    uid,
		"log_gid":    gid,
		"log_file":   "/dev/stdout",
		"nginx": map[string]interface{}{
			"headers_prefix":         "X-PLATFORM-",
			"policy":                 nil,
			"mappings":               map[string]interface{}{},
			"upstream_address":       "/run/app.sock",
			"preflight_block_policy": "FULL",
			"error_codes":            map[string]interface{}{},
		},
		"workers": 2,
	}
	switch d := d.(type) {
	case def.Service:
		{
			out["hosts"] = map[string]interface{}{}
			out["configuration"] = d.Configuration
			break
		}
	}
	return json.Marshal(out)
}

// buildConfigAppJSON builds the application section of config.json.
func (p *Project) buildConfigAppJSON(d interface{}) map[string]interface{} {
	// grab variables for given def
	name := ""
	crons := map[string]*def.AppCron{}
	hooks := def.AppHooks{}
	disk := 0
	appType := ""
	mounts := map[string]*def.AppMount{}
	worker := &def.AppWorker{}
	runtime := def.AppRuntime{}
	appWeb := &def.AppWeb{}
	dependencies := def.AppDependencies{}
	build := def.AppBuild{}
	switch d := d.(type) {
	case def.App:
		{
			name = d.Name
			if p.HasFlag(EnableCron) {
				crons = d.Crons
			}
			appType = d.Type
			hooks = d.Hooks
			disk = d.Disk
			mounts = d.Mounts
			runtime = d.Runtime
			appWebo := d.Web
			appWeb = &appWebo
			dependencies = d.Dependencies
			worker = nil
			build = d.Build
			break
		}
	case *def.AppWorker:
		{
			name = d.Name
			appType = d.Type
			disk = d.Disk
			mounts = d.Mounts
			runtime = d.Runtime
			workero := d
			worker = workero
			appWeb = nil
			break
		}
	}
	// build configuration section
	configuration := map[string]interface{}{
		"app_dir":       def.AppDir,
		"hooks":         hooks,
		"variables":     p.GetDefinitionEnvironmentVariables(d),
		"timezone":      nil,
		"disk":          disk,
		"slug_id":       "-",
		"size":          "AUTO",
		"relationships": p.GetDefinitionRelationships(d),
		"is_production": false,
		"name":          name,
		"access":        map[string]string{},
		"preflight": map[string]interface{}{
			"enabled":       true,
			"ignored_rules": []string{},
		},
		"tree_id":   "-",
		"mounts":    mounts,
		"runtime":   runtime,
		"type":      appType,
		"crons":     crons,
		"slug":      "-",
		"resources": nil,
	}
	if appWeb != nil {
		configuration["web"] = appWeb
	}
	if worker != nil {
		configuration["worker"] = worker
	}
	privKey, _ := config.PrivateKey()
	return map[string]interface{}{
		"name":                  name,
		"build":                 build,
		"crons":                 crons,
		"enable_smtp":           "false",
		"mounts":                mounts,
		"hooks":                 hooks,
		"cron_minimum_interval": "1",
		"dependencies":          dependencies,
		"configuration":         configuration,
		"project_info": map[string]interface{}{
			"ssh_key": privKey,
		},
	}
}
