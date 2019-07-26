"""
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
"""

import cryptography
from cryptography.hazmat.primitives import serialization as crypto_serialization
from platform_cc.config import PlatformConfig

class PlatformShConfig:

    """ Platform.sh configuration. """

    CONFIG_API_TOKEN_KEY = "platform_sh_api_token"
    CONFIG_ACCESS_TOKEN_KEY = "platform_sh_access_token"
    CONFIG_SSH_PUBLIC_KEY = "platform_sh_ssh_public_key"
    CONFIG_SSH_PRIVATE_KEY = "platform_sh_ssh_private_key"
    CONFIG_API_TOKEN_KEY = "platform_sh_api_token"

    def __init__(self, globalConfig = None):
        if not globalConfig:
            globalConfig = PlatformConfig()
        self.globalConfig = globalConfig

    def getApiToken(self):
        """ Retrieve API token. """
        return self.globalConfig.get(self.CONFIG_API_TOKEN_KEY, "")

    def setApiToken(self, apiToken = ""):
        """ Set API token. """
        self.globalConfig.set(self.CONFIG_API_TOKEN_KEY, apiToken)

    def getAccessToken(self):
        """ Retrieve access token. """
        return self.globalConfig.get(self.CONFIG_ACCESS_TOKEN_KEY, "")

    def setAccessToken(self, accessToken = ""):
        """ Set access token. """
        self.globalConfig.set(self.CONFIG_ACCESS_TOKEN_KEY, accessToken)

    def getSshPublicKey(self):
        """ Get SSH public key. """
        return self.globalConfig.get(self.CONFIG_SSH_PUBLIC_KEY, "") 

    def getSshPrivateKey(self):
        """ Get SSH private key. """
        return self.globalConfig.get(self.CONFIG_SSH_PRIVATE_KEY, "")
    
    def setSshPrivateKey(self, sshKey = ""):
        """ Set SSH private key. """
        # unset
        if not sshKey:
            self.globalConfig.set(self.CONFIG_SSH_PUBLIC_KEY, "")
            self.globalConfig.set(self.CONFIG_SSH_PRIVATE_KEY, "")
            return
        # load private key to verify + get public key
        privateKey = cryptography.hazmat.primitives.serialization.load_pem_private_key(
            str.encode(sshKey),
            password=None,
            backend=cryptography.hazmat.backends.default_backend()
        )
        privateKeyBytes = privateKey.private_bytes(
            crypto_serialization.Encoding.PEM,
            crypto_serialization.PrivateFormat.PKCS8,
            crypto_serialization.NoEncryption()
        )
        publicKeyBytes = privateKey.public_key().public_bytes(
            crypto_serialization.Encoding.OpenSSH,
            crypto_serialization.PublicFormat.OpenSSH
        )
        self.globalConfig.set(self.CONFIG_SSH_PUBLIC_KEY, publicKeyBytes.decode("utf-8"))
        self.globalConfig.set(self.CONFIG_SSH_PRIVATE_KEY, privateKeyBytes.decode("utf-8"))