import os
import time
import hashlib

from platform_app import PlatformApp
from platform_app_config import PlatformAppConfig

class PlatformProject:

    """ Base class for project. """

    HASH_SECRET = "4bcc181ab1f9fcc64a8c935686b55ca794e76d63"

    def __init__(self, projectPath = ""):
        self.projectPath = projectPath
        dataPath = os.path.join(
            projectPath,
            PlatformAppConfig.PLATFORM_LOCAL_DATA_PATH
        )
        if os.path.exists(projectPath):
            if not os.path.exists(dataPath):
                os.makedirs(dataPath)
        projectHashPath = os.path.join(dataPath, ".projectId")
        self.projectHash = ""
        if not os.path.isfile(projectHashPath):
            self.projectHash = hashlib.sha256(
                self.HASH_SECRET + os.getuid() + str(time.time())
            ).hexdigest()
            with open(projectHashPath, "w") as f:
                f.write(self.projectHash)
        if not self.projectHash:
            with open(projectHashPath, "r") as f:
                self.projectHash = f.read()

    def getApplications(self):
        """ Get all applications in project. """
        topPlatformAppConfigPath = os.path.join(self.projectPath, PlatformAppConfig.PLATFORM_FILENAME)
        if os.path.exists(topPlatformAppConfigPath):
            return [PlatformApp(self.projectHash, self.projectPath)]
        apps = []
        for path in os.listdir(self.projectPath):
            path = os.path.join(self.projectPath, path)
            if os.path.isdir(path):
                platformAppConfigPath = os.path.join(path, PlatformAppConfig.PLATFORM_FILENAME)
                if os.path.isfile(platformAppConfigPath):
                    apps.append(PlatformApp(self.projectHash, self.projectPath))
        return apps

