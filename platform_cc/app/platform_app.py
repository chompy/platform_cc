import os
import yaml
import yamlordereddictloader
import base64
import collections
import docker
from config.platform_app_config import PlatformAppConfig
from platform_docker import PlatformDocker
from platform_web import PlatformWeb

class PlatformApp:

    """ Base class for application. """

    def __init__(self, projectHash, appPath = "", services = [], projectVars = {}, routerConfig = {}, logger = None):
        self.services = services
        self.projectVars = projectVars
        self.config = PlatformAppConfig(projectHash, appPath, projectVars, routerConfig)
        self.docker = PlatformDocker(
            self.config,
            "%s_app" % self.config.getName(),
            self.config.getDockerImage(),
            logger
        )
        self.logger = logger
        self.logIndent = 0
        self.web = PlatformWeb(self)

    def buildServiceRelationships(self):
        """ Build service relationship list. """
        relationships = self.config.getRelationships()
        output = {}
        for relationship in relationships:
            value = relationships[relationship]
            relationshipServiceTypeName = value.split(":")[0]
            for service in self.services:
                if relationshipServiceTypeName != service.config.getName():
                    continue
                endpointName = value.split(":")[1]
                if service.config.getName() == "mysqldb":
                    serviceRelationship = service.docker.getProvisioner().getServiceRelationship(endpointName)
                else:
                    serviceRelationship = service.docker.getProvisioner().getServiceRelationship()
                output[relationship] = serviceRelationship
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
        self.docker.getContainer().exec_run(
            ["chmod", "0600", "/app/.ssh/id_rsa"]
        )
        self.docker.getContainer().exec_run(
            ["chown", "web:web", "/app/.ssh/id_rsa"]
        )
        if knownHosts:
            self.docker.getProvisioner().copyStringToFile(
                base64.b64decode(knownHosts),
                "/app/.ssh/known_hosts"
            )
            self.docker.getContainer().exec_run(
                ["chmod", "0600", "/app/.ssh/known_hosts"]
            )
            self.docker.getContainer().exec_run(
                ["chown", "web:web", "/app/.ssh/known_hosts"]
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
        if self.logger:
            self.logger.logEvent(
                "Starting main application container.",
                self.logIndent + 1
            )
        
        self.docker.logIndent += 1
        self.docker.relationships = self.buildServiceRelationships()
        self.docker.start()
        self.docker.logIndent -= 1
        self.web.start()

    def stop(self):
        """ Stop app. """
        if self.logger:
            self.logger.logEvent(
                "Stopping '%s' application." % self.config.getName(),
                self.logIndent
            )
            self.logger.logEvent(
                "Stopping main application container.",
                self.logIndent + 1
            )
        self.docker.logIndent += 1
        self.docker.stop()
        self.docker.logIndent -= 1
        self.web.stop()

    def provision(self):
        """ Provision app and run build hooks. """
        if self.logger:
            self.logger.logEvent(
                "Provision '%s' application." % self.config.getName(),
                self.logIndent
            )
        # set app file permission
        self.docker.getContainer().exec_run(
            ["chown", "-R", "web:web", "/app"]
        )
        # set vars / copy ssh key
        self.docker.relationships = self.buildServiceRelationships()
        self.logIndent += 1
        self.copySshKey()
        self.logIndent -= 1
        # provision app
        self.docker.provision(False) # no commit
        # build hooks
        if self.logger:
            self.logger.logEvent(
                "Build hooks.",
                self.logIndent + 1
            )
        self.docker.getContainer().restart()
        results = self.docker.getContainer().exec_run(
            ["sh", "-c", self.config.getBuildHooks()],
            user="web"
        )
        if results and self.logger:
            self.logger.printContainerOutput(
                results
            )
        # hack to make sites that rely on binaries being in /usr/bin
        # to work
        self.docker.getContainer().exec_run(
            ["ln", "-s", "/usr/local/bin/*", "/usr/bin/"],
            user="root"
        )
        # commit provisioned app container
        self.docker.commit()

    def deploy(self):
        """ Run deploy hooks. """
        if self.logger:
            self.logger.logEvent(
                "Deploying '%s' application." % self.config.getName(),
                self.logIndent
            )
        results = self.docker.getContainer().exec_run(
            ["sh", "-c", self.config.getDeployHooks()],
            user="web"
        )
        if results and self.logger:
            self.logger.printContainerOutput(
                results
            )

    def shell(self, cmd = "bash", user = "web"):
        """ Shell in to application container. """
        self.docker.shell(cmd, user)

    def purge(self):
        """ Purge application. """
        self.stop()
        # purge all docker instances
        self.docker.purge()
