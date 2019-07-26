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
from .base import BasePlatformVariables

class DictVariables(BasePlatformVariables):
    """
    Dict variable storage handler. Store values in a dict object.
    """

    def __init__(self, config = {}):
        BasePlatformVariables.__init__(self, config)
        self._globalVars = config.get("global_vars", {})
        self._vars = dict(self.config.get("dict_vars", {}))

    def set(self, key, value):
        self._vars[str(key)] = value

    def get(self, key, default = None):
        key = str(key)
        if key in self._vars:
            return self._vars.get(key)
        if key in self._globalVars:
            return self._globalVars.get(key)
        return default

    def delete(self, key):
        self._vars.pop(str(key), None)

    def all(self):
        output = self._globalVars
        output.update(self._vars)
        return output