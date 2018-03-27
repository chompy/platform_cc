import docker
import dockerpty

class BasePlatformService:
    """
    Base class for Platform.sh services.
    """

    """ Name prefix to use for all service containers. """
    CONTAINER_NAME_PREFIX = "pcc_"
    
    def __init__(self, project, config):
        """
        Constructor.

        :param project: Project data
        :param config: Service configuration
        """
        self.project = dict(project)
        self.config = dict(config)
        self.docker = docker.from_env()
        self._container = None

    def getName(self):
        """
        Get name of service.

        :return: Service name
        :rtype: str
        """
        return self.config.get(
            "_name",
            self.config.get(
                "_type"
            )
        )

    def getType(self):
        """
        Get service type.

        :return: Service type
        :rtype: str
        """
        return self.config.get(
            "_type"
        )

    def getDockerImage(self):
        """
        Get docker image name for service.

        :return: Docker image name
        :rtype: str
        """
        return ""

    def getContainerName(self):
        """
        Get name of service docker container.

        :return: Container name
        :rtype: str
        """
        return "%s%s_%s" % (
            self.CONTAINER_NAME_PREFIX,
            self.project.get("uid")[0:6],
            self.getName(),
        )

    def getNetworkName(self):
        """
        Get name of network to use with docker container.

        :return: Network name
        :rtype: str
        """
        return "%s%s" % (
            self.CONTAINER_NAME_PREFIX,
            self.project.get("uid")[0:6]
        )

    def getContainerCommand(self):
        """
        Get command to run to start service inside container.

        :return: Command or None if default command should be used
        """
        return None

    def getContainerEnvironmentVariables(self):
        """
        Get dictionary of environment variables to set inside container.

        :return: Dictionary of environment variables.
        :rtype: dict
        """
        return {}

    def getContainerHosts(self):
        """
        Addtional hostnames to resolve inside the container, 
        as a mapping of hostname to IP address.

        :return: Dictionary of hostname to ip map
        :rtype: dict
        """
        return {}

    def getContainerPorts(self):
        """
        Ports to bind inside the container.

        :return: Dictionary of ports to bind
        :rtype: dict
        """
        return {}

    def getContainerVolumes(self):
        """
        Dictionary of volumes to mount inside container.

        :return: Dictionary of volumes
        :rtype: dict
        """
        return {}

    def getContainerWorkingDirectory(self):
        """
        Get working directory inside container.

        :return: Working directory, or None for default
        :rtype: str
        """
        return None

    def getNetwork(self):
        """
        Get Docker network to use with service container.

        :return: Docker network
        :rtype: docker.client.networks.Network
        """
        networkName = self.getNetworkName()
        try:
            return self.docker.networks.get(
                networkName
            )
        except docker.errors.NotFound:
            pass
        return self.docker.networks.create(
            networkName
        )

    def getContainer(self):
        """
        Build Docker container for service. Retrieves container
        if it already exists.

        :return: Docker container
        :rtype: docker.client.containers.Container
        """
        containerName = self.getContainerName()
        # fetch container if already exists
        try:
            return self.docker.containers.get(containerName)
        except docker.errors.NotFound:
            pass
        self.getNetwork() # instantiate network if it does not exist
        return self.docker.containers.create(
            self.getDockerImage(),
            name = containerName,
            command = self.getContainerCommand(),
            detach = True,
            stdin_open = True,
            tty=True,
            environment = self.getContainerEnvironmentVariables(),
            extra_hosts = self.getContainerHosts(),
            network = self.getNetworkName(),
            ports = self.getContainerPorts(),
            volumes = self.getContainerVolumes(),
            working_dir = self.getContainerWorkingDirectory()
        )

    def start(self):
        """
        Start Docker container for service.
        """
        container = self.getContainer()
        if container.status == "running": return
        container.start()

    def stop(self):
        """
        Stop Docker container for service.
        """
        container = self.docker.containers.get(
            self.getContainerName()
        )
        container.stop()
        container.wait()
        container.remove()

    def restart(self):
        """
        Restart Docker container for service.
        """
        try:
            self.stop()
        except docker.errors.NotFound:
            pass
        self.start()