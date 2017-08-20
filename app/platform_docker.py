import os
import docker
import importlib
import yaml
import json
import base64

class PlatformDocker:

    """ Manage and provision a platform service. """

    DOCKER_CONTAINER_NAME_PREFIX = "pcc"

    DOCKER_COMMIT_REPO = "platform_cc"

    def __init__(self, platformConfig, image):
        self.dockerClient = docker.from_env()
        self.platformConfig = platformConfig
        self.image = str(image).strip()
        self.containerId = "%s_%s_%s_%s" % (
            self.DOCKER_CONTAINER_NAME_PREFIX,
            os.getuid(),
            self.platformConfig.getName(),
            self.image.split(":")[0]
        )
        self.networkId = "%s_%s_%s_network" % (
            self.DOCKER_CONTAINER_NAME_PREFIX,
            os.getuid(),
            self.platformConfig.getName()
        )
    
    def getTag(self):
        """ Get unique tag name for this container's configuration. """
        provisionModule = importlib.import_module("app.docker_provisioners.provision_%s" % self.image.split(":")[0])
        provisioner = provisionModule.DockerProvision(
            None,
            self.platformConfig,
            self.image
        )
        return provisioner.getUid()[:10]

    def getVariables(self):
        """ Get project variables. """
        varPath = os.path.join(self.platformConfig.getDataPath(), "vars.yaml")
        varConf = {}
        if os.path.exists(varPath):
            with open(varPath, "r") as f:
                varConf = yaml.load(f)
        varConf.update(self.platformConfig.getVariables())
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
            "PLATFORM_APPLICATION_NAME" : self.platformConfig.getName(),
            "PLATFORM_BRANCH" : "",
            "PLATFORM_DOCUMENT_ROOT" : "/",
            "PLATFORM_ENVIRONMENT" : "",
            "PLATFORM_PROJECT" : self.containerId,
            "PLATFORM_RELATIONSHIPS" : "", # TODO
            "PLATFORM_ROUTES" : "", # TODO
            "PLATFORM_TREE_ID" : "",
            "PLATFORM_VARIABLES" : base64.b64encode(json.dumps(projVars)),
            "PLATFORM_PROJECT_ENTROPY" : self.platformConfig.getEntropy()
        }
        varPath = os.path.join(self.platformConfig.getDataPath(), "vars.yaml")
        varConf = {}
        for key in projVars:
            if "env:" not in key: continue
            envVars[key.replace("env:", "")] = projVars[key]
        return envVars

    def start(self):
        """ Start docker container. """
        print "> Starting '%s' container..." % (self.image),
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
            container = self.dockerClient.containers.get(self.containerId)
            print "already running."
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
                        volumes={
                            os.path.realpath(self.platformConfig.projectPath) : {
                                "bind" : "/mnt/app",
                                "mode" : "ro"
                            }
                        },
                        environment=self.getEnvironmentVariables(),
                        working_dir="/app"
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
            print "done."
            # provision container
            if needProvision:
                self.provision()
                self.commit()

    def stop(self):
        print "> Stopping '%s' container..." % (self.image),
        try:
            container = self.dockerClient.containers.get(self.containerId)
            container.stop()
            container.wait()
            container.remove()
            print "done."
        except docker.errors.NotFound:
            print "not running, skipping."

    def provision(self):
        """ Provision current container. """
        print "> Provisioning '%s' container..." % (self.image)
        container = self.dockerClient.containers.get(self.containerId)
        provisionModule = importlib.import_module("app.docker_provisioners.provision_%s" % self.image.split(":")[0])
        provisioner = provisionModule.DockerProvision(
            container,
            self.platformConfig,
            self.image
        )
        provisioner.provision()

    def commit(self):
        """ Commit changes to container (useful after provisioning.) """
        print "> Commit '%s' container..." % (self.image),
        try:
            container = self.dockerClient.containers.get(self.containerId)
            container.commit(
                self.DOCKER_COMMIT_REPO,
                self.getTag()
            )
        except Exception as e:
            print "error."
            raise e
        print "done."