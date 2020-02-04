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

import os
import io
import tarfile
from .base import BaseTaskHandler
from ...exception.state_error import StateError

class RsyncTaskHandler(BaseTaskHandler):

    """
    Task handler for rsyncing files.
    """

    @classmethod
    def getType(cls):
        return "rsync"

    def run(self):
        # validate params
        self.checkParams(["from", "to"])

        # parse to path
        app, appPath = self.parseAppPath(self.params.get("to"))

        # get from container and path
        fromPath = self.params.get("from")

        # add host to known_hosts
        hostSplit = fromPath.split(":")[0].split("@")
        if len(hostSplit) > 1:
            app.runCommand(
                "ssh-keyscan %s >> ~/.ssh/known_hosts" % hostSplit[1],
                user="web"
            )

        # build command
        cmd = "rsync -a"

        # private key
        privateKey = self.params.get("private_key", "")
        if privateKey:
            app.runCommand(
                "chmod 0600 %s" % privateKey
            )
            cmd += " -e \"ssh -i %s\"" % privateKey

        # includes
        includes = self.params.get("includes", [])
        if includes:
            for include in includes:
                cmd += " --include=\"%s\"" % include
        
        # excludes
        excludes = self.params.get("excludes", [])
        if excludes:
            for excude in excludes:
                cmd += " --exclude=\"%s\"" % excude

        # add paths
        cmd += " %s %s" % (
            fromPath, appPath
        )

        # run command
        app.runCommand(cmd, user="web")