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
