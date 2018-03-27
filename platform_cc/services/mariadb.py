from .base import BasePlatformService
import hashlib
import base36
import docker

class MariaDbService(BasePlatformService):
    """
    Handler for MariaDB services.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "mysql":                   "mariadb:10.2",
        "mysql:10.2":              "mariadb:10.2",
        "mysql:10.1":              "mariadb:10.1",
        "mysql:10.0":              "mariadb:10.0",
        "mysql:5.5":               "mariadb:5.5",
        "mariadb":                 "mariadb:10.2",
        "mariadb:10.2":            "mariadb:10.2",
        "mariadb:10.1":            "mariadb:10.1",
        "mariadb:10.0":            "mariadb:10.0",
        "mariadb:5.5":             "mariadb:5.5"        
    }

    """ Salt used to generate passwords. """
    PASSWORD_SALT = "a62bf8b07e2abb117894442b00df02446670fBnBK&%2!2"

    def getDockerImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getRootPassword(self):
        """ 
        Get database root password.

        :return: Root password
        :rtype: str
        """
        return base36.dumps(
            int(
                hashlib.sha256(
                    (
                        "%s-%s-%s" % (
                            self.PASSWORD_SALT,
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
            "MYSQL_ROOT_PASSWORD" :         self.getRootPassword()
        }

    def getContainerVolumes(self):
        volume = self.getVolume()
        return {
            volume.name : {
                "bind" : "/var/lib/mysql",
                "mode" : "rw"
            }
        }