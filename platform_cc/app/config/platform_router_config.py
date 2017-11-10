from __future__ import absolute_import
import os
import yaml
import yamlordereddictloader
from .platform_config import PlatformConfig

class PlatformRouterConfig(PlatformConfig):

    """ Provide configuration for router. """

    ROUTER_DOCKER_IMAGE = "nginx:1.13"
    
    def __init__(self):
        PlatformConfig.__init__(self, "global")
        self.appPath = None

    def getName(self):
        return "router"

    def getMounts(self):
        return {}

    def getDockerImage(self):
        return self.ROUTER_DOCKER_IMAGE

    def getBuildFlavor(self):
        return "_router"
