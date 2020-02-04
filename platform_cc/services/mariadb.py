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
from ..exception.state_error import StateError
import hashlib
import base36
import docker
import time
import sys
import io

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

    """ Default schema list if one is not set. """
    DEFAULT_SCHEMAS = [
        "mysql"
    ]

    """ Default endpoint to provide if one is not set. """
    DEFAULT_ENDPOINT = {
        "mysql": {
            "default_schema": "mysql",
            "privileges": {
                "mysql": "admin"
            }
        }
    }

    """ Salt used to generate passwords. """
    PASSWORD_SALT = "a62bf8b07e2abb117894442b00df02446670fBnBK&%2!2"

    def getBaseImage(self):
        return self.DOCKER_IMAGE_MAP.get(self.getType())

    def getPassword(self, user="root"):
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
            "MYSQL_ROOT_PASSWORD":      self.getPassword()
        }

    def getContainerVolumes(self):
        return {
            self.getVolumeName(): {
                "bind": "/var/lib/mysql",
                "mode": "rw"
            }
        }

    def getServiceData(self):
        data = BasePlatformService.getServiceData(self)
        endpoints = self.config.get("endpoints", self.DEFAULT_ENDPOINT)
        for name, config in endpoints.items():
            data["platform_relationships"][name.strip()] = {
                "host":           self.getContainerName(),
                "ip":             data.get("ip", ""),
                "port":           3306,
                "path":           config.get(
                    "default_schema",
                    list(config.get("privileges", {}).keys())[0]
                ),
                "query": {
                    "is_master":    True
                },
                "scheme":         "mysql",
                "password":       self.getPassword(name.strip()),
                "username":       name.strip()
            }
        return data

    def getCreateUserQuery(self, user):
        """ Get database query to create a user and grant nessacary prillvileges. """
        endpoints = self.config.get("endpoints", self.DEFAULT_ENDPOINT)
        endpoint = endpoints.get(user, {})
        if not endpoint: return ""
        output = ""
        # grant privileges
        privileges = endpoint.get("privileges", {})
        for schema in privileges:
            privilege = privileges[schema]
            if privilege == "admin":
                output += "GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%' IDENTIFIED BY '%s'; " % (
                    schema,
                    user,
                    self.getPassword(user)
                )
            elif privilege == "ro":
                output += "GRANT SELECT ON %s.* TO '%s'@'%%' IDENTIFIED BY '%s'; " % (
                    schema,
                    user,
                    self.getPassword(user)
                )
            elif privilege == "rw":
                output += "GRANT SELECT, INSERT, UPDATE, DELETE ON %s.* TO '%s'@'%%' IDENTIFIED BY '%s'; " % (
                    schema,
                    user,
                    self.getPassword(user)
                )
        output += "FLUSH PRIVILEGES; "
        return output

    def getCreateUsersQuery(self):
        """ Get database query to create all users and grant them nessacary prillvileges. """
        output = ""
        endpoints = self.config.get("endpoints", self.DEFAULT_ENDPOINT)
        for endpoint in endpoints:
            output += self.getCreateUserQuery(endpoint)
        return output

    def getCreateSchemasQuery(self):
        """ Get database query to create schemas. """
        output = ""
        schemas = self.config.get("schemas", self.DEFAULT_SCHEMAS)
        for schema in schemas:
            output += "CREATE SCHEMA IF NOT EXISTS %s CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin; " % str(schema)
        return output

    def start(self):
        BasePlatformService.start(self)
        container = self.getContainer()
        if not container:
            return
        # wait for mysql to become ready
        exitCode = 1
        while exitCode != 0:
            (exitCode, _) = container.exec_run(
                "mysqladmin ping -h 127.0.0.1 --silent"
            )
            time.sleep(.35)
        # create schemas
        self.logger.info("Create database schemas.")
        self.runCommand(
            """
            mysql -h 127.0.0.1 -uroot --password="%s" \
            -e "%s"
            """ % (
                self.getPassword(),
                self.getCreateSchemasQuery()
            )
        )
        # (re)create users
        self.logger.info("(Re)create database users.")
        self.runCommand(
            """
            mysql -h 127.0.0.1 -uroot --password="%s" \
            -e "%s"
            """ % (
                self.getPassword(),
                self.getCreateUsersQuery()
            )
        )

    def executeSqlDump(self, database = "", stdin = None):
        """ Upload and execute SQL dump. """
        if not self.isRunning():
            raise StateError(
                "Service '%s' is not running." % self.getName()
            )
        cmd = "mysql -h 127.0.0.1 -uroot --password=\"%s\"" % (
            self.getPassword()
        )
        if database:
            cmd += " --database=\"%s\"" % str(database)
        self.shell(cmd, stdin=stdin)
