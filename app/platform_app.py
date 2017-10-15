import os
import yaml
import yamlordereddictloader
import base64
from config.platform_app_config import PlatformAppConfig
from platform_service import PlatformService
from config.platform_service_config import PlatformServiceConfig
from platform_docker import PlatformDocker
from platform_web import PlatformWeb

class PlatformApp:

    """ Base class for application. """

    def __init__(self, projectHash, appPath = "", projectVars = {}, logger = None):
        self.projectVars = projectVars
        self.config = PlatformAppConfig(projectHash, appPath, projectVars)
        self.docker = PlatformDocker(
            self.config,
            "%s_app" % self.config.getName(),
            self.config.getDockerImage(),
            logger
        )
        self.logger = logger
        self.logIndent = 0
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
            serviceConf = yaml.load(f, Loader=yamlordereddictloader.Loader)
        for serviceName in serviceConf:
            serviceList.append(
                PlatformService(
                    self.config,
                    serviceName,
                    self.logger
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
        if self.logger:
            self.logger.logEvent(
                "Copy SSH key.",
                self.logIndent
            )
        sshKey = self.projectVars.get("project:ssh_key")
        knownHosts = self.projectVars.get("project:known_hosts")
        if not sshKey:
            if self.logger:
                self.logger.logEvent(
                    "SSH key is not set.",
                    self.logIndent + 1
                )
            return
        self.docker.getContainer().exec_run(
            ["mkdir", "-p", "/app/.ssh"]
        )
        self.docker.getProvisioner().copyStringToFile(
            base64.b64decode(sshKey),
            "/app/.ssh/id_rsa"
        )
        if knownHosts:
            self.docker.getProvisioner().copyStringToFile(
                base64.b64decode(knownHosts),
                "/app/.ssh/known_hosts"
            )
        self.docker.getContainer().exec_run(
            ["chmod", "0600", "/app/.ssh/*"]
        )
        self.docker.getContainer().exec_run(
            ["chown", "web:web", "/app/.ssh/*"]
        )

    def deleteSshKey(self):
        """ Delete ssh key in container. """
        self.docker.getContainer().exec_run(
            ["rm", "-rf", "/app/.ssh"]
        )

    def start(self):
        """ Start app. """
        if self.logger:
            self.logger.logEvent(
                "Starting '%s' application." % self.config.getName(),
                self.logIndent
            )
        for service in self.getServices():
            service.start()
        self.docker.relationships = self.buildServiceRelationships()
        self.docker.start()
        self.web.start()

    def stop(self):
        """ Stop app. """
        if self.logger:
            self.logger.logEvent(
                "Stopping '%s' application." % self.config.getName(),
                self.logIndent
            )
        self.docker.stop()
        self.web.stop()
        for service in self.getServices():
            service.stop()

    def build(self):
        """ Run prebuild commands and build hooks. """
        if self.logger:
            self.logger.logEvent(
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
            service.preBuild()
        if self.logger:
            self.logger.logEvent(
                "Build hooks.",
                self.logIndent + 1
            )
        results = self.docker.getContainer().exec_run(
            ["sh", "-c", self.config.getBuildHooks()],
            user="web"
        )
        if results and self.logger:
            self.logger.printContainerOutput(
                results
            )
        self.deleteSshKey()

    def deploy(self):
        """ Run deploy hooks. """
        if self.logger:
            self.logger.logEvent(
                "Deploying '%s' application." % self.config.getName(),
                self.logIndent
            )
        self.docker.syncApp()
        if self.logger:
            self.logger.logEvent(
                "Deploy hooks.",
                self.logIndent + 1
            )
        results = self.docker.getContainer().exec_run(
            ["sh", "-c", self.config.getDeployHooks()],
            user="web"
        )
        if results and self.logger:
            self.logger.printContainerOutput(
                results
            )

    def shell(self):
        """ Shell in to application container. """
        if self.logger:
            self.logger.logEvent(
                "Entering shell for '%s' application." % (self.config.getName()),
                self.logIndent
            )
        self.docker.shell("bash", "web")