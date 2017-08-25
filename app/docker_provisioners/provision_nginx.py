import os
import difflib
import io
import hashlib
from provision_base import DockerProvisionBase
from ..platform_utils import print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a web/nginx container. """

    def getVolumes(self):
        if not self.appConfig.appPath: return
        volumes = DockerProvisionBase.getVolumes(self, "/app")
        volumes[os.path.realpath(self.appConfig.appPath)] = {
            "bind" : "/mnt/app",
            "mode" : "ro"
        }
        appVolumeKey = "%s_%s_%s_app" % (
            DockerProvisionBase.DOCKER_VOLUME_NAME_PREFIX,
            self.appConfig.projectHash[:6],
            self.appConfig.getName()
        )
        try:
            self.dockerClient.volumes.get(appVolumeKey)
        except docker.errors.NotFound:
            self.dockerClient.volumes.create(appVolumeKey)
        volumes[appVolumeKey] = {
            "bind" : "/app",
            "mode" : "ro"
        }
        return volumes