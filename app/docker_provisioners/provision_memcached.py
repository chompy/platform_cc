import os
import difflib
import io
import hashlib
import docker
from provision_base import DockerProvisionBase
from ..platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a memcached container. """

    def getVolumes(self):
        return {}