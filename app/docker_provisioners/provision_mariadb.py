import os
import difflib
import io
import hashlib
import docker
import json
from provision_base import DockerProvisionBase
from ..platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a maria db container. """

    PASSWORD_SALT = "dsf3$cb33FFhgkl4@567405-#@ldlUKJ^#lk43"

    # !!! TODO actually create users and passwords based on services.yaml !!!

    def getRootPassword(self):
        return hashlib.sha256(
            self.PASSWORD_SALT + self.appConfig.getEntropy()
        ).hexdigest()

    def getEnvironmentVariables(self):
        return {
            "MYSQL_ROOT_PASSWORD" : self.getRootPassword()
        }

    def getServiceRelationship(self):
        return [
            {
                "host" : self.container.attrs.get("Config", {}).get("Hostname", ""),
                "ip" : self.container.attrs.get("NetworkSettings", {}).get("IPAddress", ""),
                "password" : self.getRootPassword(),
                "path" : "main",
                "port" : "3306",
                "query": {
                    "is_master" : True
                },
                "scheme" : "mysql",
                "username" : "root"
            }
        ]