from .base import BasePlatformService
import hashlib
import base36
import docker
import time

class MariaDbService(BasePlatformService):
    """
    Handler for MariaDB services.
    """

    """ Mapping for service type to Docker image name. """
    DOCKER_IMAGE_MAP = {
        "mysql"                  : "mariadb:10.2",
        "mysql:10.2"             : "mariadb:10.2",
        "mysql:10.1"             : "mariadb:10.1",
        "mysql:10.0"             : "mariadb:10.0",
        "mysql:5.5"              : "mariadb:5.5",
        "mariadb"                : "mariadb:10.2",
        "mariadb:10.2"           : "mariadb:10.2",
        "mariadb:10.1"           : "mariadb:10.1",
        "mariadb:10.0"           : "mariadb:10.0",
        "mariadb:5.5"            : "mariadb:5.5"        
    }

    """ Salt used to generate passwords. """
    PASSWORD_SALT = "a62bf8b07e2abb117894442b00df02446670fBnBK&%2!2"

    def getDockerImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getPassword(self, user = "root"):
        """ 
        Get database password for given user.

        :param user: Database username
        :return: Password
        :rtype: str
        """
        return base36.dumps(
            int(
                hashlib.sha256(
                    (
                        "%s-%s-%s-%s-%s" % (
                            self.PASSWORD_SALT,
                            str(user),
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
            "MYSQL_ROOT_PASSWORD"       : self.getPassword()
        }

    def getContainerVolumes(self):
        volume = self.getVolume()
        return {
            volume.name : {
                "bind" : "/var/lib/mysql",
                "mode" : "rw"
            }
        }

    def getPlatformRelationship(self):
        return {
            "host"          : self.getContainerIpAddress(),
            "ip"            : self.getContainerIpAddress(),
            "query"         : {
                "is_master"     : True
            },
            "scheme"        : "mysql",
            "path"          : "main",
            "password"      : "",
            "username"      : ""
        }

    def start(self):
        BasePlatformService.start(self)
        container = self.getContainer()
        if not container: return
        # wait for mysql to become ready
        exitCode = 1
        while exitCode != 0:
            (exitCode, output) = container.exec_run(
                "mysqladmin ping -h 127.0.0.1 --silent"
            )
            time.sleep(.25)
        # create schemas
        schemas = self.config.get("schemas", [])
        for schema in schemas:
            container.exec_run(
                "mysql -h 127.0.0.1 -uroot --password=\"%s\" -e \"CREATE SCHEMA `%s` CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin;\"" % (
                    self.getPassword(),
                    schema
                )
            )
        # create users
        endpoints = self.config.get("endpoints", {})
        for endpoint in endpoints:
            container.exec_run(
                "mysql -h 127.0.0.1 -uroot --password=\"%s\" -e \"DROP USER '%s'@'%%';\"" % (
                    self.getPassword(),
                    endpoint
                )
            )
            container.exec_run(
                "mysql -h 127.0.0.1 -uroot --password=\"%s\" -e \"CREATE USER '%s'@'%%' IDENTIFIED BY '%s';\"" % (
                    self.getPassword(),
                    endpoint,
                    self.getPassword(endpoint)
                )
            )
            privileges = endpoints[endpoint].get("privileges", {})
            for schema in privileges:
                privilege = privileges[schema]
                if privilege == "admin":
                    container.exec_run(
                        "mysql -h 127.0.0.1 -uroot --password=\"%s\" -e \"GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%';\"" % (
                            self.getPassword(),
                            schema,
                            endpoint
                        )
                    )
                elif privilege == "ro":
                    container.exec_run(
                        "mysql -h 127.0.0.1 -uroot --password=\"%s\" -e \"GRANT SELECT ON %s.* TO '%s'@'%%';\"" % (
                            self.getPassword(),
                            schema,
                            endpoint
                        )
                    )                    
                elif privilege == "rw":
                    container.exec_run(
                        "mysql -h 127.0.0.1 -uroot --password=\"%s\" -e \"GRANT SELECT, INSERT, UPDATE, DELETE ON %s.* TO '%s'@'%%';\"" % (
                            self.getPassword(),
                            schema,
                            endpoint
                        )
                    )
        container.exec_run(
            "mysql -h 127.0.0.1 -uroot --password=\"%s\" -e \"FLUSH PRIVILEGES;\"" % (
                self.getPassword()
            )
        )