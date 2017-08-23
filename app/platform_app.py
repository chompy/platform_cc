import os
import yaml
from Crypto.PublicKey import RSA
from platform_app_config import PlatformAppConfig
from platform_service import PlatformService
from platform_service_config import PlatformServiceConfig
from platform_docker import PlatformDocker
from platform_web import PlatformWeb
from app.platform_utils import log_stdout, print_stdout, seperator_stdout

class PlatformApp:

    """ Base class for application. """

    def __init__(self, projectHash, appPath = ""):
        self.config = PlatformAppConfig(projectHash, appPath)
        if not os.path.exists(self.config.getDataPath()):
            os.mkdir(self.config.getDataPath())
            self.generateSshKey()
        self.docker = PlatformDocker(
            self.config,
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

    def getServices(self):
        """ Get list of service dependencies for app. """
        serviceConf = {}
        serviceList = []
        pathToServicesYaml = os.path.join(
            self.config.appPath,
            PlatformServiceConfig.PLATFORM_SERVICES_PATH
        )
        with open(pathToServicesYaml, "r") as f:
            serviceConf = yaml.load(f)
        for serviceName in serviceConf:
            serviceList.append(
                PlatformService(
                    self.config,
                    serviceName
                )
            )
        return serviceList

    def start(self):
        """ Start app. """
        log_stdout("Starting '%s' application." % self.config.getName())
        self.docker.start()
        self.web.start()
        for service in self.getServices():
            service.start()

    def stop(self):
        """ Stop app. """
        log_stdout("Stopping '%s' application." % self.config.getName())
        self.docker.stop()
        self.web.stop()
        for service in self.getServices():
            service.stop()

    def build(self):
        """ Run build and deploy hooks. """
        log_stdout("Building '%s' application." % self.config.getName())
        self.docker.syncApp()
        self.docker.preBuild()
        for service in self.getServices():
            service.docker.preBuild()
        log_stdout("Build hooks.", 1)
        results = self.docker.getContainer().exec_run(
            ["sh", "-c", self.config.getBuildHooks()],
            user="web"
        )
        seperator_stdout()
        print_stdout(results)
        seperator_stdout()

        # TODO Deploy hooks
        