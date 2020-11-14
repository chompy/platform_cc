package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/user"
	"strconv"
	"strings"

	"gitlab.com/contextualcode/platform_cc/def"
)

const entropySalt = "Dyt+&&*^dKfD9,$rZRA$|I^DLKr%<By"

// BuildConfigJSON - make config.json for container runtime
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
	hostIP, err := p.docker.GetNetworkHostIP(p.ID)
	if err != nil {
		hostIP = "-"
	}
	// grab values from app or service def
	serviceName := "-"
	hostname := "-"
	switch d.(type) {
	case *def.Service:
		{
			c := p.GetServiceContainerConfig(d.(*def.Service))
			hostname = c.GetContainerName()
			serviceName = d.(*def.Service).Name
			break
		}
	case *def.App:
		{
			c := p.GetAppContainerConfig(d.(*def.App))
			hostname = c.GetContainerName()
			serviceName = d.(*def.App).Name
			break
		}
	}
	// build application json
	appJsons := make([]map[string]interface{}, 0)
	for _, app := range p.Apps {
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
		"service":    serviceName,
		"cluster":    "-",
		"region":     "pcc.local",
		"hostname":   hostname,
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
	case *def.Service:
		{
			out["hosts"] = map[string]interface{}{}
			out["configuration"] = d.(*def.Service).Configuration
			break
		}
	}
	return json.Marshal(out)

}

// buildConfigAppJSON - build application section of config.json
func (p *Project) buildConfigAppJSON(app *def.App) map[string]interface{} {
	// build PLATFORM_ROUTES
	routes := make(map[string]*def.Route)
	for _, route := range p.Routes {
		if strings.HasPrefix(route.Path, ".") {
			continue
		}
		routes[route.Path] = route
	}
	routesJSON, _ := json.Marshal(routes)
	routesJSONB64 := base64.StdEncoding.EncodeToString(routesJSON)
	// build PLATFORM_PROJECT_ENTROPY
	entH := md5.New()
	entH.Write([]byte(entropySalt))
	entH.Write([]byte(p.ID))
	entH.Write([]byte(entropySalt))
	// build PLATFORM_VARIABLES
	appVars := p.getAppVariables(app)
	appVarsJSON, _ := json.Marshal(appVars)
	appVarsB64 := base64.StdEncoding.EncodeToString(appVarsJSON)
	// build environment vars
	envVars := map[string]string{
		"PLATFORM_DOCUMENT_ROOT":    "/app/web",
		"PLATFORM_APPLICATION":      app.BuildPlatformApplicationVar(),
		"PLATFORM_PROJECT":          "-",
		"PLATFORM_PROJECT_ENTROPY":  fmt.Sprintf("%x", entH.Sum(nil)),
		"PLATFORM_APPLICATION_NAME": app.Name,
		"PLATFORM_BRANCH":           "pcc",
		"PLATFORM_DIR":              def.AppDir,
		"PLATFORM_TREE_ID":          "-",
		"PLATFORM_ENVIRONMENT":      "pcc",
		"PLATFORM_VARIABLES":        appVarsB64,
		"PLATFORM_ROUTES":           routesJSONB64,
	}
	for k, v := range app.Variables["env"] {
		switch v.(type) {
		case int:
			{
				envVars[k] = fmt.Sprintf("%d", v.(int))
				break
			}
		case string:
			{
				envVars[k] = v.(string)
				break
			}
		}
	}
	for k, v := range p.Variables["env"] {
		envVars[k] = v
	}
	return map[string]interface{}{
		"name":                  app.Name,
		"crons":                 app.Crons,
		"enable_smtp":           "false",
		"mounts":                app.Mounts,
		"hooks":                 app.Hooks,
		"cron_minimum_interval": "1",
		"dependencies":          app.Dependencies,
		"configuration": map[string]interface{}{
			"app_dir":       def.AppDir,
			"hooks":         app.Hooks,
			"variables":     envVars,
			"timezone":      nil,
			"disk":          app.Disk,
			"slug_id":       "-",
			"size":          "AUTO",
			"relationships": app.Relationships,
			"web":           app.Web,
			"is_production": false,
			"name":          app.Name,
			"access":        map[string]string{},
			"preflight": map[string]interface{}{
				"enabled":       true,
				"ignored_rules": []string{},
			},
			"tree_id":   "-",
			"mounts":    app.Mounts,
			"runtime":   app.Runtime,
			"type":      app.Type,
			"crons":     app.Crons,
			"slug":      "-",
			"resources": nil,
		},
	}
}
