import os
import docker
import importlib
import yaml
import json
import base64
import hashlib
from app.platform_utils import log_stdout, print_stdout

class PlatformDocker:

    """ Manage and provision a platform service. """

    DOCKER_CONTAINER_NAME_PREFIX = "pcc"

    DOCKER_COMMIT_REPO = "platform_cc"

    def __init__(self, config, name = None, image = None):
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
        return provisionModule.DockerProvision(
            self.dockerClient,
            container,
            self.config,
            self.image
        )

    def getTag(self):
        """ Get unique tag name for this container's configuration. """
        return "%s_%s" % (
            self.image.split(":")[0],
            self.getProvisioner().getUid()[:6]
        )

    def getVariables(self):
        """ Get project variables. """
        varPath = os.path.join(self.config.getDataPath(), "vars.yaml")
        varConf = {}
        if os.path.exists(varPath):
            with open(varPath, "r") as f:
                varConf = yaml.load(f)
        varConf.update(self.config.getVariables())
        finalVar = {}
        for key in varConf:
            if type(varConf) is dict:
                for subKey in varConf[key]:
                    finalVar["%s:%s" % (key, subKey)] = varConf[key][subKey]
                continue
            finalVar[key] = varConf[key]
        return finalVar

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
            "PLATFORM_RELATIONSHIPS" : "", # TODO
            "PLATFORM_ROUTES" : "", # TODO
            "PLATFORM_TREE_ID" : "",
            "PLATFORM_VARIABLES" : base64.b64encode(json.dumps(projVars)),
            "PLATFORM_PROJECT_ENTROPY" : self.config.getEntropy()
        }
        envVars.update(self.getProvisioner().getEnvironmentVariables())
        varPath = os.path.join(self.config.getDataPath(), "vars.yaml")
        varConf = {}
        for key in projVars:
            if "env:" not in key: continue
            envVars[key.replace("env:", "")] = projVars[key]
        return envVars

    def start(self):
        """ Start docker container. """
        log_stdout("Starting '%s' container..." % (self.image), 1, False)
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
            print_stdout("already running.")
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
            if not container:
                if lastExcept: raise lastExcept
                return
            # add to network
            network.connect(container)
            print_stdout("done.")
            # provision container
            if needProvision:
                self.provision()
                self.commit()
        # runtime commands
        log_stdout("Execute runtime commands for '%s' container..." % (self.image), 1)
        self.getProvisioner().runtime()

    def stop(self):
        log_stdout("Stopping '%s' container..." % (self.image), 1, False)
        try:
            container = self.getContainer()
            container.stop()
            container.wait()
            container.remove()
            print_stdout("done.")
        except docker.errors.NotFound:
            print_stdout("not running, skipped.")

    def syncApp(self):
        """ Sync application files in to container. """
        log_stdout("Sync application files...", 1, False)
        container = self.getContainer()
        container.exec_run(
            ["rsync", "-a", "--exclude", ".platform", "--exclude", ".git", "--exclude", ".platform.app.yaml", "/mnt/app/", "/app"]
        )
        container.exec_run(
            ["chown", "-R", "web:web", "/app"]
        )
        print_stdout("done.")

    def provision(self):
        """ Provision current container. """
        log_stdout("Provisioning '%s' container..." % (self.image), 1)
        self.getProvisioner().provision()
        self.getContainer().restart()

    def preBuild(self):
        """ Run pre build commands. """
        self.getProvisioner().preBuild()

    def commit(self):
        """ Commit changes to container (useful after provisioning.) """
        log_stdout("Commit '%s' container..." % (self.image), 1, False)
        try:
            container = self.getContainer()
            container.commit(
                self.DOCKER_COMMIT_REPO,
                self.getTag()
            )
        except Exception as e:
            print_stdout("error.")
            raise e
        print_stdout("done.")