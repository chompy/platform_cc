import os
import random
import string
import docker
import tarfile
import time
import io

class DockerProvisionBase:

    """ Base docker container provisioning class. """

    def __init__(self, container, platformConfig, image):
        self.container = container
        self.platformConfig = platformConfig
        self.image = image

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