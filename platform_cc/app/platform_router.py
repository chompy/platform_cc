from config.platform_router_config import PlatformRouterConfig
from platform_docker import PlatformDocker

class PlatformRouter:
    """ Provide router to route request to specific app. """

    def __init__(self, logger = None):
        self.config = PlatformRouterConfig()
        self.docker = PlatformDocker(
            self.config,
            self.config.getName(),
            self.config.getDockerImage(),
            logger
        )
        self.logger = logger
        self.logIndent = 0

    def start(self):
        """ Start router. """
        if self.logger:
            self.logger.logEvent(
                "Starting router.",
                self.logIndent
            )
        self.docker.start(None, {"80/tcp" : 80, "443/tcp" : 443})
        self.docker.getProvisioner().copyStringToFile(
            self.docker.getProvisioner().config.get("router_conf", ""),
            "/etc/nginx/nginx.conf"
        )
        self.docker.getContainer().restart()

    def stop(self):
        """ Stop router. """
        if self.logger:
            self.logger.logEvent(
                "Stopping router.",
                self.logIndent
            )
        self.docker.stop()

    def addProject(self, project):
        """ Add project to router """
        if self.logger:
            self.logger.logEvent(
                "Add project '%s' to router." % (
                    project.projectHash[:6]
                ),
                self.logIndent
            )

        container = self.docker.getContainer()
        if not container: return
        self.docker.getProvisioner().copyStringToFile(
            project.generateRouterConfig(),
            "/router/project_%s.conf" % (
                project.projectHash
            )
        )
        self.docker.getContainer().restart()            

    def removeProject(self, project):
        """ Remove project from router. """
        if self.logger:
            self.logger.logEvent(
                "Remove project '%s' from router." % (
                    project.projectHash[:6]
                ),
                self.logIndent
            )

        self.docker.getContainer().exec_run(
                [
                    "rm",
                    "-f",
                    "/router/project_%s.conf" % (
                        project.projectHash,
                    )
                ],
                privileged=True
            )
        self.docker.getContainer().restart()