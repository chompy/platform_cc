from config.platform_service_config import PlatformServiceConfig
from platform_docker import PlatformDocker

class PlatformService:

    """ Platform service handler. """

    def __init__(self, projectHash, projectPath, name, logger = None):
        self.name = str(name).strip()
        self.config = PlatformServiceConfig(
            projectHash,
            projectPath,
            name
        )
        self.docker = PlatformDocker(
            self.config,
            "%s" % (self.config.getName()),
            self.config.getDockerImage(),
            logger
        )
        self.logger = logger
        self.logIndent = 0

    def start(self):
        """ Start service. """
        if self.logger:
            self.logger.logEvent(
                "Starting '%s' service." % self.config.getName(),
                self.logIndent
            )
        if not self.config.getDockerImage():
            if self.logger:
                self.logger.logEvent(
                    "No docker image available, skipped.",
                    self.logIndent + 1
                )
            return
        self.docker.start()

    def stop(self):
        """ Stop service. """
        if self.logger:
            self.logger.logEvent(
                "Stopping '%s' service." % self.config.getName(),
                self.logIndent
            )
        if not self.config.getDockerImage():
            if self.logger:
                self.logger.logEvent(
                    "No docker image available, skipped.",
                    self.logIndent + 1
                )
            return
        self.docker.stop()

    def shell(self, cmd = "bash"):
        """ Shell in to application container. """
        if self.logger:
            self.logger.logEvent(
                "Entering shell for '%s' service." % (self.config.getName()),
                self.logIndent
            )
        self.docker.shell(cmd, "root")