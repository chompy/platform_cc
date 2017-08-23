import os
import yaml
import hashlib
import time
from platform_config import PlatformConfig

class PlatformAppConfig(PlatformConfig):

    """ Provide configuration for application. """

    PLATFORM_DOCKER_IMAGES = {
        "php:5.4":          "php:5.4-fpm",
        "php:5.6":          "php:5.6-fpm"
    }

    PLATFORM_FILENAME = ".platform.app.yaml"

    def __init__(self, projectHash, appPath = ""):
        self.projectHash = projectHash
        self.appPath = appPath
        self._config = {}
        pathToPlatformYaml = os.path.join(
            self.appPath,
            self.PLATFORM_FILENAME
        )
        with open(pathToPlatformYaml, "r") as f:
            self._config = yaml.load(f)

    def getName(self):
        return self._config.get("name", "default")

    def getBuildFlavor(self):
        build = self._config.get("build", {})
        return build.get("flavor", None)

    def getRelationships(self):
        return {}

    def getBuildHooks(self):
        hooks = self._config.get("hooks", {})
        return hooks.get("build", "")

    def getDeployHooks(self):
        hooks = self._config.get("hooks", {})
        return hooks.get("deploy", "")

    def getMounts(self):
        return self._config.get("mounts", {})

    def getRuntime(self):
        return self._config.get("runtime", {})

    def getDockerImage(self):
        """ Get name of docker image for app. """
        return self.PLATFORM_DOCKER_IMAGES.get(self.getType(), None)

    def getVariables(self):
        return self._config.get("variables", {})

    def getMounts(self):
        return self._config.get("mounts", {})

    def getWeb(self):
        return self._config.get("web", {})
