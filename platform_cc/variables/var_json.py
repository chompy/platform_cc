import os
import json
from .base import BasePlatformVariables

class JsonVariables(BasePlatformVariables):
    """
    JSON Variable storage handler.
    """

    """ Filename to use when storing variables. """
    JSON_STORAGE_FILENAME = ".pcc_variables.json"

    def __init__(self, projectPath, projectConfig = None):
        BasePlatformVariables.__init__(self, projectPath, projectConfig)
        self.JSON_STORAGE_FILENAME = projectConfig.get(
            "variables_json_filename",
            self.JSON_STORAGE_FILENAME
        )
        self._jsonPath = os.path.join(
            projectPath,
            self.JSON_STORAGE_FILENAME
        )
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
        self._vars[str(key)] = str(value)
        self._saveVars()

    def get(self, key, default = None):
        return self._vars.get(str(key), default)

    def delete(self, key):
        self._vars.pop(str(key), None)
        self._saveVars()

    def all(self):
        return self._vars.copy()