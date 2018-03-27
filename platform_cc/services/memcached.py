from .base import BasePlatformService

class MemcachedService(BasePlatformService):
    """
    Handler for memcached services.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "memcached"              : "memcached:1",
        "memcached:1.4"          : "memcached:1"
    }

    def getDockerImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())
