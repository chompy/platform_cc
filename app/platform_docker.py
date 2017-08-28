import os
import docker
import importlib
import yaml
import json
import base64
import hashlib

class PlatformDocker:

    """ Manage and provision a platform service. """

    DOCKER_CONTAINER_NAME_PREFIX = "pcc"

    DOCKER_COMMIT_REPO = "platform_cc"

    def __init__(self, config, name = None, image = None, logger = None):
        self.dockerClient = docker.from_env()
        self.config = config
        self.name = name if name else self.config.getName()
        self.image = image if image else self.config.getDockerImage()
        self.containerId = "%s_%s_%s_%s" % (
            self.DOCKER_CONTAINER_NAME_PREFIX,
            self.config.projectHash[:6],
            self.image.split(":")[0],
            self.name
        )
        self.networkId = "%s_%s_network" % (
            self.DOCKER_CONTAINER_NAME_PREFIX,
            self.config.projectHash[:6]
        )
        self.relationships = {}
        self.logger = logger
        self.logIndent = 1

    def getContainer(self):
        """ Get docker container. """
        return self.dockerClient.containers.get(self.containerId)

    def getProvisioner(self):
        """ Get provisoner object for this docker container. """
        container = None
        try:
            container = self.getContainer()
        except docker.errors.NotFound:
            pass
        provisionModule = importlib.import_module("app.docker_provisioners.provision_%s" % self.image.split(":")[0])
        provisioner = provisionModule.DockerProvision(
            self.dockerClient,
            container,
            self.config,
            self.image,
            self.logger
        )
        provisioner.logIndent = self.logIndent
        return provisioner

    def getTag(self):
        """ Get unique tag name for this container's configuration. """
        return "%s_%s" % (
            self.image.split(":")[0],
            self.getProvisioner().getUid()[:6]
        )

    def getVariables(self):
        """ Get project variables. """
        return self.config.getVariables()

    def getEnvironmentVariables(self):
        """ Get environment variables to use with container. """
        projVars = self.getVariables()
        envVars = {
            "PLATFORM_APP_DIR" : "/app",
            "PLATFORM_APPLICATION" : {},
            "PLATFORM_APPLICATION_NAME" : self.config.getName(),
            "PLATFORM_BRANCH" : "",
            "PLATFORM_DOCUMENT_ROOT" : "/",
            "PLATFORM_ENVIRONMENT" : "",
            "PLATFORM_PROJECT" : self.containerId,
            "PLATFORM_RELATIONSHIPS" : base64.b64encode(json.dumps(self.relationships)),
            "PLATFORM_ROUTES" : "", # TODO
            "PLATFORM_TREE_ID" : "",
            "PLATFORM_VARIABLES" : base64.b64encode(json.dumps(projVars)),
            "PLATFORM_PROJECT_ENTROPY" : self.config.getEntropy()
        }
        envVars.update(self.getProvisioner().getEnvironmentVariables())
        varConf = {}
        for key in projVars:
            if "env:" not in key: continue
            envVars[key.replace("env:", "").strip().upper()] = projVars[key]
        return envVars

    def start(self, cmd = ""):
        """ Start docker container. """
        if self.logger:
            self.logger.logEvent(
                "Starting '%s' container." % (self.image),
                self.logIndent
            )
        # get network for docker container
        network = None
        try:
            network = self.dockerClient.networks.get(self.networkId)
        except docker.errors.NotFound:
            network = self.dockerClient.networks.create(
                self.networkId,
                driver="bridge"
            )
        # start container
        container = None
        try:
            container = self.getContainer()
            if container.status != "running":
                container.start()
        except docker.errors.NotFound:
            # create container
            # first look for committed provisioned container, if not found use unprovisioned image
            container = None
            lastExcept = None
            needProvision = True
            for image in ["%s:%s" % (self.DOCKER_COMMIT_REPO, self.getTag()), self.image]:
                if container:
                    needProvision = False
                    break
                try:
                    container = self.dockerClient.containers.run(
                        image,
                        name=self.containerId,
                        command=cmd if cmd else None,
                        detach=True,
                        volumes=self.getProvisioner().getVolumes(),
                        environment=self.getEnvironmentVariables(),
                        working_dir="/app",
                        hostname=self.containerId
                    )
                except docker.errors.ImageNotFound as e:
                    lastExcept = e
                except docker.errors.NotFound as e:
                    lastExcept = e
                except docker.errors.APIError as e:
                    lastExcept = e
            if not container:
                if lastExcept: raise lastExcept
                return
            # add to network
            network.connect(container)
            # provision container
            if needProvision:
                self.provision()
                self.commit()
        # runtime commands
        if self.logger:
            self.logger.logEvent(
                "Execute runtime commands for '%s' container." % (self.image), 
                self.logIndent
            )
        self.getProvisioner().runtime()

    def stop(self):
        if self.logger:
            self.logger.logEvent(
                "Stopping '%s' container." % (self.image),
                self.logIndent
            )
        try:
            container = self.getContainer()
            container.stop()
            container.wait()
            container.remove()
        except docker.errors.NotFound:
            pass

    def syncApp(self):
        """ Sync application files in to container. """
        if self.logger:
            self.logger.logEvent(
                "Sync application files.",
                self.logIndent
            )
        container = self.getContainer()
        container.exec_run(
            ["rsync", "-a", "--exclude", ".platform", "--exclude", ".git", "--exclude", ".platform.app.yaml", "/mnt/app/", "/app"]
        )
        container.exec_run(
            ["chown", "-R", "web:web", "/app"]
        )

    def provision(self):
        """ Provision current container. """
        if self.logger:
            self.logger.logEvent(
                "Provisioning '%s' container." % (self.image),
                self.logIndent
            )
        self.getProvisioner().provision()
        self.getContainer().restart()

    def preBuild(self):
        """ Run pre build commands. """
        self.getProvisioner().preBuild()

    def commit(self):
        """ Commit changes to container (useful after provisioning.) """
        if self.logger:
            self.logger.logEvent(
                "Commit '%s' container." % (self.image),
                self.logIndent
            )
        container = self.getContainer()
        container.commit(
            self.DOCKER_COMMIT_REPO,
            self.getTag()
        )