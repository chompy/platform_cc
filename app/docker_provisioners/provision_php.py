import os
import difflib
import io
import hashlib
import docker
from provision_base import DockerProvisionBase

class DockerProvision(DockerProvisionBase):

    """ Provision a PHP container. """

    def provision(self):
        # parent method
        DockerProvisionBase.provision(self)
        # install extensions
        if self.logger:
            self.logger.logEvent(
                "Install extensions.",
                self.logIndent
            )
        extensions = self.appConfig.getRuntime().get("extensions", [])
        extensionConfigs = self.config.get("extensions", {})
        for extensionName in extensions:
            if self.logger:
                self.logger.logEvent(
                    "%s." % (extensionName),
                    self.logIndent + 1
                )
            extensionConfig = extensionConfigs.get(extensionName, {})
            if not extensionConfig:
                if self.logger:
                    self.logger.logEvent(
                        "Extension not available.",
                        self.logIndent + 2
                    )
                continue
            if extensionConfig.get("core", False):
                if self.logger:
                    self.logger.logEvent(
                        "No additional configuration nessacary.",
                        self.logIndent + 2
                    )
                continue
            depCmdKey = difflib.get_close_matches(
                self.image,
                extensionConfig.keys(),
                1
            )
            if not depCmdKey:
                if self.logger:
                    self.logger.logEvent(
                        "Extension not available.",
                        self.logIndent + 2
                    )
                continue
            self.container.exec_run(
                ["sh", "-c", extensionConfig[depCmdKey[0]]]
            )

    def getVolumes(self):
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
            "mode" : "rw"
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
            if extensionConfig.get("core", False): continue
            hashStr += extension
        return hashlib.sha256(hashStr).hexdigest()