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

from platform_cc.config import PlatformConfig

class PlatformShConfig:

    """ Platform.sh configuration. """

    CONFIG_ACCESS_TOKEN_KEY = "platform_sh_access_token"
    CONFIG_SSH_PUBLIC_KEY = "platform_sh_ssh_public_key"
    CONFIG_SSH_PRIVATE_KEY = "platform_sh_ssh_private_key"

    def __init__(self, globalConfig = None):
        if not globalConfig:
            globalConfig = PlatformConfig()
        self.globalConfig = globalConfig

    def getAccessToken(self):
        """ Retrieve access token. """
        return self.globalConfig.get(self.CONFIG_ACCESS_TOKEN_KEY, "")

    def setAccessToken(self, accessToken = ""):
        """ Set access token. """
        self.globalConfig.set(self.CONFIG_ACCESS_TOKEN_KEY, accessToken)
        self.globalConfig.save()

    def getSshPublicKey(self):
        """ Get SSH public key. """
        return self.globalConfig.get(self.CONFIG_SSH_PUBLIC_KEY, "") 

    def setSshPublicKey(self, sshKey = ""):
        """ Set SSH public key. """
        self.globalConfig.set(self.CONFIG_SSH_PUBLIC_KEY, sshKey)
        self.globalConfig.save()

    def getSshPrivateKey(self):
        """ Get SSH private key. """
        return self.globalConfig.get(self.CONFIG_SSH_PRIVATE_KEY, "")
    
    def setSshPrivateKey(self, sshKey = ""):
        """ Set SSH private key. """
        self.globalConfig.set(self.CONFIG_SSH_PRIVATE_KEY, sshKey)
        self.globalConfig.save()