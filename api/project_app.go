package api

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// getAppContainerConfig - get container configuration for app
func (p *Project) getAppContainerConfig(app *AppDef) dockerContainerConfig {
	uid, gid := p.getUID()
	cmd := fmt.Sprintf(
		appContainerCmd,
		uid+1,
		gid+1,
		uid,
		gid,
	)
	return dockerContainerConfig{
		projectID:  p.ID,
		objectName: app.Name,
		objectType: objectContainerApp,
		command:    []string{"sh", "-c", cmd},
		Image:      app.GetContainerImage(),
		Binds: map[string]string{
			app.Path: "/app",
		},
		Volumes: map[string]string{
			fmt.Sprintf(containerVolumeNameFormat, p.ID, app.Name): "/mnt/storage",
		},
		Env: p.getAppEnvironmentVariables(app),
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
	configJSON, err := p.BuildConfigJSON(app)
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
	if err != nil {
		return err
	}
	relationshipsB64 := base64.StdEncoding.EncodeToString(relationshipsJSON)
	cmd := fmt.Sprintf(appOpenCmd, relationshipsB64)
	// run open command
	return p.docker.RunContainerCommand(
		containerConfig.GetContainerName(),
		"root",
		[]string{"sh", "-c", cmd},
		os.Stdout,
	)
}

// BuildApp - build the application
func (p *Project) BuildApp(app *AppDef) error {
	log.Printf("Building app '%s.'", app.Name)
	// get container config
	containerConfig := p.getAppContainerConfig(app)
	// upload build script
	buildScriptReader := strings.NewReader(appBuildScript)
	if err := p.docker.UploadDataToContainer(
		containerConfig.GetContainerName(),
		"/opt/build.py",
		buildScriptReader,
	); err != nil {
		return err
	}
	// build data
	uid, gid := p.getUID()
	buildData := map[string]interface{}{
		"application": p.buildConfigAppJSON(app),
		"source_dir":  appDir,
		"output_dir":  appDir,
		"cache_dir":   "/tmp/cache",
		"uid":         uid,
		"gid":         gid,
	}
	buildJSON, err := json.Marshal(buildData)
	if err != nil {
		return err
	}
	buildB64 := base64.StdEncoding.EncodeToString(buildJSON)
	// run command
	if err := p.docker.RunContainerCommand(
		containerConfig.GetContainerName(),
		"root",
		[]string{"sh", "-c",
			fmt.Sprintf(appBuildCmd, buildB64),
		},
		os.Stdout,
	); err != nil {
		return err
	}
	return nil
}

// getAppRelationships - generate relationships variable for app
func (p *Project) getAppRelationships(app *AppDef) (map[string][]map[string]interface{}, error) {
	out := make(map[string][]map[string]interface{})
	for name, rel := range app.Relationships {
		out[name] = make([]map[string]interface{}, 0)
		relSplit := strings.Split(rel, ":")
		for _, v := range p.relationships {
			if v["service"].(string) == relSplit[0] && v["rel"] == relSplit[1] {
				out[name] = append(out[name], v)
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
			case int:
				{
					out[fmt.Sprintf("%s:%s", strings.ToLower(varType), k)] = fmt.Sprintf("%d", v.(int))
					break
				}
			case string:
				{
					out[fmt.Sprintf("%s:%s", strings.ToLower(varType), k)] = v.(string)
					break
				}
			}
		}
	}
	for varType, varVal := range p.Variables {
		for k, v := range varVal {
			out[fmt.Sprintf("%s:%s", strings.ToLower(varType), k)] = v
		}
	}
	return out
}

// getAppEnvironmentVariables - get application environment variables
func (p *Project) getAppEnvironmentVariables(app *AppDef) map[string]string {
	// build PLATFORM_ROUTES
	routesJSON, _ := json.Marshal(p.Routes)
	routesJSONB64 := base64.StdEncoding.EncodeToString(routesJSON)
	// build PLATFORM_PROJECT_ENTROPY
	entH := md5.New()
	entH.Write([]byte(entropySalt))
	entH.Write([]byte(p.ID))
	entH.Write([]byte(entropySalt))
	// build PLATFORM_RELATIONSHIPS
	relationships, _ := p.getAppRelationships(app)
	relationshipsJSON, _ := json.Marshal(relationships)
	relationshipsB64 := base64.StdEncoding.EncodeToString(relationshipsJSON)
	// build PLATFORM_VARIABLES
	appVars := p.getAppVariables(app)
	appVarsJSON, _ := json.Marshal(appVars)
	appVarsB64 := base64.StdEncoding.EncodeToString(appVarsJSON)
	// build environment vars
	envVars := map[string]string{
		"PLATFORM_DOCUMENT_ROOT":    "/app/web",
		"PLATFORM_APPLICATION":      app.BuildPlatformApplicationVar(),
		"PLATFORM_PROJECT":          p.ID,
		"PLATFORM_PROJECT_ENTROPY":  fmt.Sprintf("%x", entH.Sum(nil)),
		"PLATFORM_APPLICATION_NAME": app.Name,
		"PLATFORM_BRANCH":           "pcc",
		"PLATFORM_DIR":              appDir,
		"PLATFORM_TREE_ID":          "-",
		"PLATFORM_ENVIRONMENT":      "pcc",
		"PLATFORM_VARIABLES":        appVarsB64,
		"PLATFORM_ROUTES":           routesJSONB64,
		"PLATFORM_RELATIONSHIPS":    relationshipsB64,
	}
	// append user environment vars
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
	return envVars
}

// ShellApp - shell in to application
func (p *Project) ShellApp(app *AppDef) error {
	log.Printf("Shell in to app '%s.'", app.Name)
	// get container config
	containerConfig := p.getAppContainerConfig(app)
	return p.docker.ShellContainer(
		containerConfig.GetContainerName(),
		"web",
		[]string{"sh", "-c", "cd /app && bash"},
	)
}
