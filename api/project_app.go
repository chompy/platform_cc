package api

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gitlab.com/contextualcode/platform_cc/def"
)

// GetAppContainerConfig gets container configuration for an app.
func (p *Project) GetAppContainerConfig(app *def.App) DockerContainerConfig {
	return p.GetServiceContainerConfig(app)
}

// startApp - start an app
func (p *Project) startApp(app *def.App) error {
	return p.startService(app)
}

// openApp - make the application available
func (p *Project) openApp(app *def.App) error {
	return p.openService(app)
}

// BuildApp builds the app.
func (p *Project) BuildApp(app *def.App) error {
	log.Printf("Building app '%s.'", app.Name)
	// get container config
	containerConfig := p.GetAppContainerConfig(app)
	// build flavor
	buildFlavorComposer := strings.ToLower(app.Build.Flavor) == "composer"
	// upload build script
	buildScriptReader := strings.NewReader(
		fmt.Sprintf(appBuildScript, strings.Title(strconv.FormatBool(buildFlavorComposer))),
	)
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
		"source_dir":  def.AppDir,
		"output_dir":  def.AppDir,
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

// DeployApp - deploy the application (run deploy hooks)
func (p *Project) DeployApp(app *def.App) error {
	log.Printf("Deploying app '%s.'", app.Name)
	// get container config
	containerConfig := p.GetAppContainerConfig(app)
	// run command
	if err := p.docker.RunContainerCommand(
		containerConfig.GetContainerName(),
		"root",
		[]string{"sh", "-c", appDeployCmd},
		os.Stdout,
	); err != nil {
		return err
	}
	return nil
}

// getAppRelationships - generate relationships variable for app
func (p *Project) getAppRelationships(def interface{}) (map[string][]map[string]interface{}, error) {
	return p.getServiceRelationships(def)
}

// getAppVariables - get variables to inject in to container
func (p *Project) getAppVariables(app *def.App) map[string]string {
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
func (p *Project) getAppEnvironmentVariables(app *def.App) map[string]string {
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
		"PLATFORM_DIR":              def.AppDir,
		"PLATFORM_APP_DIR":          def.AppDir,
		"PLATFORM_TREE_ID":          "-",
		"PLATFORM_ENVIRONMENT":      "pcc",
		"PLATFORM_VARIABLES":        appVarsB64,
		"PLATFORM_ROUTES":           routesJSONB64,
		"PLATFORM_RELATIONSHIPS":    relationshipsB64,
		"PLATFORM_CACHE_DIR":        "/tmp/cache",
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
func (p *Project) ShellApp(app *def.App) error {
	return p.ShellService(app, []string{"bash"})
}
