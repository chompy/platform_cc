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

package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gitlab.com/contextualcode/platform_cc/pkg/output"
	"golang.org/x/crypto/ssh"
)

const privateKeyPath = "pcc_ssh_private"
const publicKeyPath = "pcc_ssh_public"

// generateSSHKeyPair generates a public and private key pair.
func generateSSHKeyPair() (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", errors.WithStack(err)
	}
	// generate and write private key as PEM
	var privKeyBuf strings.Builder
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(&privKeyBuf, privateKeyPEM); err != nil {
		return "", "", errors.WithStack(err)
	}
	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", errors.WithStack(err)
	}
	var pubKeyBuf strings.Builder
	pubKeyBuf.Write(ssh.MarshalAuthorizedKey(pub))
	return pubKeyBuf.String(), privKeyBuf.String(), nil
}

// GenerateSSH generates SSH keypair and stores to config path.
func GenerateSSH() error {
	done := output.Duration("Generate SSH keypair.")
	pubKey, privKey, err := generateSSHKeyPair()
	if err != nil {
		return errors.WithStack(err)
	}
	// init config directory
	if err := initConfig(); err != nil {
		return errors.WithStack(err)
	}
	// save private key
	if err := ioutil.WriteFile(
		pathTo(privateKeyPath),
		[]byte(privKey),
		0600,
	); err != nil {
		return errors.WithStack(err)
	}
	// save public key
	if err := ioutil.WriteFile(
		pathTo(publicKeyPath),
		[]byte(pubKey),
		0600,
	); err != nil {
		return errors.WithStack(err)
	}
	done()
	return nil
}

func loadKey(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		// generate if not exist
		if os.IsNotExist(err) {
			if err := GenerateSSH(); err != nil {
				return nil, errors.WithStack(err)
			}
		}
		data, err = ioutil.ReadFile(path)
	}
	return data, errors.WithStack(err)
}

// PrivateKey returns the private key.
func PrivateKey() ([]byte, error) {
	return loadKey(pathTo(privateKeyPath))
}

// PrivateKeyPath returns path to the private key.
func PrivateKeyPath() string {
	// generate if not exist
	_, err := os.Stat(pathTo(privateKeyPath))
	if os.IsNotExist(err) {
		if err := GenerateSSH(); err != nil {
			return ""
		}
	}
	return pathTo(privateKeyPath)
}

// PublicKey returns the public key.
func PublicKey() ([]byte, error) {
	return loadKey(pathTo(publicKeyPath))
}
