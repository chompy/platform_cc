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

	"gitlab.com/contextualcode/platform_cc/api/container"
	"gitlab.com/contextualcode/platform_cc/api/platformsh"

	"gitlab.com/contextualcode/platform_cc/api/config"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/api/def"
	"gitlab.com/contextualcode/platform_cc/api/output"
)

const pshSyncSSHCertPath = "/mnt/pcc_ssh_cert"
const pshSyncSSHKeyPath = "/mnt/pcc_ssh_key"

func (p *Project) platformSHSyncPreflight(envName string) error {
	if p.PlatformSH == nil || p.PlatformSH.ID == "" {
		return errors.WithStack(platformsh.ErrProjectNotFound)
	}
	if len(p.Apps) < 1 {
		return errors.WithStack(ErrNoApplicationFound)
	}
	if err := p.PlatformSH.FetchEnvironments(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// PlatformSHSyncVariables syncs the given platform.sh environment's variables to the local project.
func (p *Project) PlatformSHSyncVariables(envName string) error {

	done := output.Duration("Sync variables.")

	if err := p.platformSHSyncPreflight(envName); err != nil {
		return errors.WithStack(err)
	}
	env := p.PlatformSH.GetEnvironment(envName)
	if env == nil {
		return errors.Wrapf(platformsh.ErrEnvironmentNotFound, "platform.sh environment '%s' not found", envName)
	}
	vars, err := p.PlatformSH.Variables(env, p.Apps[0].Name)
	if err != nil {
		return errors.WithStack(err)
	}
	for k, v := range vars {
		if err := p.VarSet(k, v); err != nil {
			return errors.WithStack(err)
		}
	}
	pvars, err := p.PlatformSH.PlatformVariables(env, p.Apps[0].Name)
	if err != nil {
		return errors.WithStack(err)
	}
	for k, v := range pvars {
		if err := p.VarSet(k, def.InterfaceToString(v)); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := p.Save(); err != nil {
		return errors.WithStack(err)
	}
	done()
	return nil
}

// PlatformSHSyncMounts syncs the given platform.sh environment's mounts to the local project.
func (p *Project) PlatformSHSyncMounts(envName string) error {

	done := output.Duration("Sync mounts.")

	// get psh environment
	if err := p.platformSHSyncPreflight(envName); err != nil {
		return errors.WithStack(err)
	}
	env := p.PlatformSH.GetEnvironment(envName)
	if env == nil {
		return errors.Wrapf(platformsh.ErrEnvironmentNotFound, "platform.sh environment '%s' not found", envName)
	}

	// get ssh cert
	sshCert, err := p.PlatformSH.SSHCertficiate()
	if err != nil {
		return errors.WithStack(err)
	}
	certReader := bytes.NewReader(sshCert)

	// get ssh key
	sshKey, err := config.PrivateKey()
	if err != nil {
		return errors.WithStack(err)
	}
	keyReader := bytes.NewReader(sshKey)

	// set volume mount strategy to ensure mount sync works
	p.Options[OptionMountStrategy] = MountStrategyVolume
	if err := p.Save(); err != nil {
		return errors.WithStack(err)
	}

	// mount sync function
	syncMount := func(di interface{}) error {
		name := ""
		sshURL := ""
		cont := p.NewContainer(di)
		var mounts map[string]*def.AppMount
		switch d := di.(type) {
		case def.App:
			{
				name = d.Name
				sshURL = p.PlatformSH.SSHUrl(env, d.Name)
				mounts = d.Mounts
				break
			}
		case *def.AppWorker:
			{
				name = d.Name
				sshURL = p.PlatformSH.SSHWorkerUrl(env, d.ParentApp, d.Name)
				mounts = d.Mounts
				break
			}
		default:
			{
				return errors.WithStack(ErrInvalidDefinition)
			}
		}
		// upload ssh cert and key
		if err := cont.Upload(pshSyncSSHCertPath, certReader); err != nil {
			return errors.WithStack(err)
		}
		if err := cont.Upload(pshSyncSSHKeyPath, keyReader); err != nil {
			return errors.WithStack(err)
		}
		// itterate mounts and rsync
		for dest := range mounts {
			done2 := output.Duration(fmt.Sprintf("%s:%s", name, dest))
			if _, err := cont.Shell(
				"root",
				[]string{
					"ssh-agent",
					"bash",
					"-c",
					fmt.Sprintf(
						`chmod 0600 %s && chmod 0600 %s && ssh-add %s && rsync -avzh -e "ssh -i %s" %s:/app/%s/ /app/%s/`,
						pshSyncSSHKeyPath,
						pshSyncSSHCertPath,
						pshSyncSSHKeyPath,
						pshSyncSSHCertPath,
						sshURL,
						strings.Trim(dest, "/"),
						strings.Trim(dest, "/"),
					),
				},
			); err != nil {
				if !errors.Is(err, container.ErrCommandExited) {
					return errors.WithStack(err)
				}
				output.Warn("Mount sync exited with non zero code.")
			}
			done2()
		}
		// remove ssh key
		cont.Shell("root", []string{"bash", "-c", fmt.Sprintf("rm %s && rm %s", pshSyncSSHCertPath, pshSyncSSHKeyPath)})
		return nil
	}

	// itterate apps to sync mounts
	for _, app := range p.Apps {
		if err := syncMount(app); err != nil {
			return errors.WithStack(err)
		}
		// itterate workers
		if p.HasFlag(EnableWorkers) {
			for _, worker := range app.Workers {
				if err := syncMount(worker); err != nil {
					return errors.WithStack(err)
				}
			}
		}
	}
	done()
	return nil
}

// PlatformSHSyncDatabases syncs the given platform.sh environment's databases to the local project.
func (p *Project) PlatformSHSyncDatabases(envName string) error {

	done := output.Duration("Sync databases .")

	// get psh environment
	if err := p.platformSHSyncPreflight(envName); err != nil {
		return errors.WithStack(err)
	}
	env := p.PlatformSH.GetEnvironment(envName)
	if env == nil {
		return errors.Wrapf(platformsh.ErrEnvironmentNotFound, "environment '%s' not found", envName)
	}

	// fetch relationships for dump passwords
	relationships, err := p.PlatformSH.PlatformRelationships(env, p.Apps[0].Name)
	if err != nil {
		return errors.WithStack(err)
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
						return errors.WithStack(err)
					}
					done3()
					// download dump
					done3 = output.Duration("Download dump.")
					dbPath, err := p.PlatformSH.SSHDownload(
						env, p.Apps[0].Name,
						"/tmp/db.sql.gz",
					)
					if err != nil {
						return errors.WithStack(err)
					}
					// delete remote dump
					if _, err := p.PlatformSH.SSHCommand(
						env, p.Apps[0].Name,
						"rm /tmp/db.sql.gz",
					); err != nil {
						return errors.WithStack(err)
					}
					done3()
					// import dump
					done3 = output.Duration("Import dump.")
					dbOpen, err := os.Open(dbPath)
					if err != nil {
						return errors.WithStack(err)
					}
					defer func() {
						dbOpen.Close()
						os.Remove(dbPath)
					}()
					cont := p.NewContainer(service)
					if err := cont.Upload("/mnt/data/db.sql.gz", dbOpen); err != nil {
						return errors.WithStack(err)
					}
					if _, err := cont.containerHandler.ContainerShell(
						cont.Config.GetContainerName(),
						"root",
						[]string{"sh", "-c", fmt.Sprintf("zcat /mnt/data/db.sql.gz | %s && rm /mnt/data/db.sql.gz", p.GetDatabaseShellCommand(service, db))},
						nil,
					); err != nil {
						return errors.WithStack(err)
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
