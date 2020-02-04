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

class ShellTaskHandler(BaseTaskHandler):

    """
    Task handler for running shell commands.
    """

    @classmethod
    def getType(cls):
        return "shell"

    def run(self):
        # validate params
        self.checkParams(["to", "command"])
        # get 'to' application
        toApp = self.project.getApplication(self.params.get("to"))
        # get user to run as
        user = self.params.get("user", "web")
        # run command
        toApp.runCommand(self.params.get("command"), user)
