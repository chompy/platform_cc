import os
import yaml
from platform_config import PlatformConfig

class PlatformVarsConfig(PlatformConfig):

    """ Provide configuration for vars setter/fetcher. """

    VARS_DOCKER_IMAGE = "busybox:latest"

    def __init__(self, projectHash):
        PlatformConfig.__init__(self, projectHash)

    def getName(self):
        return "vars"

    def getMounts(self):
        return {
            "/data" : "data"
        }

    def getDockerImage(self):
        return self.VARS_DOCKER_IMAGE
