import os
import difflib
import io
import hashlib
import docker
import json
import time
from provision_base import DockerProvisionBase

class DockerProvision(DockerProvisionBase):

    """ Provision a maria db container. """

    PASSWORD_SALT = "dsf3$cb33FFhgkl4@567405-#@ldlUKJ^#lk43"

    def getPassword(self, username = "root"):
        return hashlib.sha256(
            self.PASSWORD_SALT + self.appConfig.getEntropy() + username
        ).hexdigest()

    def runtime(self):
        # wait for mysqld to be ready
        while not self.healthcheck():
            time.sleep(1)
        # provision databases
        cmds = []
        config = self.appConfig.getConfiguration()
        for schema in config.get("schemas", []):
            cmds.append({
                "cmd" :     "mysql -uroot --password=\"%s\" -e \"CREATE SCHEMA %s CHARACTER SET UTF8mb4 COLLATE utf8mb4_bin;\"" % (
                    self.getPassword(),
                    schema
                ),
                "desc" :    "Create schema '%s.'" % schema
            })
        endpoints = config.get("endpoints", {})
        for endpointName in endpoints:
            endpoint = endpoints[endpointName]
            cmds.append({
                "cmd" :     "mysql -uroot --password=\"%s\" -e \"DROP USER '%s'@'%%';\"" % (
                    self.getPassword(),
                    endpointName
                ),
                "desc" :    "Drop user '%s' (if exists)." % endpointName
            })
            cmds.append({
                "cmd" :     "mysql -uroot --password=\"%s\" -e \"CREATE USER '%s'@'%%' IDENTIFIED BY '%s';\"" % (
                    self.getPassword(),
                    endpointName,
                    self.getPassword(endpointName)
                ),
                "desc" :    "Create user '%s'." % endpointName
            })
            privileges = endpoint.get("privileges", {})
            for schema in endpoint.get("privileges", {}):
                # !!TODO grant only specified permissions!!!
                cmds.append({
                    "cmd" :     "mysql -uroot --password=\"%s\" -e \"GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%';\"" % (
                        self.getPassword(),
                        schema,
                        endpointName
                    ),
                    "desc" :    "Grants privileges to '%s' on table '%s.'" % (
                        endpointName,
                        schema
                    )
                })  
        self.runCommands(cmds)

    def getEnvironmentVariables(self):
        return {
            "MYSQL_ROOT_PASSWORD" : self.getPassword()
        }

    def getServiceRelationship(self, endpointName = None):
        config = self.appConfig.getConfiguration()
        endpoints = config.get("endpoints", {})
        if endpointName != None and endpointName in endpoints.keys():
            endpoints = { endpointName: endpoints[endpointName] }
        relationships = []
        for endpointName in endpoints:
            endpoint = endpoints[endpointName]
            relationships.append({
                "host" : self.container.attrs.get("Config", {}).get("Hostname", ""),
                "ip" : self.container.attrs.get("NetworkSettings", {}).get("IPAddress", ""),
                "password" : self.getPassword(endpointName),
                "path" : endpoint.get("default_schema", ""),
                "port" : "3306",
                "query": {
                    "is_master" : True # uncertain what makes this true (maybe first endpoint is master?)
                },
                "scheme" : "mysql",
                "username" : endpointName
            })
        return relationships

    def healthcheck(self):
        if not self.container or not self.container.status == "running":
            return False
        pingResult = self.container.exec_run(
            [
                "sh",
                "-c", 
                "/usr/bin/mysqladmin ping -uroot --password='%s' --silent" % (
                    self.getPassword()
                )
            ],
            user="root",
            privileged=True,
        )
        return True if pingResult else False
