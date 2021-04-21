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

package platformsh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/helloyi/go-sshclient"
	"golang.org/x/term"

	"gitlab.com/contextualcode/platform_cc/api/config"
	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/melbahja/goph"
	"github.com/ztrue/tracerr"
)

const sshKeyTitle = "Platform.CC V2 (%s)"
const sshWaitCheckIntveral = 30
const sshWaitRetryCount = 20

// storeSSHKey generates a new ssh key, sends it to platform.sh, and stores the private key locally.
func (p *Project) storeSSHKey() error {
	// load public key, generate if not exist
	pubKey, err := config.PublicKey()
	if err != nil {
		if os.IsNotExist(err) {
			if err := config.GenerateSSH(); err != nil {
				return tracerr.Wrap(err)
			}
			pubKey, err = config.PublicKey()
			if err != nil {
				return tracerr.Wrap(err)
			}
		} else {
			return tracerr.Wrap(err)
		}
	}
	// upload public key to platform.sh
	done := output.Duration("Upload public key to Platform.sh.")
	res := make(map[string]interface{})
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "-"
	}
	if err := p.request(
		"ssh_keys",
		map[string]interface{}{
			"value": string(pubKey),
			"title": fmt.Sprintf(sshKeyTitle, hostname),
		},
		&res,
	); err != nil {
		// if error returns 'fingerprint already in use' then ssh key was previously uploaded
		// and this should not be treated as an error
		if res["title"] != nil && strings.Contains(res["title"].(string), "is already in use") {
			done()
			return nil
		}
		return tracerr.Wrap(err)
	}
	if res["value"] == nil || res["title"] == nil {
		return tracerr.Errorf("recieved malformed response when sending ssh key")
	}
	done()
	return nil
}

// waitForSSH waits for a newly uploaded Platform.sh SSH key to be accepted.
func (p *Project) waitForSSH(env *Environment, service string) error {
	if env == nil {
		return tracerr.Errorf("invalid environment")
	}
	done := output.Duration("Waiting for Platform.sh to accept new key. (This can take a while.)")
	i := 0
	for i = 0; i < sshWaitRetryCount; i++ {
		time.Sleep(time.Second * sshWaitCheckIntveral)
		if _, err := p.SSHCommand(env, service, "true"); err == nil {
			break
		}
	}
	if i >= sshWaitRetryCount {
		return tracerr.Errorf("timed out")
	}
	done()
	return nil
}

// openSSH opens SSH connection and returns client.
func (p *Project) openSSH(env *Environment, service string) (*goph.Client, error) {
	// load goph key auth, create key if not exist
	auth, err := config.KeyAuth()
	if err != nil {
		if os.IsNotExist(err) {
			if err := config.GenerateSSH(); err != nil {
				return nil, tracerr.Wrap(err)
			}
			auth, err = config.KeyAuth()
			if err != nil {
				return nil, tracerr.Wrap(err)
			}
		} else {
			return nil, tracerr.Wrap(err)
		}
	}
	// open ssh connection
	sshURL := strings.Split(p.SSHUrl(env, service), "@")
	client, err := goph.NewUnknown(
		sshURL[0],
		sshURL[1],
		auth,
	)
	if err != nil {
		// handshake failed, upload ssh key to platform.sh api
		if strings.Contains(err.Error(), "handshake failed") && !env.hasSSH {
			if err := p.storeSSHKey(); err != nil {
				return nil, tracerr.Wrap(err)
			}
			env.hasSSH = true
			if err := p.waitForSSH(env, service); err != nil {
				return nil, tracerr.Wrap(err)
			}
			// try again
			client, err = goph.NewUnknown(
				p.SSHUser(env, service),
				env.EdgeHostname,
				auth,
			)
			if err != nil {
				return nil, tracerr.Wrap(err)
			}
		} else {
			return nil, tracerr.Wrap(err)
		}
	}
	return client, nil
}

// SSHUrl returns the SSH url for the environment.
func (p Project) SSHUrl(env *Environment, service string) string {
	urlSplit := strings.Split(strings.TrimPrefix(env.Links.SSH.HREF, "ssh://"), "@")
	return fmt.Sprintf("%s--%s@%s", urlSplit[0], service, urlSplit[1])
}

// SSHUser returns the SSH username for given environment and service.
func (p Project) SSHUser(env *Environment, service string) string {
	return fmt.Sprintf(
		"%s-%s--%s",
		p.ID,
		env.MachineName,
		service,
	)
}

// SSHCommand sends a command to given Platform.sh environment over SSH and returns the output.
func (p *Project) SSHCommand(env *Environment, service string, command string) ([]byte, error) {
	// open ssh connection
	client, err := p.openSSH(env, service)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	// send command, return results
	defer client.Close()
	out, err := client.Run(command)
	return out, tracerr.Wrap(err)
}

// SSHDownload downloads given remote file local and returns path to downloaded file.
func (p *Project) SSHDownload(env *Environment, service string, path string) (string, error) {
	// open ssh connection
	client, err := p.openSSH(env, service)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	// prepare download
	outPath := filepath.Join(os.TempDir(), service+"-"+filepath.Base(path))
	if err := client.Download(path, outPath); err != nil {
		return "", tracerr.Wrap(err)
	}
	return outPath, nil
}

// SSHTerminal creates an interactive SSH terminal.
func (p *Project) SSHTerminal(env *Environment, service string) error {
	// open ssh connection
	sshURL := strings.Split(p.SSHUrl(env, service), "@")
	output.Info(fmt.Sprintf("Open SSH terminal to %s.", sshURL[1]))
	client, err := sshclient.DialWithKey(
		sshURL[1]+":22",
		sshURL[0],
		config.PrivateKeyPath(),
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	// create interactive shell
	// make raw
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()
	// create terminal
	if err := client.Terminal(&sshclient.TerminalConfig{
		Term: "xterm",
	}).Start(); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}
