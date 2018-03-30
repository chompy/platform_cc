import docker
from dockerpty import PseudoTerminal, ExecOperation
from exception.state_error import StateError

class Container:
    """
    Base class for all Docker container instances.
    """
    
    """ Name prefix to use for all service containers. """
    CONTAINER_NAME_PREFIX = "pcc_"

    def __init__(self, project, name, dockerClient = None):
        """
        Constructor.

        :param project: Project data
        :param name: Name of this container
        :param dockerClient: Docker client
        """        
        self.project = dict(project)
        self.name = str(name)
        self.docker = dockerClient
        if not self.docker:
            self.docker = docker.from_env()

    def getName(self):
        """
        Get name of container service/application/etc.

        :return: Service or application name
        :rtype: str
        """
        return self.name

    def getDockerImage(self):
        """
        Get docker image name for container.

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
            self.name
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
        Get service Docker container if it exists.

        :return: Docker container
        :rtype: docker.client.containers.Container
        """
        try:
            return self.docker.containers.get(
                self.getContainerName()
            )
        except docker.errors.NotFound:
            pass
        return None

    def getVolume(self, name = ""):
        """
        Get a Docker volume to use with this service.

        :return: Docker volume
        :rtype: docker.client.volumes.Volume
        """
        volumeId = "%s%s_%s" %(
            self.CONTAINER_NAME_PREFIX,
            self.project.get("uid")[0:6],
            self.getName()
        )
        if name:
            volumeId = "%s_%s" % (
                volumeId,
                str(name)
            )
        try:
            return self.docker.volumes.get(volumeId)
        except docker.errors.NotFound:
            pass
        return self.docker.volumes.create(
            volumeId
        )

    def isRunning(self):
        """
        Determine if service container is currently running.

        :return: True if running
        :rtype: bool
        """
        container = self.getContainer()
        return container and container.status == "running"

    def getContainerIpAddress(self):
        """
        Get container local IP address.
        
        :return: Container IP address
        :rtype: str
        """
        container = self.getContainer()
        if not container: return ""
        return str(
            container.attrs.get(
                "NetworkSettings", {}
            ).get(
                "Networks", {}
            ).get(
                self.getNetworkName(), {}
            ).get(
                "IPAddress", ""
            )
        ).strip()

    def start(self):
        """
        Start Docker container for service.
        """
        container = self.getContainer()
        if not container:
            self.getNetwork() # instantiate if not created
            container = self.docker.containers.create(
                self.getDockerImage(),
                name = self.getContainerName(),
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
        if container.status == "running": return
        container.start()

    def stop(self):
        """
        Stop Docker container for service.
        """
        container = self.getContainer()
        if not container: return
        container.stop()
        container.wait()
        container.remove()

    def restart(self):
        """
        Restart Docker container for service.
        """
        self.stop()
        self.start()

    def shell(self, cmd = "bash", user = "root"):
        """
        Create an interactive shell inside container.

        :param cmd: Command to run
        :param user: User to run as
        """
        if not self.isRunning():
            raise StateError("Service '%s' is not running." % self.getName())
        container = self.getContainer()
        execId = self.docker.api.exec_create(
            container.id,
            cmd,
            tty = True,
            stdin = True,
            user = user
        )
        operation = ExecOperation(self.docker.api, execId)
        PseudoTerminal(self.docker.api, operation).start()