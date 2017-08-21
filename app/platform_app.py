import os
import yaml
from Crypto.PublicKey import RSA
from platform_config import PlatformConfig
from platform_service import PlatformService
from platform_docker import PlatformDocker
from platform_utils import print_stdout

class PlatformNoAvailableDockerImageException(Exception):
    """ 
    Exception that signifies that the current 
    app has no available docker image.
    """
    pass

class PlatformApp:

    """ Base class for platform.sh application. """

    def __init__(self, projectPath = ""):
        self.config = PlatformConfig(projectPath)
        if not os.path.exists(self.config.getDataPath()):
            os.mkdir(self.config.getDataPath())
            self.generateSshKey()

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
        baseDocker = PlatformDocker(
            self.config,
            self.config.getDockerImage()
        )
        baseDocker.start()

    def stop(self):
        """ Stop app. """
        print_stdout("> Stopping '%s' app." % self.config.getName())
        baseImage = self.config.getDockerImage()
        if not baseImage:
            raise PlatformNoAvailableDockerImageException(
                "No Docker image available for app type '%s.'" % self.getType()
            )
        baseDocker = PlatformDocker(
            self.config,
            self.config.getDockerImage()
        )
        baseDocker.stop()
        print_stdout("> Done.")

    def build(self):
        """ Run build and deploy hooks. """
        print_stdout("> Building application...")
        baseDocker = PlatformDocker(
            self.config,
            self.config.getDockerImage()
        )
        baseDocker.preBuild()
        print_stdout("  - Build hooks...", False)
        baseDocker.getContainer().exec_run(
            ["sh", "-c", self.config.getBuildHooks()],
            user="web"
        )
        print_stdout("done.")
        # TODO Deploy hooks
        