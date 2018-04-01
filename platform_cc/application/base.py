import os
import docker
from container import Container

class BasePlatformApplication(Container):

    """
    Base class for managing Platform.sh applications.
    """

    """ Name of repository for all committed application images. """
    COMMIT_REPOSITORY_NAME = "platform_cc"

    """
    Directory inside container to mount application to.
    """
    APPLICATION_DIRECTORY = "/app"

    """
    Directory inside container to mount storage to.
    """
    STORAGE_DIRECTORY = "/mnt/storage"

    def __init__(self, project, config):
        """
        Constructor.

        :param project: Project data
        :param config: Service configuration
        """
        self.config = dict(config)
        Container.__init__(
            self,
            project,
            self.config.get(
                "name",
                ""
            )
        )
        self._hasCommitImage = None

    def getDockerImage(self):
        commitImage = self.getCommitImage()
        if self._hasCommitImage == None:
            try:
                self.docker.images.get(commitImage)
                self._hasCommitImage = True
            except docker.errors.ImageNotFound:
                self._hasCommitImage = False
        if self._hasCommitImage == True:
            return commitImage
        return self.getBaseImage()

    def getContainerVolumes(self):
        return {
            self.project.get("path") : {
                "bind" : self.APPLICATION_DIRECTORY,
                "mode" : "rw"                
            },
            self.getVolume().name : {
                "bind" : self.STORAGE_DIRECTORY,
                "mode" : "rw"
            }
        }

    def getContainerWorkingDirectory(self):
        return self.APPLICATION_DIRECTORY

    def getBaseImage(self):
        """
        Get base Docker image name to use for application
        prior to build.

        :return: Docker image name
        :rtype: str
        """
        return "busybox:latest"

    def getType(self):
        """
        Get application type.

        :return: Application type
        :rtype: str
        """        
        return self.config.get("type")

    def getCommitImage(self):
        """
        Get name of committed Docker image to use
        instead of main image if it exists.

        :return: Committed Docker image name
        :rtype: str
        """
        return "%s:%s_%s" % (
            self.COMMIT_REPOSITORY_NAME,
            self.getName(),
            self.project.get("short_uid")
        )

    def setupMounts(self):
        """
        Setup application defined mount points.
        """
        configMounts = self.config.get("mounts", {})
        for mountDest, config in configMounts.items():
            mountSrc = ""
            if type(config) is dict:
                if not config.get("source") == "local": continue
                mountSrc = config.get("source_path", "").strip("/")
            elif type(config) is str:
                localMountPrefx = "shared:files/"
                if not config.startswith(localMountPrefx): continue
                mountSrc = config[len(localMountPrefx):].strip("/")
            else:
                continue
            mountSrc = os.path.join(
                self.STORAGE_DIRECTORY,
                mountSrc.strip("/")
            )
            mountDest = os.path.join(
                self.APPLICATION_DIRECTORY,
                mountDest.strip("/")
            )
            self.runCommand(
                "mkdir -p %s && mount --bind %s %s" % (
                    mountSrc,
                    mountSrc,
                    mountDest
                )
            )

    def build(self):
        """
        Run first time commands on application container. Then
        commit the container.
        """
        pass