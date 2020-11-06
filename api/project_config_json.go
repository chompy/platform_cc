package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os/user"
	"strconv"
	"strings"
)

const entropySalt = "Dyt+&&*^dKfD9,$rZRA$|I^DLKr%<By"

// BuildConfigJSON - make config.json for container runtime
func (p *Project) BuildConfigJSON(def interface{}) ([]byte, error) {
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
	switch def.(type) {
	case *ServiceDef:
		{
			c := p.getServiceContainerConfig(def.(*ServiceDef))
			hostname = c.GetContainerName()
			serviceName = def.(*ServiceDef).Name
			break
		}
	case *AppDef:
		{
			c := p.getAppContainerConfig(def.(*AppDef))
			hostname = c.GetContainerName()
			serviceName = def.(*AppDef).Name
			break
		}
	}
	// build application json
	appJsons := make([]map[string]interface{}, 0)
	for _, app := range p.Apps {
		appJsons = append(appJsons, p.buildConfigAppJSON(app))
	}
	return json.Marshal(map[string]interface{}{
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
		"nameserver": "1.1.1.1",
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
	})

}

// buildConfigAppJSON - build application section of config.json
func (p *Project) buildConfigAppJSON(app *AppDef) map[string]interface{} {
	// build PLATFORM_ROUTES
	routes := make(map[string]*RouteDef)
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
	// build environment vars
	envVars := map[string]string{
		"PLATFORM_DOCUMENT_ROOT":    "/app/web",
		"PLATFORM_APPLICATION":      app.BuildPlatformApplicationVar(),
		"PLATFORM_PROJECT":          "-",
		"PLATFORM_PROJECT_ENTROPY":  fmt.Sprintf("%x", entH.Sum(nil)),
		"PLATFORM_APPLICATION_NAME": app.Name,
		"PLATFORM_BRANCH":           "pcc",
		"PLATFORM_DIR":              appDir,
		"PLATFORM_TREE_ID":          "-",
		"PLATFORM_ENVIRONMENT":      "pcc",
		"PLATFORM_VARIABLES":        app.BuildPlatformVariablesVar(),
		"PLATFORM_ROUTES":           routesJSONB64,
	}
	for k, v := range app.Variables["env"] {
		envVars[k] = v.(string)
	}
	return map[string]interface{}{
		"crons":                 app.Crons,
		"enable_smtp":           "false",
		"mounts":                app.Mounts,
		"cron_minimum_interval": "1",
		"configuration": map[string]interface{}{
			"app_dir":       appDir,
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
