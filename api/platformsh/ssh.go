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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/melbahja/goph"
	"github.com/ztrue/tracerr"
	"golang.org/x/crypto/ssh"
)

const privateKeyPath = "~/.pcc/psh.key"
const sshKeyTitle = "Platform.CC V2 (%s)"
const sshWaitCheckIntveral = 30
const sshWaitRetryCount = 10

// makeSSHKeyPair generates a public and private key pair.
func makeSSHKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", tracerr.Wrap(err)
	}
	// generate and write private key as PEM
	var privKeyBuf strings.Builder
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(&privKeyBuf, privateKeyPEM); err != nil {
		return "", "", tracerr.Wrap(err)
	}
	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", tracerr.Wrap(err)
	}
	var pubKeyBuf strings.Builder
	pubKeyBuf.Write(ssh.MarshalAuthorizedKey(pub))
	return pubKeyBuf.String(), privKeyBuf.String(), nil
}

// hasSSHKey returns true if local private key is found.
func hasSSHKey() bool {
	_, err := os.Stat(expandPath(privateKeyPath))
	return !os.IsNotExist(err)
}

// PrivateKey returns the private key.
func PrivateKey() ([]byte, error) {
	if !hasSSHKey() {
		return nil, tracerr.Errorf("private key not found")
	}
	return ioutil.ReadFile(expandPath(privateKeyPath))
}

// storeSSHKey generates a new ssh key, sends it to platform.sh, and stores the private key locally.
func (p *Project) storeSSHKey() error {
	// generate keys
	done := output.Duration("Generate Platform.sh SSH key.")
	pubKey, privKey, err := makeSSHKeyPair()
	if err != nil {
		return tracerr.Wrap(err)
	}
	output.LogDebug("Platform.sh SSH key generated.", map[string]string{"public": pubKey, "private": privKey})
	done()
	// upload public key to platform.sh
	done = output.Duration("Upload public key to Platform.sh.")
	res := make(map[string]interface{})
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "-"
	}
	if err := p.request(
		"ssh_keys",
		map[string]interface{}{
			"value": pubKey,
			"title": fmt.Sprintf(sshKeyTitle, hostname),
		},
		&res,
	); err != nil {
		return tracerr.Wrap(err)
	}
	if res["value"] == nil || res["title"] == nil {
		return tracerr.Errorf("recieved malformed response when sending ssh key")
	}
	done()
	// save private key
	done = output.Duration("Save private key.")
	if err := ioutil.WriteFile(
		expandPath(privateKeyPath),
		[]byte(privKey),
		0600,
	); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// waitForSSH waits for a newly generated Platform.sh SSH key to be accepted.
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
	// generate ssh key if not exist
	if !hasSSHKey() {
		if err := p.storeSSHKey(); err != nil {
			return nil, tracerr.Wrap(err)
		}
		if err := p.waitForSSH(env, service); err != nil {
			return nil, tracerr.Wrap(err)
		}
	}
	// parse private key
	auth, err := goph.Key(expandPath(privateKeyPath), "")
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	// open ssh connection
	client, err := goph.NewUnknown(
		p.SSHUser(env, service),
		env.EdgeHostname,
		auth,
	)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return client, nil
}

// SSHUrl returns the SSH url for the environment.
func (p Project) SSHUrl(env *Environment, service string) string {
	return fmt.Sprintf(
		"%s@%s",
		p.SSHUser(env, service),
		env.EdgeHostname,
	)
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
