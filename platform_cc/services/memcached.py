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

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        data["platform_relationships"]["memcached"] = {
            "host"          : data.get("ip", ""),
            "ip"            : data.get("ip", ""),
            "scheme"        : "memcached",
            "port"          : 11211
        }
        return data