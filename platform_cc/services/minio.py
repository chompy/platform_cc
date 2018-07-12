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
        data["platform_relationships"]["athenapdf"] = {
            "host"          : self.getContainerName(),
            "ip"            : data.get("ip", ""),
            "port"          : 9000,
            "access_key"    : self.getKey(),
            "secret_key"    : self.getKey("secret_key")
        }
        return data
