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

from __future__ import absolute_import
from cleo import Command
from platform_cc.platform_sh.config import PlatformShConfig
from platform_cc.platform_sh.api import PlatformShApi
from platform_cc.platform_sh.cloner import PlatformShCloner
import os
import sys

class PlatformShLogin(Command):
    """
    Login to Platform.sh using an API token.

    platform_sh:login
        {token : Access token.}
    """

    def handle(self):
        pshConfig = PlatformShConfig()
        accessToken = PlatformShApi.getAccessToken(
            self.argument("token")
        )
        if accessToken:
            pshConfig.setAccessToken(accessToken)

class PlatformShLogout(Command):
    """
    Logout of Platform.sh and delete all stored credientials.

    platform_sh:logout
    """

    def handle(self):
        pshConfig = PlatformShConfig()
        pshConfig.setAccessToken("")

class PlatformShSetSsh(Command):
    """
    Set SSH private key to use to clone Platform.sh Git repository.

    platform_sh:set_ssh
    {--p|path=? : Path to SSH private key file. (Default=read from STDIN)}
    """

    def handle(self):
        # read private key
        path = self.option("path")
        sshPrivateKey = ""
        if path:
            with open(path, "r") as f:
                sshPrivateKey = f.read(4096)
        if not sshPrivateKey:
            sshPrivateKey = sys.stdin.read(4096)
        pshConfig = PlatformShConfig()
        pshConfig.setSshPrivateKey(sshPrivateKey)

class PlatformShClone(Command):
    """
    Clone an Platform.sh project.

    platform_sh:clone
        {project_id : Project ID.}
        {--p|path=? : Path to clone project to. (Default=current directory)}
        {--e|environment=? : Environment ID. (Default=master)}
    """

    def handle(self):
        projectId = self.argument("project_id")
        environment = self.option("environment")
        path = self.option("path")
        if not path:
            path = os.getcwd()
        if not environment: environment = "master"
        pshCloner = PlatformShCloner(
            projectId,
            environment,
            path
        )
        pshCloner.clone()
