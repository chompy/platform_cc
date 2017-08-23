import os
import difflib
import io
import hashlib
import docker
from provision_base import DockerProvisionBase
from ..platform_utils import log_stdout, print_stdout

class DockerProvision(DockerProvisionBase):

    """ Provision a PHP container. """

    def provision(self):
        # parent method
        DockerProvisionBase.provision(self)
        # install extensions
        log_stdout("Install extensions.", 2)
        extensions = self.appConfig.getRuntime().get("extensions", [])
        extensionConfigs = self.config.get("extensions", {})
        for extensionName in extensions:
            log_stdout("%s..." % (extensionName), 3, False)
            extensionConfig = extensionConfigs.get(extensionName, {})
            if not extensionConfig:
                print_stdout("not available.")
                continue
            if extensionConfig.get("core", False):
                print_stdout("already installed (core extension).")
                continue
            depCmdKey = difflib.get_close_matches(
                self.image,
                extensionConfig.keys(),
                1
            )
            if not depCmdKey:
                print_stdout("not available.")
                continue
            self.container.exec_run(
                ["sh", "-c", extensionConfig[depCmdKey[0]]]
            )
            print_stdout("done.")

    def getVolumes(self):
        volumes = DockerProvisionBase.getVolumes(self)
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

    def getUid(self):
        """ Generate unique id based on configuration. """
        hashStr = self.image
        hashStr += str(self.appConfig.getBuildFlavor())
        extensions = self.appConfig.getRuntime().get("extensions", [])
        extensions.sort()
        extensionConfigs = self.config.get("extensions", {})
        for extension in extensions:
            extensionConfig = extensionConfigs.get(extension, {})
            if not extensionConfig: continue
            if not extensionConfig.get("core", False): continue
            hashStr += extension
        return hashlib.sha256(hashStr).hexdigest()