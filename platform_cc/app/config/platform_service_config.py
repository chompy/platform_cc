from __future__ import absolute_import
import os
import yaml
import yamlordereddictloader
from .platform_config import PlatformConfig

class PlatformServiceConfig(PlatformConfig):

    """ Provide configuration for service. """

    PLATFORM_SERVICE_DOCKER_IMAGES = {
        "mysql":                   "mariadb:10.2",
        "mysql:10.2":              "mariadb:10.2",
        "mysql:10.1":              "mariadb:10.1",
        "mysql:10.0":              "mariadb:10.0",
        "mysql:5.5":               "mariadb:5.5",
        "memcached":               "memcached:1",
        "memcached:1.4":           "memcached:1",
        "solr:6.6":                "solr:6.6",
        "solr.6.3":                "solr:6.3",
        "redis:2.8":                "redis:2.8",
        "redis:3.0":                "redis:3.0",
        "redis:3.2":                "redis:3.2"
    }

    PLATFORM_SERVICES_PATH = ".platform/services.yaml"

    def __init__(self, projectHash, projectPath = "", name = ""):
        self.name = name.strip()
        self.projectPath = projectPath
        PlatformConfig.__init__(self, projectHash)
        pathToServiceYaml = os.path.join(
            self.projectPath,
            self.PLATFORM_SERVICES_PATH
        )
        serviceConfigs = {}
        with open(pathToServiceYaml, "r") as f:
            serviceConfigs = yaml.load(f, Loader=yamlordereddictloader.Loader)
        for serviceName in serviceConfigs:
            if serviceName == name:
                self._config = serviceConfigs[serviceName]

    def getName(self):
        return self.name

    def getMounts(self):
        return {}

    def getDockerImage(self):
        return self.PLATFORM_SERVICE_DOCKER_IMAGES.get(self.getType(), None)

    def getConfiguration(self):
        return self._config.get("configuration", {})
