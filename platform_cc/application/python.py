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

import io
import os
import json
from .base import BasePlatformApplication
from ..exception.container_command_error import ContainerCommandError
from ..core.version import PCC_VERSION

class PythonApplication(BasePlatformApplication):
    """
    Handler for Python applications.
    """


    """ Mapping for application type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "python:3.7"         : "chompy/platform_cc:%s-python37" % PCC_VERSION,  
    }

    """ Default user id to assign for user 'web' """
    DEFAULT_WEB_USER_ID = 1000

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def build(self):
        self.prebuild()
        output = ""
        # add web user
        self.logger.info(
            "Add and configure 'web' user."
        )
        output += self.runCommand(
            """
            useradd -d /app -m -p secret~ --uid %s web
            usermod -a -G staff web
            mkdir -p /var/lib/gems
            chown -R web:web /var/lib/gems
            chown -R root:staff /usr/bin
            chmod -R g+rw /usr/bin
            """ % (
                self.project.get("config", {}).get("web_user_id", self.DEFAULT_WEB_USER_ID)
            )
        )
        output += self.runCommand(
            "usermod -u %s web" % (
                self.project.get("config", {}).get("web_user_id", self.DEFAULT_WEB_USER_ID)
            )            
        )
        # install ssh key + known_hosts
        self.installSsh()
        output += self.runCommand(
            "chown -f -R web /app/.ssh"
        )
        self.logger.info(
            "Setup/fix user permission."
        )
        try:
            output += self.runCommand(
                """
                chown -f -R web %s
                chown -f -R web %s
                """ % (self.STORAGE_DIRECTORY, self.APPLICATION_DIRECTORY)
            )
        except ContainerCommandError:
            pass
        # build hooks
        self.logger.info(
            "Run build hooks."
        )
        try:
            output += self.runCommand(
                self.config.get("hooks", {}).get("build", ""),
                "root"
            )
        # allow build hooks to fail...for now
        except ContainerCommandError:
            pass
        # clean up
        self.logger.info(
            "Clean up."
        )
        output += self.runCommand(
            """
            apt-get clean
            """
        )
        # commit container
        self.logger.info(
            "Commit container."
        )
        self.commit()
        return output

    def start(self, requireServices = True):
        BasePlatformApplication.start(self, requireServices)
        container = self.getContainer()
        if not container: return
        self.startServices()