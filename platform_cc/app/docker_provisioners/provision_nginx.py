from __future__ import absolute_import
import os
import difflib
import io
import hashlib
import docker
import sys
from .provision_base import DockerProvisionBase
from app.platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a web/nginx container. """

    def getVolumes(self):
        volumes = {}
        # app volume
        if self.appConfig.appPath != None:
            volumes = DockerProvisionBase.getVolumes(self, "/app")
            appPath = os.path.realpath(self.appConfig.appPath)
            # hack for docker toolbox for windows, use unix path
            if sys.platform in ["msys", "win32"]:
                appPath = appPath.split(":")
                appPath = "/%s/%s" % (
                    appPath[0].lower(),
                    ("/".join(appPath[1].split("\\"))).lstrip("/")
                )
            volumes[appPath] = {
                "bind" : "/app",
                "mode" : "rw"
            }
        return volumes