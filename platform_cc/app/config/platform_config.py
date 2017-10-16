import os
import time
import hashlib
import json

class PlatformConfig:

    """ Base config class. """

    ENTROPY_SALT = "M0igz7x2nr0cCXitjPvbF5eQhsf01F#!"

    def __init__(self, projectHash):
        self.projectHash = projectHash
        self._config = {}

    def getName(self):
        return "default"

    def getType(self):
        return self._config.get("type", None)

    def getBuildFlavor(self):
        return None

    def getMounts(self):
        return {}

    def getDockerImage(self):
        return None

    def getVariables(self):
        return {}

    def getEntropy(self):
        return hashlib.sha256(
            self.projectHash + self.ENTROPY_SALT
        ).hexdigest()
