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
	"bytes"
	"fmt"
	"os"
	"strings"

	"gitlab.com/contextualcode/platform_cc/api/config"

	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

func (p *Project) platformSHSyncPreflight(envName string) error {
	if p.PlatformSH == nil || p.PlatformSH.ID == "" {
		return tracerr.Errorf("platform.sh project not found")
	}
	if len(p.Apps) < 1 {
		return tracerr.Errorf("project should have at least one application")
	}
	if err := p.PlatformSH.FetchEnvironments(); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// PlatformSHSyncVariables syncs the given platform.sh environment's variables to the local project.
func (p *Project) PlatformSHSyncVariables(envName string) error {

	done := output.Duration(fmt.Sprintf("Sync variables from Platform.sh '%s' environment.", envName))

	if err := p.platformSHSyncPreflight(envName); err != nil {
		return tracerr.Wrap(err)
	}
	env := p.PlatformSH.GetEnvironment(envName)
	if env == nil {
		return tracerr.Errorf("environment '%s' not found", envName)
	}
	vars, err := p.PlatformSH.Variables(env)
	if err != nil {
		return tracerr.Wrap(err)
	}
	for k, v := range vars {
		if err := p.VarSet(k, v); err != nil {
			return tracerr.Wrap(err)
		}
	}
	pvars, err := p.PlatformSH.PlatformVariables(env, p.Apps[0].Name)
	if err != nil {
		return tracerr.Wrap(err)
	}
	for k, v := range pvars {
		if err := p.VarSet(k, def.InterfaceToString(v)); err != nil {
			return tracerr.Wrap(err)
		}
	}

	if err := p.Save(); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// PlatformSHSyncMounts syncs the given platform.sh environment's mounts to the local project.
func (p *Project) PlatformSHSyncMounts(envName string) error {

	done := output.Duration(fmt.Sprintf("Sync mounts from Platform.sh '%s' environment.", envName))

	// get psh environment
	if err := p.platformSHSyncPreflight(envName); err != nil {
		return tracerr.Wrap(err)
	}
	env := p.PlatformSH.GetEnvironment(envName)
	if env == nil {
		return tracerr.Errorf("environment '%s' not found", envName)
	}

	// get private key
	privKey, err := config.PrivateKey()
	if err != nil {
		return tracerr.Wrap(err)
	}
	privKeyReader := bytes.NewReader(privKey)

	// set volume mount strategy to ensure mount sync works
	p.Options[OptionMountStrategy] = MountStrategyVolume
	if err := p.Save(); err != nil {
		return tracerr.Wrap(err)
	}

	// itterate apps to sync mounts
	for _, app := range p.Apps {
		sshURL := p.PlatformSH.SSHUrl(env, app.Name)
		cont := p.NewContainer(app)
		if err := cont.Upload("/mnt/priv.key", privKeyReader); err != nil {
			return tracerr.Wrap(err)
		}
		for dest := range app.Mounts {
			done2 := output.Duration(fmt.Sprintf("%s:%s", app.Name, dest))
			if err := cont.Shell(
				"root",
				[]string{"sh", "-c",
					fmt.Sprintf(
						`chmod 0600 /mnt/priv.key && rsync -avzh -e "ssh -i /mnt/priv.key" %s:/app/%s /app/%s`,
						sshURL,
						strings.TrimLeft(dest, "/"),
						strings.TrimLeft(dest, "/"),
					),
				},
			); err != nil {
				return tracerr.Wrap(err)
			}
			done2()
		}
		cont.Shell("root", []string{"rm", "/mnt/priv.key"})
		// TODO workers
	}
	done()
	return nil
}

// PlatformSHSyncDatabases syncs the given platform.sh environment's databases to the local project.
func (p *Project) PlatformSHSyncDatabases(envName string) error {

	done := output.Duration(fmt.Sprintf("Sync databases from Platform.sh '%s' environment.", envName))

	// get psh environment
	if err := p.platformSHSyncPreflight(envName); err != nil {
		return tracerr.Wrap(err)
	}
	env := p.PlatformSH.GetEnvironment(envName)
	if env == nil {
		return tracerr.Errorf("environment '%s' not found", envName)
	}

	// fetch relationships for dump passwords
	relationships, err := p.PlatformSH.PlatformRelationships(env, p.Apps[0].Name)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// itterate services to find database services
	for _, service := range p.Services {
		for _, dbType := range GetDatabaseTypeNames() {
			if dbType == service.GetTypeName() {
				// itterate databases
				for _, dbIf := range service.Configuration["schemas"].([]interface{}) {
					db := dbIf.(string)
					done2 := output.Duration(fmt.Sprintf("%s:%s", service.Name, db))
					// create dump
					done3 := output.Duration("Create dump.")
					if _, err := p.PlatformSH.SSHCommand(
						env, p.Apps[0].Name,
						p.GetPlatformSHDatabaseDumpCommand(service, db, relationships)+" | gzip > /tmp/db.sql.gz",
					); err != nil {
						return tracerr.Wrap(err)
					}
					done3()
					// download dump
					done3 = output.Duration("Download dump.")
					dbPath, err := p.PlatformSH.SSHDownload(
						env, p.Apps[0].Name,
						"/tmp/db.sql.gz",
					)
					if err != nil {
						return tracerr.Wrap(err)
					}
					// delete remote dump
					if _, err := p.PlatformSH.SSHCommand(
						env, p.Apps[0].Name,
						"rm /tmp/db.sql.gz",
					); err != nil {
						return tracerr.Wrap(err)
					}
					done3()
					// import dump
					done3 = output.Duration("Import dump.")
					dbOpen, err := os.Open(dbPath)
					if err != nil {
						return tracerr.Wrap(err)
					}
					defer func() {
						dbOpen.Close()
						os.Remove(dbPath)
					}()
					cont := p.NewContainer(service)
					if err := cont.Upload("/mnt/data/db.sql.gz", dbOpen); err != nil {
						return tracerr.Wrap(err)
					}
					if err := cont.containerHandler.ContainerShell(
						cont.Config.GetContainerName(),
						"root",
						[]string{"sh", "-c", fmt.Sprintf("zcat /mnt/data/db.sql.gz | %s && rm /mnt/data/db.sql.gz", p.GetDatabaseShellCommand(service, db))},
						nil,
					); err != nil {
						return tracerr.Wrap(err)
					}
					done3()
					done2()
				}
				break
			}
		}
	}

	done()
	return nil

}
