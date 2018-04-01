from container import Container

class BasePlatformService(Container):
    """
    Base class for Platform.sh services.
    """

    def __init__(self, project, config, dockerClient = None):
        """
        Constructor.

        :param project: Project data
        :param config: Service configuration
        :param dockerClient: Docker client
        """
        self.config = dict(config)
        Container.__init__(
            self,
            project,
            self.config.get(
                "_name",
                self.config.get(
                    "_type"
                )
            ),
            dockerClient
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

    def getServiceData(self):
        """
        Get data needed to access service for use by applications.

        :return: Dictionary containing service data
        :rtype: dict
        """
        return {
            "running"                   : self.isRunning(),
            "ip"                        : self.getContainerIpAddress(),
            "platform_relationships"    : {}
        }