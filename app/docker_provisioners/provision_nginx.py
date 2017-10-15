import os
import difflib
import io
import hashlib
import docker
from provision_base import DockerProvisionBase
from ..platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a web/nginx container. """

    def getVolumes(self):
        if self.appConfig.appPath == None: return
        volumes = DockerProvisionBase.getVolumes(self, "/app")
        volumes[os.path.realpath(self.appConfig.appPath)] = {
            "bind" : "/app",
            "mode" : "rw"
        }
        return volumes