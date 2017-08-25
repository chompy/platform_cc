import os
import time
import hashlib
from Crypto.PublicKey import RSA
from platform_app import PlatformApp
from config.platform_app_config import PlatformAppConfig
from platform_vars import PlatformVars
from platform_router import PlatformRouter

class ProjectNotFoundException(Exception):
    pass

class PlatformProject:

    """ Base class for project. """

    HASH_SECRET = "4bcc181ab1f9fcc64a8c935686b55ca794e76d63"

    def __init__(self, projectPath = ""):
        self.projectPath = projectPath
        if not os.path.isdir(os.path.realpath(projectPath)):
            raise ProjectNotFoundException("Could not find project at '%s.'" % os.path.realpath(projectPath))
        projectHashPath = os.path.join(projectPath, ".pcc_project_id")
        self.projectHash = ""
        if not os.path.isfile(projectHashPath):
            self.projectHash = hashlib.sha256(
                self.HASH_SECRET + str(os.getuid()) + str(time.time())
            ).hexdigest()
            with open(projectHashPath, "w") as f:
                f.write(self.projectHash)
        if not self.projectHash:
            with open(projectHashPath, "r") as f:
                self.projectHash = f.read()
        self.vars = PlatformVars(self.projectHash)
        self.router = PlatformRouter(
            self.projectHash,
            self.projectPath
        )

    def getApplications(self, withVars = True):
        """ Get all applications in project. """
        topPlatformAppConfigPath = os.path.join(self.projectPath, PlatformAppConfig.PLATFORM_FILENAME)
        projectVars = {}
        if withVars:
            projectVars = self.vars.all()
        if os.path.exists(topPlatformAppConfigPath):
            return [PlatformApp(self.projectHash, self.projectPath, projectVars)]
        apps = []
        for path in os.listdir(self.projectPath):
            path = os.path.join(self.projectPath, path)
            if os.path.isdir(path):
                platformAppConfigPath = os.path.join(path, PlatformAppConfig.PLATFORM_FILENAME)
                if os.path.isfile(platformAppConfigPath):
                    apps.append(PlatformApp(self.projectHash, self.projectPath))
        return apps

    def generateSshKey(self):
        """ Generate SSH key for use inside containers. """
        key = RSA.generate(2048)
        self.vars.set(
            "private_key",
            key.exportKey('PEM')
        )
        pubkey = key.publickey()
        self.vars.set(
            "public_key",
            pubkey.exportKey('OpenSSH')
        )
