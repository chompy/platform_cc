import os
import yaml
import hashlib
import time

class PlatformConfig:

    """ Provide configuration for the application. """

    PLATFORM_FILENAME = ".platform.app.yaml"

    PLATFORM_SERVICES_PATH = ".platform/services.yaml"

    PLATFORM_ROUTES_PATH = ".platform/routes.yaml"

    PLATFORM_LOCAL_DATA_PATH = ".platform/.pcclocal"

    PLATFORM_DOCKER_IMAGES = {
        "php:5.4":          "php:5.4-fpm",
        "php:5.6":          "php:5.6-fpm"
    }

    def __init__(self, projectPath = ""):
        self.projectPath = projectPath
        self._platformConfig = {}
        pathToPlatformYaml = os.path.join(
            self.projectPath,
            self.PLATFORM_FILENAME
        )
        with open(pathToPlatformYaml, "r") as f:
            self._platformConfig = yaml.load(f)

    def getName(self):
        return self._platformConfig.get("name", "default")

    def getType(self):
        return self._platformConfig.get("type", "php:7.0")

    def getBuildFlavor(self):
        build = self._platformConfig.get("build", {})
        return build.get("flavor", None)

    def getRelationships(self):
        return {}

    def getBuildHooks(self):
        hooks = self._platformConfig.get("hooks", {})
        return hooks.get("build", "")

    def getDeployHooks(self):
        hooks = self._platformConfig.get("hooks", {})
        return hooks.get("deploy", "")

    def getMounts(self):
        return self._platformConfig.get("mounts", {})

    def getRuntime(self):
        return self._platformConfig.get("runtime", {})

    def getServices(self):
        """ Get list of service dependencies for app. """
        serviceConf = {}
        serviceList = []
        pathToServicesYaml = os.path.join(
            self.projectPath,
            self.PLATFORM_SERVICES_PATH
        )
        with open(pathToServicesYaml, "r") as f:
            serviceConf = yaml.load(f)
        for serviceName in serviceConf:
            serviceList.append(
                PlatformService(
                    serviceName,
                    serviceConf[serviceName]
                )
            )
        return serviceList

    def getDockerImage(self):
        """ Get name of docker image for app. """
        return self.PLATFORM_DOCKER_IMAGES.get(self.getType(), None)

    def getDataPath(self):
        return os.path.join(
            self.projectPath,
            self.PLATFORM_LOCAL_DATA_PATH
        )

    def getVariables(self):
        return self._platformConfig.get("variables", {})

    def getMounts(self):
        return self._platformConfig.get("mounts", {})

    def getEntropy(self):
        entropyPath = os.path.join(self.getDataPath(), ".entropy")
        if not os.path.exists(entropyPath):
            entropy = hashlib.sha256(
                self.projectPath + yaml.dump(self._platformConfig) + self.PLATFORM_LOCAL_DATA_PATH + str(time.time())
            ).hexdigest()
            with open(entropyPath, "w") as f:
                f.write(entropy)
            return entropy
        with open(entropyPath, "r") as f:
            return f.read()