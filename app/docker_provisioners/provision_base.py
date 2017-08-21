import os
import random
import string
import docker
import tarfile
import time
import io
import hashlib
import yaml
from ..platform_utils import print_stdout

class DockerProvisionBase:

    """ Base docker container provisioning class. """

    CONFIG_DIRECTORY = "%s/../../config" % (os.path.dirname(__file__))

    def __init__(self, container, platformConfig, image):
        self.container = container
        self.config = {}
        configPath = os.path.join(
            self.CONFIG_DIRECTORY,
            "%s.yaml" % (image.split(":")[0])
        )
        if os.path.isfile(configPath):
            with open(configPath, "r") as f:
                self.config = yaml.load(f)
        self.platformConfig = platformConfig
        self.image = image

    def runCommands(self, cmdList):
        """ Run commands in container. """
        for cmd in cmdList:
            requiredBuildFlavor = cmd.get("build_flavor", "")
            if requiredBuildFlavor and requiredBuildFlavor != self.platformConfig.getBuildFlavor():
                continue
            print_stdout(
                "  - %s" % (
                    cmd.get("desc", "Run command in '%s' container." % self.image)
                )
            )
            results = self.container.exec_run(
                ["sh", "-c", cmd.get("cmd", "")],
                user=cmd.get("user", "root")
            )
            if results:
                print_stdout("=======================================\n%s\n=======================================" % results)

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

    def provision(self):
        """ Provision the container. """
        self.runCommands(
            self.config.get("provision", {})
        )

    def preBuild(self):
        """ Prebuild commands. """
        self.runCommands(
            self.config.get("pre_build", {})
        )

    def getUid(self):
        """ Generate unique id based on configuration. """
        return hashlib.sha256(self.image).hexdigest()