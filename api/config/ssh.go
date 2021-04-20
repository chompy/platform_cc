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
	"strings"

	"github.com/melbahja/goph"
	"github.com/ztrue/tracerr"
	"gitlab.com/contextualcode/platform_cc/api/output"
	"golang.org/x/crypto/ssh"
)

const privateKeyPath = "id_rsa"
const publicKeyPath = "id_rsa.pub"

// generateSSHKeyPair generates a public and private key pair.
func generateSSHKeyPair() (string, string, error) {
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

// GenerateSSH generates SSH keypair and stores to config path.
func GenerateSSH() error {
	done := output.Duration("Generate SSH keypair.")
	pubKey, privKey, err := generateSSHKeyPair()
	if err != nil {
		return tracerr.Wrap(err)
	}
	// init config directory
	if err := initConfig(); err != nil {
		return tracerr.Wrap(err)
	}
	// save private key
	if err := ioutil.WriteFile(
		pathTo(privateKeyPath),
		[]byte(privKey),
		0600,
	); err != nil {
		return tracerr.Wrap(err)
	}
	// save public key
	if err := ioutil.WriteFile(
		pathTo(publicKeyPath),
		[]byte(pubKey),
		0600,
	); err != nil {
		return tracerr.Wrap(err)
	}
	done()
	return nil
}

// PrivateKey returns the private key.
func PrivateKey() ([]byte, error) {
	return ioutil.ReadFile(pathTo(privateKeyPath))
}

// PublicKey returns the public key.
func PublicKey() ([]byte, error) {
	return ioutil.ReadFile(pathTo(publicKeyPath))
}

// KeyAuth returns goph auth for private key.
func KeyAuth() (goph.Auth, error) {
	return goph.Key(pathTo(privateKeyPath), "")
}
