import os
import yaml
from Crypto.PublicKey import RSA
from platform_app_config import PlatformAppConfig
from platform_service import PlatformService
from platform_docker import PlatformDocker
from platform_web import PlatformWeb
from platform_utils import print_stdout

class PlatformNoAvailableDockerImageException(Exception):
    """ 
    Exception that signifies that the current 
    app has no available docker image.
    """
    pass

class PlatformApp:

    """ Base class for application. """

    def __init__(self, projectHash, appPath = ""):
        self.config = PlatformAppConfig(projectHash, appPath)
        if not os.path.exists(self.config.getDataPath()):
            os.mkdir(self.config.getDataPath())
            self.generateSshKey()
        self.docker = PlatformDocker(
            self.config,
            self.config.getDockerImage(),
            "app"
        )
        self.web = PlatformWeb(self)

    def generateSshKey(self):
        """ Generate SSH key for use inside containers. """
        key = RSA.generate(2048)
        sshKeyPath = os.path.join(self.config.getDataPath(), "id_rsa")
        with open(sshKeyPath, 'w') as f:
            os.chmod(sshKeyPath, 0600)
            f.write(key.exportKey('PEM'))
        pubkey = key.publickey()
        pubkeyPath = os.path.join(self.config.getDataPath(), "id_rsa.pub")
        with open(pubkeyPath, 'w') as f:
            f.write(pubkey.exportKey('OpenSSH'))

    def start(self):
        """ Start app. """
        print_stdout("> Starting '%s' app." % self.config.getName())
        baseImage = self.config.getDockerImage()
        if not baseImage:
            raise PlatformNoAvailableDockerImageException(
                "No Docker image available for app type '%s.'" % self.getType()
            )
        self.docker.start()
        self.web.start()

    def stop(self):
        """ Stop app. """
        print_stdout("> Stopping '%s' app." % self.config.getName())
        baseImage = self.config.getDockerImage()
        if not baseImage:
            raise PlatformNoAvailableDockerImageException(
                "No Docker image available for app type '%s.'" % self.getType()
            )
        self.docker.stop()
        self.web.stop()

    def build(self):
        """ Run build and deploy hooks. """
        print_stdout("> Building application...")
        self.docker.syncApp()
        self.docker.preBuild()
        print_stdout("  - Build hooks.")
        results = self.docker.getContainer().exec_run(
            ["sh", "-c", self.config.getBuildHooks()],
            user="web"
        )
        print_stdout("=======================================\n%s\n=======================================" % results)
        # TODO Deploy hooks
        