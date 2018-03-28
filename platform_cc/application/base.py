from container import Container

class BasePlatformApplication(Container):

    """
    Base class for managing Platform.sh applications.
    """

    def __init__(self, project, config):
        """
        Constructor.

        :param project: Project data
        :param config: Service configuration
        """
        self.config = dict(config)
        Container.__init__(
            self,
            project,
            self.config.get(
                "_name",
                ""
            )
        )        
