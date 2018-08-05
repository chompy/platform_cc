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
from .var_dict import DictVariables

class JsonFileVariables(DictVariables):
    """
    JSON file variable storage handler.
    """

    """ Filename to use when storing variables. """
    JSON_STORAGE_FILENAME = ".pcc_variables.json"

    def __init__(self, config = {}):
        
        DictVariables.__init__(self, config)

        # get path to json file
        self._jsonPath = self.config.get(
            "json_path"
        )
        # not provided or path doesn't exist
        if not self._jsonPath or not os.path.isdir(os.path.dirname(self._jsonPath)):
            raise ValueError("Invalid or missing json_path provided.")

        # load
        self._loadVars()

    def _loadVars(self):
        """ Load variables from json file. """
        self._vars = {}
        if os.path.exists(self._jsonPath):
            with open(self._jsonPath, "r") as f:
                self._vars = json.load(f)

    def _saveVars(self):
        """ Save variables to json file. """
        with open(self._jsonPath, "w") as f:
            json.dump(
                self._vars,
                f,
                sort_keys=True,
                indent=4,
                separators=(',', ': ')
            )

    def set(self, key, value):
        DictVariables.set(self, key, value)
        self._saveVars()

    def delete(self, key):
        DictVariables.delete(self, key)
        self._saveVars()
