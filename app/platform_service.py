
class PlatformService:

    """ Platform service definition. """

    PLATFORM_SERVICE_DOCKER_IMAGES = {
        "mysql":                   "mariadb:10.2",
        "mysql:10.2":              "mariadb:10.2",
        "mysql:10.1":              "mariadb:10.1",
        "mysql:10.0":              "mariadb:10.0",
        "mysql:5.5":               "mariadb:5.5",
        "memcached":               "memcached:1",
        "memcached:1.14":          "memcached:1"
    }

    def __init__(self, name, serviceDefinition = {}):
        self.name = str(name).strip()
        self.type = serviceDefinition.get("type", "mysql:10.2")
        self.config = serviceDefinition.get("configuration", {})

    def getName(self):
        return self.name

    def getType(self):
        return self.type

    def getConfig(self):
        return self.config

    def getDockerImage(self):
        return self.PLATFORM_SERVICE_DOCKER_IMAGES.get(self.getType(), None)

    def startDocker(self):
        return

    def stopDocker(self):
        return