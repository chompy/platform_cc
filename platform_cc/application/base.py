import os
import json
import base64
import docker
import logging
from container import Container
from exception.state_error import StateError

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
        self.logger = logging.getLogger(__name__)
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

    def getContainerEnvironmentVariables(self):
        # get platform relationships
        platformRelationships = {}
        for key, value in self.config.get("relationships", {}).items():
            value = value.strip().split(":")
            platformRelationships[key] = [
                self.project.get("services", {})
                    .get(value[0], {})
                    .get("platform_relationships", {})
                    .get(value[1])
            ]
        envVars = {
            "PLATFORM_APP_DIR"          : self.APPLICATION_DIRECTORY,
            "PLATFORM_APPLICATION"      : "",
            "PLATFORM_APPLICATION_NAME" : self.getName(),
            "PLATFORM_BRANCH"           : "",
            "PLATFORM_DOCUMENT_ROOT"    : "/",
            "PLATFORM_ENVIRONMENT"      : "",
            "PLATFORM_PROJECT"          : self.project.get("uid", ""),
            "PLATFORM_RELATIONSHIPS"    : base64.b64encode(
                bytes(str(json.dumps(platformRelationships)).encode("utf-8"))
            ).decode("utf-8"),
            "PLATFORM_ROUTES"           : "", # TODO
            "PLATFORM_TREE_ID"          : "",
            "PLATFORM_VARIABLES"        : base64.b64encode(
                bytes(str(json.dumps(self.project.get("variables", {}))).encode("utf-8"))
            ).decode("utf-8"),
            "PLATFORM_PROJECT_ENTROPY"  : self.project.get("entropy", ""),
            "TRUSTED_PROXIES"           : "172.0.0.0/8,127.0.0.1"
        }
        # set env vars from app variables
        for key, value in self.config.get("variables", {}).get("env", {}).items():
            envVars[key.strip().upper()] = str(value)
        # set env vars from project variables
        for key, value in self.project.get("variables", {}).items():
            if not key.startswith("env:"): continue
            key = key[4:]
            envVars[key.strip().upper()] = str(value)
        
        return envVars

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
        Run commands needed to get container ready for given
        application. Also runs build hooks commands.
        """
        self.logger.info("Build '%s' application." % self.getName())

    def deploy(self):
        """
        Run deploy hook commands.
        """
        self.logger.info("Run deploy hooks for '%s' application." % self.getName())

    def start(self):
        self.logger.info("Start '%s' application." % self.getName())
        # ensure all required services are available
        projectServices = self.project.get("services", {})
        serviceNames = list(self.config.get("relationships", {}).values())
        for serviceName in serviceNames:
            serviceName = serviceName.strip().split(":")[0]
            projectService = projectServices.get(serviceName)
            if not projectService or not projectService.get("running"):
                raise StateError(
                    "Application '%s' depends on service '%s' which is not running." % (
                        self.getName(),
                        serviceName
                    )
                )            
        Container.start(self)