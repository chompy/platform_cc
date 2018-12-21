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

class MinioService(BasePlatformService):
    """
    Handler for Minio object store service.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "minio"              : "minio/minio"
    }

    """ Salt used to generate access keys. """
    KEY_SALT = "5ACN5ncNxaLTAOsBuSgqzb2ySzIyxK8F$5^&4F"

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def isPlatformShCompatible(self):
        return False

    def getKey(self, type = "access_key"):
        """ 
        Get key to use for for minio access

        :rtype: str
        """
        return base36.dumps(
            int(
                hashlib.sha256(
                    (
                        "%s-%s-%s-%s-%s" % (
                            self.KEY_SALT,
                            self.getName(),
                            self.project.get("entropy", ""),
                            self.project.get("uid", ""),
                            type
                        )
                    ).encode("utf-8")
                ).hexdigest(),
                16
            )
        )

    def getContainerEnvironmentVariables(self):
        return {
            "MINIO_ACCESS_KEY"       : self.getKey(),
            "MINIO_SECRET_KEY"       : self.getKey("secret_key")
        }

    def getContainerCommand(self):
        return "server /data"

    def getContainerVolumes(self):
        return {
            self.getVolumeName() : {
                "bind" : "/data",
                "mode" : "rw"
            }
        }

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        data["platform_relationships"][self.getName()] = {
            "host"          : self.getContainerName(),
            "ip"            : data.get("ip", ""),
            "port"          : 9000,
            "access_key"    : self.getKey(),
            "secret_key"    : self.getKey("secret_key")
        }
        data["platform_relationships"]["minio"] = data["platform_relationships"][self.getName()]
        return data
