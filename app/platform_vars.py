import os
import base64
from platform_docker import PlatformDocker
from config.platform_vars_config import PlatformVarsConfig
from app.platform_utils import log_stdout

class PlatformVars:

    """ Handles storage and retrival of user defined variables. """

    DOCKER_CMD = "sleep 1"
    LIST_ALL_DELIMITER = "-------!!!--------"

    def __init__(self, projectHash):
        self.config = PlatformVarsConfig(projectHash)
        self.docker = PlatformDocker(
            self.config,
            "app_vars",
            PlatformVarsConfig.VARS_DOCKER_IMAGE
        )

    def _sanitizeKey(self, key):
        return key.replace(" ", "_").strip().lower()

    def set(self, key, value):
        log_stdout("Set variable '%s'..." % key)
        self.docker.start(self.DOCKER_CMD)
        self.docker.getProvisioner().copyStringToFile(
            str(value),
            "/data/%s" % self._sanitizeKey(key)
        )
        self.docker.stop()

    def get(self, key):
        log_stdout("Get variable '%s'..." % key)
        self.docker.start(self.DOCKER_CMD)
        results = self.docker.getContainer().exec_run(
            ["sh", "-c", "cat /data/%s 2> /dev/null" % self._sanitizeKey(key)]
        )        
        self.docker.stop()
        return results

    def all(self):
        log_stdout("Fetch all variables...")
        self.docker.start(self.DOCKER_CMD)
        cmd = "find /data -type f -exec sh -c 'printf \"$0%s\" && cat $0 | base64 && printf \"%s\"' {} \\;" % (
            self.LIST_ALL_DELIMITER,
            self.LIST_ALL_DELIMITER
        )
        results = self.docker.getContainer().exec_run(
            ["sh", "-c", cmd]
        )
        self.docker.stop()
        results = results.strip().strip(self.LIST_ALL_DELIMITER).replace("\n", "").split(self.LIST_ALL_DELIMITER)
        output = {}
        i = 0;
        while i < len(results):
            key = results[i][6:]
            output[key] = base64.b64decode(results[i + 1])
            i += 2
        return output