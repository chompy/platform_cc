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
        volumes = {}
        # app volume
        if self.appConfig.appPath != None:
            volumes = DockerProvisionBase.getVolumes(self, "/app")
            volumes[os.path.realpath(self.appConfig.appPath)] = {
                "bind" : "/app",
                "mode" : "rw"
            }
        # router volume
        routerVolumeKey = "%s_router_data" % (
            DockerProvisionBase.DOCKER_VOLUME_NAME_PREFIX
        )
        try:
            self.dockerClient.volumes.get(routerVolumeKey)
        except docker.errors.NotFound:
            self.dockerClient.volumes.create(routerVolumeKey)
        volumes[routerVolumeKey] = {
            "bind" : "/router",
            "mode" : "rw"
        }
        return volumes