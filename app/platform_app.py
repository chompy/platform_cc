import os
import yaml
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

    def start(self):
        """ Start app. """
        for service in self.getServices():
            service.start()
        self.docker.relationships = self.buildServiceRelationships()
        log_stdout("Starting '%s' application." % self.config.getName())
        self.docker.start()
        self.web.start()

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
        self.docker.relationships = self.buildServiceRelationships()
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

        # TODO move to stand alone function
        log_stdout("Deploy hooks.", 1)
        results = self.docker.getContainer().exec_run(
            ["sh", "-c", self.config.getDeployHooks()],
            user="web"
        )
        seperator_stdout()
        print_stdout(results)
        seperator_stdout()
        