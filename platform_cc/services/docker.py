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


class DockerService(BasePlatformService):
    """
    Handler for docker service.
    """

    def getBaseImage(self):
        return self.config.get("image", "")

    def isPlatformShCompatible(self):
        return False

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        data["platform_relationships"][self.getName()] = {
            "host":         self.getContainerName(),
            "ip":           data.get("ip", "")
        }
        data["platform_relationships"]["docker"] = (
            data["platform_relationships"][self.getName()]
        )
        return data

    def getContainerEnvironmentVariables(self):
        evars = BasePlatformService.getContainerEnvironmentVariables(self)
        for key, value in self.config.get("environment", {}).items():
            evars[key] = str(value)
        return evars

    def getContainerVolumes(self):
        volumes = {}
        for key, value in self.config.get("volumes", {}).items():
            volumes[self.getVolumeName(str(key))] = {
                "bind":     str(value),
                "mode":     "rw"
            }
        return volumes

    def getContainerCommand(self):
        return self.config.get("command")
