import os
import time
import hashlib
import json

class PlatformConfig:

    """ Base config class. """

    PLATFORM_LOCAL_DATA_PATH = ".platform/.pcclocal"

    def __init__(self, projectHash, appPath = ""):
        self.projectHash = projectHash
        self.appPath = appPath
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

    def getDataPath(self):
        return os.path.join(
            self.appPath,
            self.PLATFORM_LOCAL_DATA_PATH
        )

    def getEntropy(self):
        entropyPath = os.path.join(self.getDataPath(), ".entropy")
        if not os.path.exists(entropyPath):
            entropy = hashlib.sha256(
                self.appPath + json.dumps(self._config) + self.PLATFORM_LOCAL_DATA_PATH + str(time.time())
            ).hexdigest()
            with open(entropyPath, "w") as f:
                f.write(entropy)
            return entropy
        with open(entropyPath, "r") as f:
            return f.read()