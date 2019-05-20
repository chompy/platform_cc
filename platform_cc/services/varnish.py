from .base import BasePlatformService
import hashlib
import base36
import docker
import time
import requests


class VarnishService(BasePlatformService):
    """
    Handler for Varnish service.
    Just a placeholder/dummy service for now.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "varnish:5.2":            "busybox",
        "varnish:6.0":            "busybox"
    }

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getContainerVolumes(self):
        return {}

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        data["platform_relationships"][self.getName()] = {
            "host":           self.getContainerName(),
            "ip":             data.get("ip", ""),
            "scheme":         "http",
            "port":           80
        }
        data["platform_relationships"]["varnish"] = (
            data["platform_relationships"][self.getName()]
        )
        return data

    def start(self):
        BasePlatformService.start(self)
