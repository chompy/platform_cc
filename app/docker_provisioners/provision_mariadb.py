import os
import difflib
import io
import hashlib
import docker
from provision_base import DockerProvisionBase
from ..platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a maria db container. """

    def getEnvironmentVariables(self):
        """ Get environment variables for container. """
        return {
            "MYSQL_ALLOW_EMPTY_PASSWORD" : "true"
        }