import os
import difflib
import io
import hashlib
from provision_base import DockerProvisionBase
from ..platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a web/nginx container. """

    def provision(self):
        # parent method
        DockerProvisionBase.provision(self)
