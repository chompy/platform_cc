from .base import BasePlatformApplication

class PhpApplication(BasePlatformApplication):
    """
    Handler for PHP applications.
    """

    """ Mapping for application type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "php"             : "php:5.6-fpm",
        "php:5.4"         : "php:5.4-fpm",
        "php:5.6"         : "php:5.6-fpm",
        "php:7.0"         : "php:7.0-fpm"   
    }

    def getDockerImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())