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

class RabbitMqService(BasePlatformService):
    """
    Handler for Rabbitmq services.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "rabbitmq:3.5"           : "rabbitmq:3",
        "rabbitmq:3.5"           : "rabbitmq:3"
    }

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getContainerVolumes(self):
        return {
            self.getVolumeName() : {
                "bind" : "/var/lib/rabbitmq",
                "mode" : "rw"
            }
        }

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        data["platform_relationships"][self.getName()] = {
            "host"          : self.getContainerName(),
            "ip"            : data.get("ip", ""),
            "scheme"        : "amqp",
            "port"          : 5672,
            "username"      : "guest",
            "password"      : "guest"
        }
        return data
