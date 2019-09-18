"""
This file is part of Platform.CC.

Platform.CC is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

Platform.CC is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Platform.CC.  If not, see <https://www.gnu.org/licenses/>.
"""

from .base import BasePlatformService


class RedisService(BasePlatformService):
    """
    Handler for redis services.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "redis":                  "redis:3.2-alpine",
        "redis-persistent":       "redis:3.2-alpine",
        "redis:3.2":              "redis:3.2-alpine",
        "redis-persistent:3.2":   "redis:3.2-alpine",
        "redis:2.8":              "redis:2.8-alpine",
        "redis-persistent:2.8":   "redis:2.8-alpine",
        "redis:3.0":              "redis:3.0-alpine",
        "redis-persistent:3.0":   "redis:3.0-alpine",
        "redis:4.0":              "redis:4.0-alpine",
        "redis-persistent:4.0":   "redis:4.0-alpine",
        "redis:5.0":              "redis:5.0-alpine",
        "redis-persistent:5.0":   "redis:5.0-alpine"
    }

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def isPersistent(self):
        """ Return true if this redis service is persistent. """
        return (
            len(self.getType()) > 11 and
            self.getType()[-11:] == "-persistent"
        )

    def getContainerVolumes(self):
        if not self.isPersistent():
            return {}
        return {
            self.getVolumeName(): {
                "bind": "/data",
                "mode": "rw"
            }
        }

    def getContainerCommand(self):
        if not self.isPersistent():
            return BasePlatformService.getContainerCommand(self)
        return ["redis-server", "--appendonly", "yes"]

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        data["platform_relationships"][self.getName()] = {
            "host":           self.getContainerName(),
            "ip":             data.get("ip", ""),
            "scheme":         "redis",
            "port":           6379
        }
        data["platform_relationships"]["redis"] = (
            data["platform_relationships"][self.getName()]
        )
        return data
