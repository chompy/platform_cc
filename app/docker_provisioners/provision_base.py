import os
import random
import string
import docker
import tarfile
import time
import io
import hashlib
import yaml
from ..platform_utils import log_stdout, print_stdout, seperator_stdout

class DockerProvisionBase:

    """ Base docker container provisioning class. """

    CONFIG_DIRECTORY = "%s/../../config" % (os.path.dirname(__file__))

    DOCKER_VOLUME_NAME_PREFIX = "pcc"

    def __init__(self, dockerClient, container, appConfig, image = None):
        self.dockerClient = dockerClient
        self.container = container
        self.provisionConfig = {}
        self.image = image if image else appConfig.getDockerImage()
        configPath = os.path.join(
            self.CONFIG_DIRECTORY,
            "%s.yaml" % (self.image.split(":")[0])
        )
        self.config = {}
        if os.path.isfile(configPath):
            with open(configPath, "r") as f:
                self.config = yaml.load(f)
                if not self.config: self.config = {}
        self.appConfig = appConfig

    def runCommands(self, cmdList):
        """ Run commands in container. """
        for cmd in cmdList:
            requiredBuildFlavor = cmd.get("build_flavor", "")
            if requiredBuildFlavor and requiredBuildFlavor != self.appConfig.getBuildFlavor():
                continue
            log_stdout(
                cmd.get("desc", "Run command in '%s' container." % self.image),
                2
            )
            results = self.container.exec_run(
                ["sh", "-c", cmd.get("cmd", "")],
                user=cmd.get("user", "root")
            )
            if results:
                seperator_stdout()
                print_stdout(results)
                seperator_stdout()

    def randomString(self, length):
        """ Utility method. Generate random string. """
        return ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(length))

    def copyFile(self, localSrc, containerDest):
        """ Copy a local file in to the container. """
        tarData = io.BytesIO()
        with tarfile.open(fileobj=tarData, mode="w") as tar:
            containerFilename = os.path.basename(containerDest)
            if not containerFilename:
                containerFilename = os.path.basename(localSrc)
            tarFileInfo = tarfile.TarInfo(name=containerFilename)
            tarFileInfo.size = os.path.getsize(localSrc)
            tarFileInfo.mtime = time.time()
            with open(localSrc, mode="rb") as f:
                tar.addfile(
                    tarFileInfo,
                    f
                )
        tarData.seek(0)
        self.container.put_archive(
            os.path.dirname(containerDest),
            data=tarData
        )

    def copyStringToFile(self, stringData, containerDest):
        """ Copy a string to a file inside the container. """
        tarData = io.BytesIO()
        with tarfile.open(fileobj=tarData, mode="w") as tar:
            containerFilename = os.path.basename(containerDest)
            if not containerFilename:
                containerFilename = os.path.basename(localSrc)
            tarFileInfo = tarfile.TarInfo(name=containerFilename)
            tarFileInfo.size = len(stringData)
            tarFileInfo.mtime = time.time()
            tar.addfile(
                tarFileInfo,
                io.BytesIO(stringData)
            )
        tarData.seek(0)
        self.container.put_archive(
            os.path.dirname(containerDest),
            data=tarData
        )

    def fetchFile(self, containerDest):
        """ Retrieve contents of file from within container. """
        tarData, stat = self.container.get_archive(containerDest)
        tarData = io.BytesIO(tarData.read())
        results = ""
        with tarfile.open(fileobj=tarData, mode="r") as tar:
            f = tar.extractfile(os.path.basename(containerDest))
            results = f.read()
            f.close()
        return results

    def provision(self):
        """ Provision the container. """
        self.runCommands(
            self.config.get("provision", {})
        )

    def preBuild(self):
        """ Prebuild commands. (Run prior to build hooks.) """
        self.runCommands(
            self.config.get("pre_build", {})
        )

    def runtime(self):
        """ Runtime commands. """
        self.runCommands(
            self.config.get("runtime", {})
        )

    def getVolumes(self, destPrefix = ""):
        """ Get/create volumes to mount. """
        mounts = self.appConfig.getMounts()
        mounts.update(self.config.get("mounts", {}))
        volumes = {}
        for mountDest in mounts:
            mountKey = mounts[mountDest]
            dockerVolumeKey = ("%s_%s_%s_%s" % (
                self.DOCKER_VOLUME_NAME_PREFIX,
                self.appConfig.projectHash[:6],
                self.appConfig.getName(),
                os.path.basename(mountKey)
            )).rstrip("_")
            try:
                self.dockerClient.volumes.get(dockerVolumeKey)
            except docker.errors.NotFound:
                self.dockerClient.volumes.create(dockerVolumeKey)
            volumes[dockerVolumeKey] = {
                "bind" : ("%s/%s" % (
                    destPrefix,
                    mountDest.lstrip("/")
                )).rstrip("/"),
                "mode" : "rw"
            }
        return volumes

    def getEnvironmentVariables(self):
        """ Get environment variables for container. """
        return {}

    def getServiceRelationship(self):
        """ Get service configuration to expose to application via PLATFORM_RELATIONSHIPS. """
        return []

    def getUid(self):
        """ Generate unique id based on configuration. """
        hashStr = self.image
        hashStr += str(self.appConfig.getBuildFlavor())
        return hashlib.sha256(hashStr).hexdigest()