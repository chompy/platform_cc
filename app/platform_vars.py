import os
import base64
import json
from platform_docker import PlatformDocker
from config.platform_vars_config import PlatformVarsConfig
from app.platform_utils import log_stdout

class PlatformVars:

    """ Handles storage and retrival of user defined variables. """

    DOCKER_CMD = "sleep 1"
    VAR_PATH = "/data/vars"

    def __init__(self, projectHash):
        self.config = PlatformVarsConfig(projectHash)
        self.docker = PlatformDocker(
            self.config,
            "app_vars",
            PlatformVarsConfig.VARS_DOCKER_IMAGE
        )
        self.docker.logIndent = -1

    def _sanitizeKey(self, key):
        return key.replace(" ", "_").strip().lower()

    def set(self, key, value):
        allVars = self.all()
        allVars[self._sanitizeKey(key)] = value
        self.docker.start(self.DOCKER_CMD)
        self.docker.getProvisioner().copyStringToFile(
            json.dumps(allVars),
            self.VAR_PATH
        )
        self.docker.stop()

    def get(self, key):
        allVars = self.all()
        return allVars.get(
            self._sanitizeKey(key),
            None
        )

    def all(self):
        self.docker.start(self.DOCKER_CMD)
        results = self.docker.getProvisioner().fetchFile(self.VAR_PATH)
        self.docker.stop()
        if not results: return {}
        return json.loads(results)