import os
import difflib
import io
import hashlib
import docker
import json
from provision_base import DockerProvisionBase
from ..platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a busybox container. This is used as a place to store project variables. """

    def getEnvironmentVariables(self):
        return {}
