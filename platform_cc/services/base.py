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
            self.getName(),
            self.project.get("uid")[0:6]
        )

    def _buildContainer(self):
        pass

    def start(self):
        containerName = self.getContainerName()
        dockerImage = self.getDockerImage()
        print(containerName)

    def stop(self):
        pass