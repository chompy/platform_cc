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
	"os/user"
	"strconv"

	"gitlab.com/contextualcode/platform_cc/api/def"
)

// BuildConfigJSON makes config.json for container runtime.
func (p *Project) BuildConfigJSON(d interface{}) ([]byte, error) {
	// determine uid/gid to set
	uid := 0
	gid := 0
	currentUser, _ := user.Current()
	if currentUser != nil {
		uid, _ = strconv.Atoi(currentUser.Uid)
		gid, _ = strconv.Atoi(currentUser.Gid)
	}
	if uid == 0 {
		uid = 1000
	}
	if gid == 0 {
		gid = 1000
	}
	// get host ip
	hostIP, err := p.docker.GetNetworkHostIP()
	if err != nil {
		hostIP = "-"
	}
	// get name
	name := ""
	appName := ""
	switch d.(type) {
	case def.App:
		{
			appName = d.(def.App).Name
			name = appName
			break
		}
	case def.AppWorker:
		{
			name = d.(def.AppWorker).Name
			break
		}
	case def.Service:
		{
			name = d.(def.Service).Name
			break
		}
	}
	// build application json
	appJsons := make([]map[string]interface{}, 0)
	if appName != "" {
		appJsons = append(appJsons, p.buildConfigAppJSON(d))
	}
	for _, app := range p.Apps {
		if appName == "" && app.Name == appName {
			continue
		}
		appJsons = append(appJsons, p.buildConfigAppJSON(app))
	}
	out := map[string]interface{}{
		"primary_ip":   "127.0.0.1",
		"features":     []string{},
		"domainname":   fmt.Sprintf("%s.pcc.local", p.ID),
		"host_ip":      hostIP,
		"applications": appJsons,
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
				"ssh_key": "-",
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
		"hostname":   hostIP,
		"instance":   p.ID,
		"nameserver": "127.0.0.11",
		"web_uid":    uid,
		"web_gid":    gid,
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
	switch d.(type) {
	case def.Service:
		{
			out["hosts"] = map[string]interface{}{}
			out["configuration"] = d.(def.Service).Configuration
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
	switch d.(type) {
	case def.App:
		{
			name = d.(def.App).Name
			if p.Flags.Has(EnableCron) {
				crons = d.(def.App).Crons
			}
			appType = d.(def.App).Type
			hooks = d.(def.App).Hooks
			disk = d.(def.App).Disk
			mounts = d.(def.App).Mounts
			runtime = d.(def.App).Runtime
			appWebo := d.(def.App).Web
			appWeb = &appWebo
			dependencies = d.(def.App).Dependencies
			worker = nil
			break
		}
	case def.AppWorker:
		{
			name = d.(def.AppWorker).Name
			appType = d.(def.AppWorker).Type
			disk = d.(def.AppWorker).Disk
			mounts = d.(def.AppWorker).Mounts
			workero := d.(def.AppWorker)
			worker = &workero
			appWeb = nil
			break
		}
	}
	// build configuration section
	configuration := map[string]interface{}{
		"app_dir":       def.AppDir,
		"hooks":         hooks,
		"variables":     p.GetPlatformEnvironmentVariables(d),
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
		"project_info": map[string]interface{}{
			"ssh_key": "",
		},
	}
	if appWeb != nil {
		configuration["web"] = appWeb
	}
	if worker != nil {
		configuration["worker"] = worker
		configuration["mounts"] = worker.Mounts
	}
	return map[string]interface{}{
		"name":                  name,
		"crons":                 crons,
		"enable_smtp":           "false",
		"mounts":                mounts,
		"hooks":                 hooks,
		"cron_minimum_interval": "1",
		"dependencies":          dependencies,
		"configuration":         configuration,
	}

}
