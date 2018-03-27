from .base import BasePlatformService

class MariaDbService(BasePlatformService):
    """
    Handler for MariaDB services.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "mysql":                   "mariadb:10.2",
        "mysql:10.2":              "mariadb:10.2",
        "mysql:10.1":              "mariadb:10.1",
        "mysql:10.0":              "mariadb:10.0",
        "mysql:5.5":               "mariadb:5.5",
        "mariadb":                 "mariadb:10.2",
        "mariadb:10.2":            "mariadb:10.2",
        "mariadb:10.1":            "mariadb:10.1",
        "mariadb:10.0":            "mariadb:10.0",
        "mariadb:5.5":             "mariadb:5.5"        
    }

    def getDockerImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())