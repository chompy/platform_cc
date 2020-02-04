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
from ..platform_sh.config import PlatformShConfig
from ..platform_sh.api import PlatformShApi
from ..platform_sh.cloner import PlatformShCloner
from .common import getProject
from ..platform_sh.exception.config_error import PlatformShConfigError
import os
import sys
import time

class PlatformShLogin(Command):
    """
    Login to Platform.sh using an API token.

    platform_sh:login
        {token : API token. (Generate your at https://accounts.platform.sh/user/api-tokens.) }
    """

    def handle(self):
        pshConfig = PlatformShConfig()
        pshConfig.setApiToken(self.argument("token"))
        pshConfig.setAccessToken("")

class PlatformShLogout(Command):
    """
    Logout of Platform.sh and delete all stored credientials.

    platform_sh:logout
    """

    def handle(self):
        pshConfig = PlatformShConfig()
        pshConfig.setAccessToken("")
        pshConfig.setApiToken("")

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
        {--skip-mount-sync : Skip syncing mount directories.}
        {--skip-service-sync : Skip syncing service assets.}
    """

    def handle(self):
        projectId = self.argument("project_id")
        environment = self.option("environment")
        # warn user about clone
        self.line(
            "<question>!!! Clone of Platform.sh project '%s:%s' will commence in 5 seconds. Press CTRL+C to cancel. !!!</question>" % (
                projectId,
                "master" if not environment else environment
            )
        )
        time.sleep(5)
        path = self.option("path")
        if not path:
            path = os.getcwd()
        if not environment: environment = "master"
        pshCloner = PlatformShCloner(
            projectId,
            environment,
            path
        )
        pshCloner.clone(
            skipMountSync=self.option("skip-mount-sync"),
            skipServiceSync=self.option("skip-service-sync")
        )

class PlatformShSync(Command):
    """
    Sync an Platform.sh project.

    platform_sh:sync
        {--p|path=? : Path to project root. (Default=current directory)}
        {--e|environment=? : Environment ID. (Default=master)}
        {--skip-var-sync : Skip syncing project variables.}
        {--skip-mount-sync : Skip syncing mount directories.}
        {--skip-service-sync : Skip syncing service assets.}
    """

    def handle(self):
        project = getProject(self)
        pshProjectId = project.variables.get("env:PSH_PROJECT_ID")
        if not pshProjectId:
            raise PlatformShConfigError(
                "Project '%s' does not have a Platform.sh project id set. You can set it with 'var:set env:PSH_PROJECT_ID <project_id>'." % (
                    project.getShortUid()
                )
            )
        environment = self.option("environment")
        if not environment: environment = "master"
        # warn user about sync
        self.line(
            "<question>!!! Sync of project '%s' with Platform.sh (%s:%s) will commence in 5 seconds. Press CTRL+C to cancel. !!!</question>" % (
                project.getShortUid(),
                pshProjectId,
                environment
            )
        )
        time.sleep(5)
        # sync using cloner class
        pshCloner = PlatformShCloner(
            pshProjectId,
            environment,
            project.path
        )
        
        try:
            pshCloner.start()
            if not self.option("skip-var-sync"):
                pshCloner.syncVars(project)
            project.start()
            if not self.option("skip-mount-sync"):
                pshCloner.syncMounts(project)
            if not self.option("skip-service-sync"):
                pshCloner.syncServices(project)
        except Exception as e:
            pshCloner.logger.error("An error occured, stopping container...")
            project.stop()
            pshCloner.stop()
            raise e
        project.stop()
        pshCloner.stop()
