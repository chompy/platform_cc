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


class NetworkStorageService(BasePlatformService):
    """
    Handler for network-storage services.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "network-storage:1.0": "busybox:1"
    }

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getContainerVolumes(self):
        return {
            self.getVolumeName(): {
                "bind": "/data",
                "mode": "rw"
            }
        }

    def getContainerCommand(self):
        return BasePlatformService.getContainerCommand(self)

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        return data
