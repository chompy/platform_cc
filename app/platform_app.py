import os
import yaml
import base64
from config.platform_app_config import PlatformAppConfig
from platform_service import PlatformService
from config.platform_service_config import PlatformServiceConfig
from platform_docker import PlatformDocker
from platform_web import PlatformWeb
from app.platform_utils import log_stdout, print_stdout, seperator_stdout

class PlatformApp:

    """ Base class for application. """

    def __init__(self, projectHash, appPath = "", projectVars = {}):
        self.projectVars = projectVars
        self.config = PlatformAppConfig(projectHash, appPath, projectVars)
        self.docker = PlatformDocker(
            self.config,
            "app",
        )
        self.web = PlatformWeb(self)
        self.logIndent = 0

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

    def buildServiceRelationships(self):
        """ Build service relationship list. """
        services = self.getServices()
        relationships = self.config.getRelationships()
        output = {}
        for relationship in relationships:
            value = relationships[relationship]
            for service in services:
                serviceTypeName = service.config.getType().split(":")[0]
                if value != ("%s:%s" % (service.config.getName(), serviceTypeName)):
                    continue
                output[relationship] = service.docker.getProvisioner().getServiceRelationship()
                break
        return output

    def copySshKey(self):
        """ Copy ssh key in to container. """
        log_stdout(
            "Copy SSH key...",
            self.logIndent,
            False
        )
        sshKey = self.projectVars.get("ssh:id_rsa")
        if not sshKey:
            print_stdout("not set.")
            return
        self.docker.getContainer().exec_run(
            ["mkdir", "-p", "/app/.ssh"]
        )
        self.docker.getProvisioner().copyStringToFile(
            base64.b64decode(sshKey),
            "/app/.ssh/id_rsa"
        )
        self.docker.getContainer().exec_run(
            ["chmod", "0600", "/app/.ssh/id_rsa"]
        )
        self.docker.getContainer().exec_run(
            ["chown", "web:web", "/app/.ssh/id_rsa"]
        )
        print_stdout("done.")

    def deleteSshKey(self):
        """ Delete ssh key in container. """
        self.docker.getContainer().exec_run(
            ["rm", "-rf", "/app/.ssh"]
        )

    def start(self):
        """ Start app. """
        for service in self.getServices():
            service.start()
        self.docker.relationships = self.buildServiceRelationships()
        log_stdout(
            "Starting '%s' application." % self.config.getName(),
            self.logIndent
        )
        self.docker.start()
        self.web.start()

    def stop(self):
        """ Stop app. """
        log_stdout(
            "Stopping '%s' application." % self.config.getName(),
            self.logIndent
        )
        self.docker.stop()
        self.web.stop()
        for service in self.getServices():
            service.stop()

    def build(self):
        """ Run build and deploy hooks. """
        log_stdout(
            "Building '%s' application." % self.config.getName(),
            self.logIndent
        )
        self.docker.relationships = self.buildServiceRelationships()
        self.docker.syncApp()
        self.logIndent += 1
        self.copySshKey()
        self.logIndent -= 1
        self.docker.preBuild()
        for service in self.getServices():
            service.docker.preBuild()
        log_stdout("Build hooks.", self.logIndent)
        results = self.docker.getContainer().exec_run(
            ["sh", "-c", self.config.getBuildHooks()],
            user="web"
        )
        seperator_stdout()
        print_stdout(results)
        seperator_stdout()

        # TODO move to stand alone function
        log_stdout("Deploy hooks.", self.logIndent)
        results = self.docker.getContainer().exec_run(
            ["sh", "-c", self.config.getDeployHooks()],
            user="web"
        )
        seperator_stdout()
        print_stdout(results)
        seperator_stdout()

        self.deleteSshKey()