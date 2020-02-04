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
import json
from ..variables.var_json_file import JsonFileVariables

class PlatformConfig(JsonFileVariables):

    """ Global configuration. """

    GLOBAL_PROJECT_VAR_PREFIX = "global_project_var:"

    def __init__(self, filename = None):
        if not filename:
            filename = os.path.join(os.path.expanduser("~"), ".platform_cc_config.json")
        JsonFileVariables.__init__(
            self, {
                "json_path" : filename
            }
        )

    def getGlobalProjectVars(self):
        """ Retrieve all global project vars. """
        allVars = self.all()
        output = {}
        for key in allVars:
            if key.startswith(self.GLOBAL_PROJECT_VAR_PREFIX):
                output[key[len(self.GLOBAL_PROJECT_VAR_PREFIX):]] = allVars[key]
        return output