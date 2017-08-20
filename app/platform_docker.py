import os
import docker
import importlib

class PlatformDocker:

    def __init__(self, platformConfig, image):
        self.dockerClient = docker.from_env()
        self.platformConfig = platformConfig
        self.image = str(image).strip()
        self.containerId = "psh_%s_%s_%s" % (
            os.getuid(),
            self.platformConfig.getName(),
            self.image.split(":")[0]
        )
        self.networkId = "psh_%s_%s_network" % (
            os.getuid(),
            self.platformConfig.getName()
        )
    
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
            container = self.dockerClient.containers.run(
                self.image,
                name=self.containerId,
                detach=True,
                volumes={
                    os.path.realpath(self.platformConfig.projectPath) : {
                        "bind" : "/mnt/app",
                        "mode" : "ro"
                    }
                }
            )
            # add to network
            network.connect(container)
            print "done."
            # provision container
            print "> Provisioning '%s' container..." % (self.image)
            provisionModule = importlib.import_module("app.docker_provisioners.provision_%s" % self.image.split(":")[0])
            provisioner = provisionModule.DockerProvision(
                container,
                self.platformConfig
            )
            provisioner.provision()

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

