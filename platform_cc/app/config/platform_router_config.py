import os
import yaml
import yamlordereddictloader
from platform_config import PlatformConfig

class PlatformRouterConfig(PlatformConfig):

    """ Provide configuration for router. """

    ROUTER_DOCKER_IMAGE = "nginx:1.13"
    PLATFORM_ROUTES_PATH = ".platform/routes.yaml"

    def __init__(self, projectHash, projectPath):
        PlatformConfig.__init__(self, projectHash)
        routeYamlPath = os.path.join(
            projectPath,
            self.PLATFORM_ROUTES_PATH
        )
        if os.path.exists(routeYamlPath):
            with open(routeYamlPath, "r") as f:
                self.config = yaml.load(f, Loader=yamlordereddictloader.Loader)
        self.appPath = None

    def getName(self):
        return "router"

    def getMounts(self):
        return {}

    def getDockerImage(self):
        return self.ROUTER_DOCKER_IMAGE

    def getBuildFlavor(self):
        return "_router"

    def getRoutes(self):
        return self.config