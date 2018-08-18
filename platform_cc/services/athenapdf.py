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
import hashlib
import base36

class AthenaPdfService(BasePlatformService):
    """
    Handler for Athena PDF services.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "athenapdf"              : "arachnysdocker/athenapdf-service"
    }

    """ Salt used to generate auth key. """
    AUTH_SALT = "9IxcRXqOxYTbChKf0CQ1apzw26oyKrbsBfd#$2!2"

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def isPlatformShCompatible(self):
        return False

    def getAuthKey(self):
        """ 
        Get auth key to use when authenticating with weaver service.

        :rtype: str
        """
        return base36.dumps(
            int(
                hashlib.sha256(
                    (
                        "%s-%s-%s-%s" % (
                            self.AUTH_SALT,
                            self.getName(),
                            self.project.get("entropy", ""),
                            self.project.get("uid", "")
                        )
                    ).encode("utf-8")
                ).hexdigest(),
                16
            )
        )

    def getContainerEnvironmentVariables(self):
        return {
            "WEAVER_AUTH_KEY"       : self.getAuthKey()
        }

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        data["platform_relationships"]["athenapdf"] = {
            "host"          : self.getContainerName(),
            "ip"            : data.get("ip", ""),
            "port"          : 8080,
            "auth"          : self.getAuthKey()
        }
        return data
