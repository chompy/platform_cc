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
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gitlab.com/contextualcode/platform_cc/api/output"

	"github.com/helloyi/go-sshclient"

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"

	"golang.org/x/term"

	"gitlab.com/contextualcode/platform_cc/api/config"

	"github.com/ztrue/tracerr"
)

const sshAPIURL = "https://ssh.api.platform.sh/"
const sshCertificateFile = "psh_ssh_cert"

// sshCertificate defines ssh certificate storage for Platform.sh.
type sshCertificate []byte

func (k sshCertificate) data() (map[string]interface{}, error) {
	// parse cert
	key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(k))
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	// read raw data
	data := key.Marshal()
	pos := 0
	// read function
	readUint32 := func() uint32 {
		if len(data)-pos < 4 {
			return 0
		}
		dataString := data[pos : pos+4]
		pos += 4
		return binary.BigEndian.Uint32(dataString)
	}
	readUint64 := func() uint64 {
		if len(data)-pos < 8 {
			return 0
		}
		dataString := data[pos : pos+8]
		pos += 8
		return binary.BigEndian.Uint64(dataString)
	}
	readString := func() string {
		length := int(readUint32())
		if length == 0 || len(data)-pos < length {
			return ""
		}
		output := string(data[pos : pos+length])
		pos += length
		return output
	}
	out := map[string]interface{}{
		"type": key.Type,
	}
	readString() // ignore key type
	readString() // ignore nonce
	if key.Type() == "ssh-ed25519-cert-v01@openssh.com" {
		readString() //ignore ED25519 public key
	} else {
		readString() // ignore RSA exponent
		readString() // ignore RSA modulus
	}
	readUint64() // ignore serial number
	readUint32() // ignore certificate type
	out["keyId"] = readString()
	readString() // ignore valid principals
	out["validAfter"] = readUint64()
	out["validBefore"] = readUint64()
	return out, nil
}

// valid returns true if given ssh private key is still valid.
func (k sshCertificate) valid() bool {
	data, err := k.data()
	if err != nil {
		return false
	}
	validBefore := time.Unix(int64(data["validBefore"].(uint64)), 0)
	return !validBefore.IsZero() && time.Now().Before(validBefore)
}

// save stores the ssh private key.
func (k sshCertificate) save() error {
	return tracerr.Wrap(ioutil.WriteFile(
		filepath.Join(config.Path(), sshCertificateFile),
		k,
		0600,
	))
}

// fetchSSHCertficiate retrieves a SSH certificate from the Platform.sh API.
func (p *Project) fetchSSHCertficiate() error {
	p.apiURL = sshAPIURL
	defer func() {
		p.apiURL = apiURL
	}()
	// fetch pcc ssh public key
	pubKey, err := config.PublicKey()
	if err != nil {
		return tracerr.Wrap(err)
	}
	// make api request
	resp := make(map[string]interface{})
	if err := p.request(
		"ssh",
		map[string]interface{}{
			"key": string(pubKey),
		},
		&resp,
	); err != nil {
		return tracerr.Wrap(err)
	}
	// ensure key was returned
	if resp["certificate"] == nil || resp["certificate"] == "" {
		return fmt.Errorf("recieved invalid ssh certificate")
	}
	// make private key and save
	out := sshCertificate([]byte(resp["certificate"].(string)))
	if err := out.save(); err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// SSHCertficiate returns the ssh certficiate to be used with Platform.sh.
func (p *Project) SSHCertficiate() ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Join(config.Path(), sshCertificateFile))
	if err != nil {
		// fetch if cert doesn't exist
		if os.IsNotExist(err) {
			done := output.Duration("Retrieve Platform.sh SSH certificate.")
			if err := p.fetchSSHCertficiate(); err != nil {
				return nil, tracerr.Wrap(err)
			}
			data, err = ioutil.ReadFile(filepath.Join(config.Path(), sshCertificateFile))
			if err != nil {
				return nil, tracerr.Wrap(err)
			}
			done()
		} else {
			return nil, tracerr.Wrap(err)
		}
	}
	out := sshCertificate(data)
	// fetch if cert expired
	if !out.valid() {
		done := output.Duration("Refresh Platform.sh SSH certificate.")
		if err := p.fetchSSHCertficiate(); err != nil {
			return nil, tracerr.Wrap(err)
		}
		data, err = ioutil.ReadFile(filepath.Join(config.Path(), sshCertificateFile))
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		out = sshCertificate(data)
		done()
	}
	return []byte(out), nil
}

// sshClientConfig returns the SSH client config for accessing a Platform.sh environment.
func (p *Project) sshClientConfig(env *Environment, service string) (*ssh.ClientConfig, error) {
	privKey, err := config.PrivateKey()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	signer, err := ssh.ParsePrivateKey(privKey)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	cert, err := p.SSHCertficiate()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	pk, _, _, _, err := ssh.ParseAuthorizedKey(cert)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	certSigner, err := ssh.NewCertSigner(pk.(*ssh.Certificate), signer)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return &ssh.ClientConfig{
		User: p.SSHUser(env, service),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(certSigner),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}, nil
}

// openSSH opens SSH connection and returns client.
func (p *Project) openSSH(env *Environment, service string) (*ssh.Client, error) {
	sshURL := strings.Split(p.SSHUrl(env, service), "@")
	config, err := p.sshClientConfig(env, service)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	client, err := ssh.Dial("tcp", sshURL[1]+":22", config)
	if err != nil {
		return nil, tracerr.Wrap(err)
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
	// open ssh client
	client, err := p.openSSH(env, service)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer client.Close()
	// start session
	sess, err := client.NewSession()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer sess.Close()
	// run command
	out, err := sess.Output(command)
	return out, tracerr.Wrap(err)
}

// SSHDownload downloads given remote file local and returns path to downloaded file.
func (p *Project) SSHDownload(env *Environment, service string, path string) (string, error) {
	// open ssh connection
	client, err := p.openSSH(env, service)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	// open sftp connection
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	defer sftpClient.Close()
	sftpFile, err := sftpClient.Open(path)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	// open output file
	outPath := filepath.Join(os.TempDir(), service+"-"+filepath.Base(path))
	outFile, err := os.Create(outPath)
	if err != nil {
		return "", tracerr.Wrap(err)
	}
	defer outFile.Close()
	// download
	if _, err := sftpFile.WriteTo(outFile); err != nil {
		return "", tracerr.Wrap(err)
	}
	return outPath, nil
}

// SSHTerminal creates an interactive SSH terminal.
func (p *Project) SSHTerminal(env *Environment, service string) error {
	output.Info(fmt.Sprintf("SSH in to Platform.sh environment %s-%s--%s.", p.ID, env.Name, service))
	// open ssh client
	clientConfig, err := p.sshClientConfig(env, service)
	if err != nil {
		return tracerr.Wrap(err)
	}
	sshURL := strings.Split(p.SSHUrl(env, service), "@")
	client, err := sshclient.Dial(
		"tcp", sshURL[1]+":22", clientConfig,
	)
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer client.Close()
	// create interactive shell
	// make raw
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }()
	// start
	if err := client.Terminal(nil).Start(); err != nil {
		if !strings.Contains(err.Error(), "exited with status") {
			return tracerr.Wrap(err)
		}
	}
	return nil
}
